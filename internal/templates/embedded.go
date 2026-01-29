// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// A cross-platform CLI tool for managing AI agent templates

package templates

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// CurrentVersion is the version of the embedded templates
const CurrentVersion = "2.0.0"

// Templates holds all loaded agent templates
type Templates struct {
	Version   string
	Agents    map[string]Agent
	Skills    map[string]Skill
	Workflows map[string]Workflow
}

// Agent represents a specialist agent
type Agent struct {
	Name        string
	Description string
	Skills      []string
	Tools       []string
	Content     string // full markdown content
}

// Skill represents a domain skill
type Skill struct {
	Name        string
	Description string
	Content     string
	Scripts     []string // available scripts
}

// Workflow represents a slash command workflow
type Workflow struct {
	Name        string
	Description string
	Content     string
}

//go:embed all:data
var embeddedFS embed.FS

// LoadEmbedded loads templates from the embedded filesystem.
//
// How it works:
// 1. Walk the embedded "data" directory
// 2. Parse each .md file looking for YAML frontmatter
// 3. Extract metadata (name, description, skills, etc.)
// 4. Build the Templates struct with all loaded data
//
// Why embedded? The binary is self-contained - no need to download
// templates on first run. Users can still update from network later.
func LoadEmbedded() (*Templates, error) {
	tmpl := &Templates{
		Version:   CurrentVersion,
		Agents:    make(map[string]Agent),
		Skills:    make(map[string]Skill),
		Workflows: make(map[string]Workflow),
	}

	// Load agents
	if err := loadAgents(tmpl); err != nil {
		return nil, err
	}

	// load skills
	if err := loadSkills(tmpl); err != nil {
		return nil, err
	}

	// load workflows
	if err := loadWorkflows(tmpl); err != nil {
		return nil, err
	}

	return tmpl, nil
}

// loadAgents reads all agent files from embedded FS
func loadAgents(tmpl *Templates) error {
	agentsDir := "data/agents"

	entries, err := embeddedFS.ReadDir(agentsDir)
	if err != nil {
		// agents dir might not exist in embedded data yet
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		content, err := embeddedFS.ReadFile(filepath.Join(agentsDir, entry.Name()))
		if err != nil {
			continue
		}

		agent := parseAgentFile(string(content))
		name := strings.TrimSuffix(entry.Name(), ".md")
		agent.Name = name
		tmpl.Agents[name] = agent
	}

	return nil
}

// loadSkills reads all skill directories from embedded FS
func loadSkills(tmpl *Templates) error {
	skillsDir := "data/skills"

	entries, err := embeddedFS.ReadDir(skillsDir)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillFile := filepath.Join(skillsDir, entry.Name(), "SKILL.md")
		content, err := embeddedFS.ReadFile(skillFile)
		if err != nil {
			continue
		}

		skill := parseSkillFile(string(content))
		skill.Name = entry.Name()
		tmpl.Skills[entry.Name()] = skill
	}

	return nil
}

// loadWorkflows reads all workflow files from embedded FS
func loadWorkflows(tmpl *Templates) error {
	workflowsDir := "data/workflows"

	entries, err := embeddedFS.ReadDir(workflowsDir)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		content, err := embeddedFS.ReadFile(filepath.Join(workflowsDir, entry.Name()))
		if err != nil {
			continue
		}

		workflow := parseWorkflowFile(string(content))
		name := strings.TrimSuffix(entry.Name(), ".md")
		workflow.Name = name
		tmpl.Workflows[name] = workflow
	}

	return nil
}

// parseFrontmatter extracts YAML frontmatter from markdown
func parseFrontmatter(content string) (map[string]interface{}, string) {
	if !strings.HasPrefix(content, "---") {
		return nil, content
	}

	parts := strings.SplitN(content[3:], "---", 2)
	if len(parts) < 2 {
		return nil, content
	}

	var frontmatter map[string]interface{}
	if err := yaml.Unmarshal([]byte(parts[0]), &frontmatter); err != nil {
		return nil, content
	}

	return frontmatter, strings.TrimSpace(parts[1])
}

func parseAgentFile(content string) Agent {
	fm, body := parseFrontmatter(content)

	agent := Agent{
		Content: content,
	}

	if fm != nil {
		if desc, ok := fm["description"].(string); ok {
			agent.Description = desc
		}
		if skills, ok := fm["skills"].(string); ok {
			agent.Skills = strings.Split(skills, ",")
			for i := range agent.Skills {
				agent.Skills[i] = strings.TrimSpace(agent.Skills[i])
			}
		}
		if tools, ok := fm["tools"].(string); ok {
			agent.Tools = strings.Split(tools, ",")
			for i := range agent.Tools {
				agent.Tools[i] = strings.TrimSpace(agent.Tools[i])
			}
		}
	}

	// Extract description from first paragraph if not in frontmatter
	if agent.Description == "" && body != "" {
		lines := strings.Split(body, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				agent.Description = line
				break
			}
		}
	}

	return agent
}

func parseSkillFile(content string) Skill {
	fm, _ := parseFrontmatter(content)

	skill := Skill{
		Content: content,
	}

	if fm != nil {
		if desc, ok := fm["description"].(string); ok {
			skill.Description = desc
		}
	}

	return skill
}

func parseWorkflowFile(content string) Workflow {
	fm, _ := parseFrontmatter(content)

	workflow := Workflow{
		Content: content,
	}

	if fm != nil {
		if desc, ok := fm["description"].(string); ok {
			workflow.Description = desc
		}
	}

	return workflow
}

// Filter returns a new Templates with only the specified agents and skills.
// if empty slices are passed, all are included.
func (t *Templates) Filter(agents []string, skills []string) *Templates {
	filtered := &Templates{
		Version:   t.Version,
		Agents:    make(map[string]Agent),
		Skills:    make(map[string]Skill),
		Workflows: t.Workflows, // always include all workflows
	}

	// filter agents
	if len(agents) == 0 {
		filtered.Agents = t.Agents
	} else {
		for _, name := range agents {
			if agent, ok := t.Agents[name]; ok {
				filtered.Agents[name] = agent
			}
		}
	}

	// filter skills
	if len(skills) == 0 {
		filtered.Skills = t.Skills
	} else {
		for _, name := range skills {
			if skill, ok := t.Skills[name]; ok {
				filtered.Skills[name] = skill
			}
		}
	}

	return filtered
}

// InstallTo copies templates to the specified directory.
// creates the directory structure and writes all files.
func (t *Templates) InstallTo(targetDir string) error {
	// Create directories
	dirs := []string{
		filepath.Join(targetDir, "agents"),
		filepath.Join(targetDir, "skills"),
		filepath.Join(targetDir, "workflows"),
		filepath.Join(targetDir, "rules"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Write agents
	for name, agent := range t.Agents {
		file := filepath.Join(targetDir, "agents", name+".md")
		if err := os.WriteFile(file, []byte(agent.Content), 0644); err != nil {
			return err
		}
	}

	// Write skills
	for name, skill := range t.Skills {
		skillDir := filepath.Join(targetDir, "skills", name)
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			return err
		}
		file := filepath.Join(skillDir, "SKILL.md")
		if err := os.WriteFile(file, []byte(skill.Content), 0644); err != nil {
			return err
		}
	}

	// Write workflows
	for name, workflow := range t.Workflows {
		file := filepath.Join(targetDir, "workflows", name+".md")
		if err := os.WriteFile(file, []byte(workflow.Content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// GetLatestVersion returns the current embedded version
func GetLatestVersion() string {
	return CurrentVersion
}

// CopyEmbeddedToFS copies embedded templates to the real filesystem.
// used for initial installation or recovery.
func CopyEmbeddedToFS(targetDir string) error {
	return fs.WalkDir(embeddedFS, "data", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate target path by removing "data" prefix
		relPath := strings.TrimPrefix(path, "data")
		if relPath == "" {
			return nil
		}
		targetPath := filepath.Join(targetDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		content, err := embeddedFS.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, content, 0644)
	})
}
