// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Plugin system for custom agents and skills

package plugin

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Plugin represents an installed plugin
type Plugin struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Source      string            `json:"source"`
	Type        PluginType        `json:"type"`
	InstalledAt string            `json:"installed_at"`
	Agents      []string          `json:"agents,omitempty"`
	Skills      []string          `json:"skills,omitempty"`
	Workflows   []string          `json:"workflows,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// PluginType indicates what kind of plugin this is
type PluginType string

const (
	PluginTypeAgent    PluginType = "agent"
	PluginTypeSkill    PluginType = "skill"
	PluginTypeWorkflow PluginType = "workflow"
	PluginTypeBundle   PluginType = "bundle" // contains multiple types
)

// Manager handles plugin installation and management
type Manager struct {
	pluginDir string
	registry  *Registry
}

// Registry stores information about installed plugins
type Registry struct {
	Plugins map[string]*Plugin `json:"plugins"`
	path    string
}

// NewManager creates a new plugin manager
func NewManager() (*Manager, error) {
	// Get user config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	pluginDir := filepath.Join(configDir, "agen", "plugins")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return nil, err
	}

	registry, err := loadRegistry(pluginDir)
	if err != nil {
		registry = &Registry{
			Plugins: make(map[string]*Plugin),
			path:    filepath.Join(pluginDir, "registry.json"),
		}
	}

	return &Manager{
		pluginDir: pluginDir,
		registry:  registry,
	}, nil
}

// Install installs a plugin from a source
//
// Supported sources:
// - GitHub: github.com/user/repo
// - Local: path/to/plugin
// - URL: https://example.com/plugin.zip
func (m *Manager) Install(source string) (*Plugin, error) {
	var plugin *Plugin
	var err error

	if strings.HasPrefix(source, "github.com/") {
		plugin, err = m.installFromGitHub(source)
	} else if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		plugin, err = m.installFromURL(source)
	} else {
		plugin, err = m.installFromLocal(source)
	}

	if err != nil {
		return nil, err
	}

	// Register the plugin
	m.registry.Plugins[plugin.Name] = plugin
	if err := m.registry.save(); err != nil {
		return nil, fmt.Errorf("failed to save registry: %w", err)
	}

	return plugin, nil
}

// installFromGitHub clones a plugin from GitHub
func (m *Manager) installFromGitHub(source string) (*Plugin, error) {
	// Parse source: github.com/user/repo[@version]
	parts := strings.Split(source, "@")
	repoPath := strings.TrimPrefix(parts[0], "github.com/")
	version := "main"
	if len(parts) > 1 {
		version = parts[1]
	}

	// Clone directory name
	repoParts := strings.Split(repoPath, "/")
	if len(repoParts) < 2 {
		return nil, fmt.Errorf("invalid GitHub source: %s", source)
	}
	pluginName := repoParts[1]

	targetDir := filepath.Join(m.pluginDir, pluginName)

	// Clone or update
	if _, err := os.Stat(targetDir); err == nil {
		// Already exists, pull updates
		cmd := exec.Command("git", "-C", targetDir, "pull", "origin", version)
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("failed to update: %w", err)
		}
	} else {
		// Clone fresh
		gitURL := fmt.Sprintf("https://github.com/%s.git", repoPath)
		cmd := exec.Command("git", "clone", "--depth", "1", "--branch", version, gitURL, targetDir)
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("failed to clone: %w", err)
		}
	}

	// Load plugin metadata
	return m.loadPluginMetadata(targetDir)
}

// installFromURL downloads and extracts a plugin from a URL
func (m *Manager) installFromURL(source string) (*Plugin, error) {
	// Create temp directory for download
	tempDir, err := os.MkdirTemp("", "agen-plugin-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Download the file
	resp, err := http.Get(source)
	if err != nil {
		return nil, fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Determine filename from URL or Content-Disposition
	filename := filepath.Base(source)
	if filename == "" || filename == "/" {
		filename = "plugin.zip"
	}

	// Save to temp file
	tempFile := filepath.Join(tempDir, filename)
	out, err := os.Create(tempFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Extract if zip
	if strings.HasSuffix(filename, ".zip") {
		extractDir := filepath.Join(tempDir, "extracted")
		if err := extractZip(tempFile, extractDir); err != nil {
			return nil, fmt.Errorf("failed to extract: %w", err)
		}

		// Find the plugin directory (first directory with plugin.json or agents/)
		entries, _ := os.ReadDir(extractDir)
		pluginSrc := extractDir
		for _, e := range entries {
			if e.IsDir() {
				pluginSrc = filepath.Join(extractDir, e.Name())
				break
			}
		}

		// Copy to plugins directory
		pluginName := filepath.Base(strings.TrimSuffix(filename, ".zip"))
		targetDir := filepath.Join(m.pluginDir, pluginName)
		if err := copyDir(pluginSrc, targetDir); err != nil {
			return nil, fmt.Errorf("failed to install: %w", err)
		}

		return m.loadPluginMetadata(targetDir)
	}

	return nil, fmt.Errorf("unsupported file format: %s", filename)
}

// extractZip extracts a zip file to a directory
func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	os.MkdirAll(dest, 0755)

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

// copyDir copies a directory recursively
func copyDir(src, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(src, path)
		destPath := filepath.Join(dest, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(destPath, data, info.Mode())
	})
}

// installFromLocal installs from a local path
func (m *Manager) installFromLocal(source string) (*Plugin, error) {
	absPath, err := filepath.Abs(source)
	if err != nil {
		return nil, err
	}

	// Check path exists
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("path not found: %s", source)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("source must be a directory: %s", source)
	}

	// Load plugin metadata
	return m.loadPluginMetadata(absPath)
}

// loadPluginMetadata reads plugin.json from a plugin directory
func (m *Manager) loadPluginMetadata(dir string) (*Plugin, error) {
	metadataPath := filepath.Join(dir, "plugin.json")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		// Try to infer from directory structure
		return m.inferPluginMetadata(dir)
	}

	var plugin Plugin
	if err := json.Unmarshal(data, &plugin); err != nil {
		return nil, fmt.Errorf("invalid plugin.json: %w", err)
	}

	return &plugin, nil
}

// inferPluginMetadata creates metadata from directory structure
func (m *Manager) inferPluginMetadata(dir string) (*Plugin, error) {
	name := filepath.Base(dir)
	plugin := &Plugin{
		Name:    name,
		Version: "0.0.0",
		Type:    PluginTypeBundle,
		Source:  dir,
	}

	// Scan for agents, skills, workflows
	if entries, err := os.ReadDir(filepath.Join(dir, "agents")); err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
				plugin.Agents = append(plugin.Agents, strings.TrimSuffix(e.Name(), ".md"))
			}
		}
	}

	if entries, err := os.ReadDir(filepath.Join(dir, "skills")); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				plugin.Skills = append(plugin.Skills, e.Name())
			}
		}
	}

	if entries, err := os.ReadDir(filepath.Join(dir, "workflows")); err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
				plugin.Workflows = append(plugin.Workflows, strings.TrimSuffix(e.Name(), ".md"))
			}
		}
	}

	return plugin, nil
}

// Uninstall removes a plugin
func (m *Manager) Uninstall(name string) error {
	plugin, ok := m.registry.Plugins[name]
	if !ok {
		return fmt.Errorf("plugin not found: %s", name)
	}

	// Remove plugin directory
	pluginDir := filepath.Join(m.pluginDir, name)
	if err := os.RemoveAll(pluginDir); err != nil {
		return fmt.Errorf("failed to remove plugin: %w", err)
	}

	// Update registry
	delete(m.registry.Plugins, name)
	if err := m.registry.save(); err != nil {
		return err
	}

	_ = plugin // for future use
	return nil
}

// List returns all installed plugins
func (m *Manager) List() []*Plugin {
	plugins := make([]*Plugin, 0, len(m.registry.Plugins))
	for _, p := range m.registry.Plugins {
		plugins = append(plugins, p)
	}
	return plugins
}

// Get returns a specific plugin
func (m *Manager) Get(name string) (*Plugin, error) {
	plugin, ok := m.registry.Plugins[name]
	if !ok {
		return nil, fmt.Errorf("plugin not found: %s", name)
	}
	return plugin, nil
}

// Create initializes a new plugin project
func (m *Manager) Create(name, pluginType string) (string, error) {
	targetDir := filepath.Join(".", name)

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}

	// Create plugin.json
	plugin := Plugin{
		Name:        name,
		Version:     "0.1.0",
		Description: fmt.Sprintf("Custom %s plugin", pluginType),
		Author:      "",
		Type:        PluginType(pluginType),
		Metadata:    map[string]string{},
	}

	data, _ := json.MarshalIndent(plugin, "", "  ")
	if err := os.WriteFile(filepath.Join(targetDir, "plugin.json"), data, 0644); err != nil {
		return "", err
	}

	// Create subdirectories based on type
	switch PluginType(pluginType) {
	case PluginTypeAgent:
		os.MkdirAll(filepath.Join(targetDir, "agents"), 0755)
		// Create sample agent
		sample := `---
name: my-agent
description: Custom agent description
---

# My Agent

Your agent instructions here.
`
		os.WriteFile(filepath.Join(targetDir, "agents", name+".md"), []byte(sample), 0644)

	case PluginTypeSkill:
		os.MkdirAll(filepath.Join(targetDir, "skills", name), 0755)
		sample := `---
name: my-skill
description: Custom skill description
---

# My Skill

Your skill instructions here.
`
		os.WriteFile(filepath.Join(targetDir, "skills", name, "SKILL.md"), []byte(sample), 0644)

	case PluginTypeBundle:
		os.MkdirAll(filepath.Join(targetDir, "agents"), 0755)
		os.MkdirAll(filepath.Join(targetDir, "skills"), 0755)
		os.MkdirAll(filepath.Join(targetDir, "workflows"), 0755)
	}

	// Create README
	readme := fmt.Sprintf("# %s\n\nAGEN plugin.\n\n## Installation\n\n```bash\nagen plugin install ./%s\n```\n", name, name)
	os.WriteFile(filepath.Join(targetDir, "README.md"), []byte(readme), 0644)

	return targetDir, nil
}

// Registry functions

func loadRegistry(pluginDir string) (*Registry, error) {
	path := filepath.Join(pluginDir, "registry.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var registry Registry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, err
	}
	registry.path = path

	return &registry, nil
}

func (r *Registry) save() error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.path, data, 0644)
}
