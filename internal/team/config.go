// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Team collaboration features

package team

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// TeamConfig represents shared team configuration
type TeamConfig struct {
	Name           string            `json:"name"`
	Version        string            `json:"version"`
	Description    string            `json:"description,omitempty"`
	Maintainers    []string          `json:"maintainers,omitempty"`
	RequiredAgents []string          `json:"required_agents,omitempty"`
	RequiredSkills []string          `json:"required_skills,omitempty"`
	LockedVersions map[string]string `json:"locked_versions,omitempty"`
	Settings       TeamSettings      `json:"settings"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

// TeamSettings contains team-wide settings
type TeamSettings struct {
	EnforceAgents  bool   `json:"enforce_agents"`
	EnforceSkills  bool   `json:"enforce_skills"`
	AllowPlugins   bool   `json:"allow_plugins"`
	DefaultIDE     string `json:"default_ide,omitempty"`
	TemplateSource string `json:"template_source,omitempty"`
	SyncInterval   string `json:"sync_interval,omitempty"`
}

// TeamMember represents a team member
type TeamMember struct {
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

const teamConfigFile = ".agen-team.json"

// InitTeam initializes team configuration in a project
func InitTeam(dir, name string) (*TeamConfig, error) {
	configPath := filepath.Join(dir, teamConfigFile)

	if _, err := os.Stat(configPath); err == nil {
		return nil, fmt.Errorf("team config already exists")
	}

	config := &TeamConfig{
		Name:           name,
		Version:        "1.0.0",
		RequiredAgents: []string{},
		RequiredSkills: []string{},
		LockedVersions: make(map[string]string),
		Settings: TeamSettings{
			EnforceAgents: false,
			EnforceSkills: false,
			AllowPlugins:  true,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := config.Save(dir); err != nil {
		return nil, err
	}

	return config, nil
}

// LoadTeamConfig loads team configuration from a project
func LoadTeamConfig(dir string) (*TeamConfig, error) {
	configPath := filepath.Join(dir, teamConfigFile)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("no team config found: %w", err)
	}

	var config TeamConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("invalid team config: %w", err)
	}

	return &config, nil
}

// Save writes the team config to disk
func (c *TeamConfig) Save(dir string) error {
	c.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	configPath := filepath.Join(dir, teamConfigFile)
	return os.WriteFile(configPath, data, 0644)
}

// AddRequired adds a required agent or skill
func (c *TeamConfig) AddRequired(itemType, name string) error {
	switch itemType {
	case "agent":
		for _, a := range c.RequiredAgents {
			if a == name {
				return fmt.Errorf("agent already required: %s", name)
			}
		}
		c.RequiredAgents = append(c.RequiredAgents, name)

	case "skill":
		for _, s := range c.RequiredSkills {
			if s == name {
				return fmt.Errorf("skill already required: %s", name)
			}
		}
		c.RequiredSkills = append(c.RequiredSkills, name)

	default:
		return fmt.Errorf("unknown type: %s", itemType)
	}

	return nil
}

// RemoveRequired removes a required agent or skill
func (c *TeamConfig) RemoveRequired(itemType, name string) error {
	switch itemType {
	case "agent":
		for i, a := range c.RequiredAgents {
			if a == name {
				c.RequiredAgents = append(c.RequiredAgents[:i], c.RequiredAgents[i+1:]...)
				return nil
			}
		}

	case "skill":
		for i, s := range c.RequiredSkills {
			if s == name {
				c.RequiredSkills = append(c.RequiredSkills[:i], c.RequiredSkills[i+1:]...)
				return nil
			}
		}
	}

	return fmt.Errorf("not found: %s", name)
}

// LockVersion locks a specific template version
func (c *TeamConfig) LockVersion(name, version string) {
	if c.LockedVersions == nil {
		c.LockedVersions = make(map[string]string)
	}
	c.LockedVersions[name] = version
}

// UnlockVersion removes a version lock
func (c *TeamConfig) UnlockVersion(name string) {
	delete(c.LockedVersions, name)
}

// Sync synchronizes local installation with team config
type SyncResult struct {
	Added   []string `json:"added"`
	Updated []string `json:"updated"`
	Removed []string `json:"removed"`
	Errors  []string `json:"errors"`
}

// Sync synchronizes the project with team requirements
func (c *TeamConfig) Sync(projectDir string) (*SyncResult, error) {
	result := &SyncResult{
		Added:   []string{},
		Updated: []string{},
		Removed: []string{},
		Errors:  []string{},
	}

	// Create .agent directories if needed
	agentDir := filepath.Join(projectDir, ".agent", "agents")
	skillDir := filepath.Join(projectDir, ".agent", "skills")
	os.MkdirAll(agentDir, 0755)
	os.MkdirAll(skillDir, 0755)

	// Check required agents are installed
	for _, agent := range c.RequiredAgents {
		agentPath := filepath.Join(agentDir, agent+".md")
		if _, err := os.Stat(agentPath); os.IsNotExist(err) {
			// Try to copy from embedded templates
			if err := copyEmbeddedAgent(agent, agentPath); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("agent:%s - %v", agent, err))
			} else {
				result.Added = append(result.Added, "agent:"+agent)
			}
		}
	}

	// Check required skills are installed
	for _, skill := range c.RequiredSkills {
		skillPath := filepath.Join(skillDir, skill)
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			// Try to copy from embedded templates
			if err := copyEmbeddedSkill(skill, skillPath); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("skill:%s - %v", skill, err))
			} else {
				result.Added = append(result.Added, "skill:"+skill)
			}
		}
	}

	return result, nil
}

// copyEmbeddedAgent copies an agent from embedded templates
func copyEmbeddedAgent(name, destPath string) error {
	// Read from embedded templates data
	srcPath := filepath.Join("agents", name+".md")
	content, err := readEmbeddedFile(srcPath)
	if err != nil {
		return fmt.Errorf("agent not found in templates: %s", name)
	}
	return os.WriteFile(destPath, content, 0644)
}

// copyEmbeddedSkill copies a skill from embedded templates
func copyEmbeddedSkill(name, destPath string) error {
	os.MkdirAll(destPath, 0755)
	srcPath := filepath.Join("skills", name, "SKILL.md")
	content, err := readEmbeddedFile(srcPath)
	if err != nil {
		return fmt.Errorf("skill not found in templates: %s", name)
	}
	return os.WriteFile(filepath.Join(destPath, "SKILL.md"), content, 0644)
}

// readEmbeddedFile reads a file from embedded templates
func readEmbeddedFile(path string) ([]byte, error) {
	// Try to read from internal/templates/data
	fullPath := filepath.Join("internal", "templates", "data", path)
	return os.ReadFile(fullPath)
}

// Validate checks if the project meets team requirements
type ValidationResult struct {
	Valid    bool     `json:"valid"`
	Missing  []string `json:"missing"`
	Warnings []string `json:"warnings"`
}

func (c *TeamConfig) Validate(projectDir string) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Missing:  []string{},
		Warnings: []string{},
	}

	// Check required agents
	agentDir := filepath.Join(projectDir, ".agent", "agents")
	for _, agent := range c.RequiredAgents {
		agentPath := filepath.Join(agentDir, agent+".md")
		if _, err := os.Stat(agentPath); os.IsNotExist(err) {
			result.Missing = append(result.Missing, "agent:"+agent)
			if c.Settings.EnforceAgents {
				result.Valid = false
			}
		}
	}

	// Check required skills
	skillDir := filepath.Join(projectDir, ".agent", "skills")
	for _, skill := range c.RequiredSkills {
		skillPath := filepath.Join(skillDir, skill)
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			result.Missing = append(result.Missing, "skill:"+skill)
			if c.Settings.EnforceSkills {
				result.Valid = false
			}
		}
	}

	return result
}
