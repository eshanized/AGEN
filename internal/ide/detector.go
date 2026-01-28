// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// A cross-platform CLI tool for managing AI agent templates

package ide

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/eshanized/agen/internal/templates"
)

// Adapter is the interface that each IDE must implement.
// This allows us to support multiple IDEs with different config formats.
//
// Why an interface? Each IDE stores config differently:
// - Antigravity uses .agent/ folder with multiple files
// - Cursor uses a single .cursorrules file
// - Windsurf uses a single .windsurfrules file
// - Zed uses .zed/settings.json
//
// The adapter pattern lets us handle these differences cleanly.
type Adapter interface {
	// Name returns the human-readable name of the IDE
	Name() string

	// Detect checks if this IDE is being used in the given project
	// Returns true if IDE-specific files/folders are found
	Detect(projectPath string) bool

	// Install copies templates to the project in IDE-specific format
	Install(tmpl *templates.Templates, opts InstallOptions) error

	// Update updates installed templates with conflict detection
	Update(tmpl *templates.Templates, opts UpdateOptions) (*UpdateChanges, error)

	// GetRulesPath returns the path to the main rules/config file
	GetRulesPath() string
}

// InstallOptions configures template installation
type InstallOptions struct {
	TargetDir string
	DryRun    bool
	Force     bool // overwrite without prompting
	Verbose   bool
}

// UpdateOptions configures template updates
type UpdateOptions struct {
	TargetDir string
	DryRun    bool
	Force     bool // overwrite modified files
	Verbose   bool
}

// UpdateChanges tracks what changed during an update
type UpdateChanges struct {
	Added   []string // new files
	Updated []string // changed files
	Skipped []string // user-modified, not overwritten
}

// InstalledInfo contains info about installed templates
type InstalledInfo struct {
	IDE           string
	Version       string
	AgentCount    int
	SkillCount    int
	WorkflowCount int
	ModifiedFiles int
	Agents        []string // list of installed agent names
	Skills        []string // list of installed skill names
}

// adapters holds all registered IDE adapters
var adapters = make(map[string]Adapter)

// RegisterAdapter adds an IDE adapter to the registry
func RegisterAdapter(name string, adapter Adapter) {
	adapters[name] = adapter
}

// GetAdapter returns the adapter for a specific IDE name.
// returns nil if not found.
func GetAdapter(name string) Adapter {
	return adapters[name]
}

// Detect attempts to auto-detect which IDE is being used in the project.
//
// How it works:
// 1. First we look for IDE-specific marker files (.cursorrules, .windsurfrules, etc)
// 2. If none found, check for .agent/ folder (means Antigravity/Claude)
// 3. As a fallback, we check environment variables that IDEs often set
// 4. If all else fails, return nil and let the user pick manually
//
// Why this order? IDE-specific files are most reliable because users explicitly
// created them. The .agent folder might exist from a previous agen run.
// Environment vars are least reliable since they vary by IDE version.
//
// Returns: detected IDE adapter or nil if no IDE detected
func Detect(projectPath string) Adapter {
	// Priority order for detection
	// check for explicit IDE config files first
	detectionOrder := []string{
		"cursor",      // .cursorrules
		"windsurf",    // .windsurfrules
		"zed",         // .zed/
		"antigravity", // .agent/ (check last since other IDEs might also have agents)
	}

	for _, name := range detectionOrder {
		if adapter, ok := adapters[name]; ok {
			if adapter.Detect(projectPath) {
				return adapter
			}
		}
	}

	return nil
}

// GetInstalledInfo retrieves information about installed templates.
// this is used for status and health commands.
func GetInstalledInfo(projectPath string, adapter Adapter) (*InstalledInfo, error) {
	info := &InstalledInfo{
		IDE:     adapter.Name(),
		Version: "unknown",
	}

	// For Antigravity/Zed/Windsurf (folder based or similar), we can look in .agent
	// For Cursor/Windsurf (single file), we might need to parse the file
	// Since all adapters currently populate .agent/ internally or are based on it,
	// checking .agent/ is a reasonable default, but let's be more specific for single-file IDEs.

	agentDir := filepath.Join(projectPath, ".agent", "agents")
	skillDir := filepath.Join(projectPath, ".agent", "skills")
	workflowDir := filepath.Join(projectPath, ".agent", "workflows")

	if adapter.Name() == "Cursor" {
		// Check .cursorrules content
		content, err := os.ReadFile(filepath.Join(projectPath, ".cursorrules"))
		if err == nil {
			s := string(content)
			info.AgentCount = strings.Count(s, "### ") // Rough heuristic
			info.SkillCount = strings.Count(s, "- **") // Rough heuristic
		}
	} else if adapter.Name() == "Windsurf" {
		// Check .windsurfrules content
		content, err := os.ReadFile(filepath.Join(projectPath, ".windsurfrules"))
		if err == nil {
			s := string(content)
			info.AgentCount = strings.Count(s, "### ")
			info.SkillCount = strings.Count(s, "- **")
		}
	} else {
		// Default: check directories
		info.Version = "1.0.0"
	}

	// Count agents
	if entries, err := os.ReadDir(agentDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() && filepath.Ext(e.Name()) == ".md" {
				info.AgentCount++
				info.Agents = append(info.Agents, e.Name())
			}
		}
	}

	// count skills
	if entries, err := os.ReadDir(skillDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				info.SkillCount++
				info.Skills = append(info.Skills, e.Name())
			}
		}
	}

	// count workflows
	if entries, err := os.ReadDir(workflowDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() && filepath.Ext(e.Name()) == ".md" {
				info.WorkflowCount++
			}
		}
	}

	return info, nil
}

// init registers all built-in adapters
func init() {
	RegisterAdapter("antigravity", &AntigravityAdapter{})
	RegisterAdapter("cursor", &CursorAdapter{})
	RegisterAdapter("windsurf", &WindsurfAdapter{})
	RegisterAdapter("zed", &ZedAdapter{})
}
