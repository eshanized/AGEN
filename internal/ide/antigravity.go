// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// A cross-platform CLI tool for managing AI agent templates

package ide

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/eshanized/agen/internal/templates"
)

// AntigravityAdapter handles Antigravity/Claude Code format.
// This is the "native" format - uses full .agent/ folder structure
// with agents/, skills/, workflows/, and rules/ subdirectories.
type AntigravityAdapter struct{}

// Name returns the adapter name
func (a *AntigravityAdapter) Name() string {
	return "Antigravity"
}

// Detect checks if this is an Antigravity project.
// We look for .agent/ folder or GEMINI.md in the project root.
func (a *AntigravityAdapter) Detect(projectPath string) bool {
	// check for .agent directory
	agentDir := filepath.Join(projectPath, ".agent")
	if info, err := os.Stat(agentDir); err == nil && info.IsDir() {
		return true
	}

	// also check for GEMINI.md in root (older convention)
	geminiFile := filepath.Join(projectPath, "GEMINI.md")
	if _, err := os.Stat(geminiFile); err == nil {
		return true
	}

	return false
}

// Install copies all templates to the .agent/ folder.
//
// How it works:
// 1. Create .agent/ directory if it doesn't exist
// 2. Copy agents/, skills/, workflows/, rules/ folders
// 3. Handle conflicts if --force not set
//
// The Antigravity format is the most complete - it includes everything.
// Other adapters convert FROM this format to their specific format.
func (a *AntigravityAdapter) Install(tmpl *templates.Templates, opts InstallOptions) error {
	agentDir := filepath.Join(opts.TargetDir, ".agent")

	// check if already exists
	if _, err := os.Stat(agentDir); err == nil && !opts.Force {
		if opts.DryRun {
			return nil
		}
		return fmt.Errorf("target directory %s already exists (use --force to overwrite)", agentDir)
	}

	// Create directory structure
	dirs := []string{
		agentDir,
		filepath.Join(agentDir, "agents"),
		filepath.Join(agentDir, "skills"),
		filepath.Join(agentDir, "workflows"),
		filepath.Join(agentDir, "rules"),
		filepath.Join(agentDir, "scripts"),
	}

	for _, dir := range dirs {
		if opts.DryRun {
			continue
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Install templates
	if !opts.DryRun {
		return tmpl.InstallTo(agentDir)
	}

	return nil
}

// Update updates installed templates with conflict detection.
func (a *AntigravityAdapter) Update(tmpl *templates.Templates, opts UpdateOptions) (*UpdateChanges, error) {
	changes := &UpdateChanges{
		Added:   []string{},
		Updated: []string{},
		Skipped: []string{},
	}

	agentDir := filepath.Join(opts.TargetDir, ".agent")

	// Compare agents
	for name, agent := range tmpl.Agents {
		agentPath := filepath.Join(agentDir, "agents", name+".md")
		if _, err := os.Stat(agentPath); os.IsNotExist(err) {
			// New agent - add it
			if !opts.DryRun {
				os.WriteFile(agentPath, []byte(agent.Content), 0644)
			}
			changes.Added = append(changes.Added, "agents/"+name+".md")
		} else {
			// Existing agent - check if different
			existing, err := os.ReadFile(agentPath)
			if err == nil && string(existing) != agent.Content {
				if opts.Force {
					if !opts.DryRun {
						os.WriteFile(agentPath, []byte(agent.Content), 0644)
					}
					changes.Updated = append(changes.Updated, "agents/"+name+".md")
				} else {
					changes.Skipped = append(changes.Skipped, "agents/"+name+".md (modified locally)")
				}
			}
		}
	}

	// Compare skills
	for name, skill := range tmpl.Skills {
		skillPath := filepath.Join(agentDir, "skills", name, "SKILL.md")
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			// New skill - add it
			if !opts.DryRun {
				os.MkdirAll(filepath.Dir(skillPath), 0755)
				os.WriteFile(skillPath, []byte(skill.Content), 0644)
			}
			changes.Added = append(changes.Added, "skills/"+name+"/SKILL.md")
		} else {
			// Existing skill - check if different
			existing, err := os.ReadFile(skillPath)
			if err == nil && string(existing) != skill.Content {
				if opts.Force {
					if !opts.DryRun {
						os.WriteFile(skillPath, []byte(skill.Content), 0644)
					}
					changes.Updated = append(changes.Updated, "skills/"+name+"/SKILL.md")
				} else {
					changes.Skipped = append(changes.Skipped, "skills/"+name+"/SKILL.md (modified locally)")
				}
			}
		}
	}

	return changes, nil
}

// GetRulesPath returns the path to the rules file
func (a *AntigravityAdapter) GetRulesPath() string {
	return ".agent/rules/GEMINI.md"
}
