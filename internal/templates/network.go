// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// A cross-platform CLI tool for managing AI agent templates

package templates

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GitHubRepo holds the repository info for fetching templates
const (
	defaultOwner  = "eshanized"
	defaultRepo   = "agen"
	defaultBranch = "main"
)

// GitHubContentsResponse represents the GitHub API response for directory contents
type GitHubContentsResponse struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"` // "file" or "dir"
	DownloadURL string `json:"download_url,omitempty"`
}

// FetchFromGitHub downloads templates from the GitHub repository.
//
// How it works:
// 1. Try to download as a ZIP archive (faster for full download)
// 2. If that fails, fall back to API-based file-by-file download
// 3. Parse and return as Templates struct
//
// Why ZIP first? Downloading the whole repo as ZIP is faster than
// making 50+ individual API requests for each file. GitHub has
// rate limits, so we want to be efficient.
func FetchFromGitHub(branch string) (*Templates, error) {
	if branch == "" {
		branch = defaultBranch
	}

	// Try ZIP download first
	tmpl, err := fetchViaZip(branch)
	if err == nil {
		return tmpl, nil
	}

	// Fall back to API
	return fetchViaAPI(branch)
}

// fetchViaZip downloads the repo as a ZIP and extracts templates.
//
// GitHub provides ZIP downloads at:
// https://github.com/{owner}/{repo}/archive/{branch}.zip
func fetchViaZip(branch string) (*Templates, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	zipURL := fmt.Sprintf("https://github.com/%s/%s/archive/%s.zip",
		defaultOwner, defaultRepo, branch)

	req, err := http.NewRequestWithContext(ctx, "GET", zipURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	// Save to temp file
	tmpFile, err := os.CreateTemp("", "agen-templates-*.zip")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return nil, err
	}

	// Extract and parse
	return extractTemplatesFromZip(tmpFile.Name(), branch)
}

// extractTemplatesFromZip reads templates from a downloaded ZIP file.
func extractTemplatesFromZip(zipPath, branch string) (*Templates, error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	tmpl := &Templates{
		Version:   CurrentVersion,
		Agents:    make(map[string]Agent),
		Skills:    make(map[string]Skill),
		Workflows: make(map[string]Workflow),
	}

	// The ZIP contains a top-level directory like "agen-main/"
	// We need to look for templates/data/ inside that
	prefix := fmt.Sprintf("%s-%s/internal/templates/data/", defaultRepo, branch)

	for _, file := range reader.File {
		// skip directories
		if file.FileInfo().IsDir() {
			continue
		}

		// Check if this is a template file
		if !strings.HasPrefix(file.Name, prefix) {
			continue
		}

		relativePath := strings.TrimPrefix(file.Name, prefix)

		// Read file content
		rc, err := file.Open()
		if err != nil {
			continue
		}
		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			continue
		}

		// Parse based on path
		parts := strings.Split(relativePath, "/")
		if len(parts) < 2 {
			continue
		}

		switch parts[0] {
		case "agents":
			if strings.HasSuffix(parts[1], ".md") {
				name := strings.TrimSuffix(parts[1], ".md")
				agent := parseAgentFile(string(content))
				agent.Name = name
				tmpl.Agents[name] = agent
			}

		case "skills":
			// skills/skill-name/SKILL.md
			if len(parts) >= 2 && parts[len(parts)-1] == "SKILL.md" {
				name := parts[1]
				skill := parseSkillFile(string(content))
				skill.Name = name
				tmpl.Skills[name] = skill
			}

		case "workflows":
			if strings.HasSuffix(parts[1], ".md") {
				name := strings.TrimSuffix(parts[1], ".md")
				workflow := parseWorkflowFile(string(content))
				workflow.Name = name
				tmpl.Workflows[name] = workflow
			}
		}
	}

	return tmpl, nil
}

// fetchViaAPI uses GitHub's API to download files individually.
// this is slower but works if ZIP download fails.
func fetchViaAPI(branch string) (*Templates, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	tmpl := &Templates{
		Version:   CurrentVersion,
		Agents:    make(map[string]Agent),
		Skills:    make(map[string]Skill),
		Workflows: make(map[string]Workflow),
	}

	// Fetch agents
	agentFiles, err := listGitHubDir(ctx, "internal/templates/data/agents", branch)
	if err == nil {
		for _, file := range agentFiles {
			if file.Type == "file" && strings.HasSuffix(file.Name, ".md") {
				content, err := downloadFile(ctx, file.DownloadURL)
				if err != nil {
					continue
				}
				name := strings.TrimSuffix(file.Name, ".md")
				agent := parseAgentFile(content)
				agent.Name = name
				tmpl.Agents[name] = agent
			}
		}
	}

	// Fetch workflows
	workflowFiles, err := listGitHubDir(ctx, "internal/templates/data/workflows", branch)
	if err == nil {
		for _, file := range workflowFiles {
			if file.Type == "file" && strings.HasSuffix(file.Name, ".md") {
				content, err := downloadFile(ctx, file.DownloadURL)
				if err != nil {
					continue
				}
				name := strings.TrimSuffix(file.Name, ".md")
				workflow := parseWorkflowFile(content)
				workflow.Name = name
				tmpl.Workflows[name] = workflow
			}
		}
	}

	// Fetch skills (need to list directories first)
	skillDirs, err := listGitHubDir(ctx, "internal/templates/data/skills", branch)
	if err == nil {
		for _, dir := range skillDirs {
			if dir.Type == "dir" {
				// Get SKILL.md from this directory
				skillFiles, err := listGitHubDir(ctx, dir.Path, branch)
				if err != nil {
					continue
				}
				for _, file := range skillFiles {
					if file.Name == "SKILL.md" {
						content, err := downloadFile(ctx, file.DownloadURL)
						if err != nil {
							continue
						}
						skill := parseSkillFile(content)
						skill.Name = dir.Name
						tmpl.Skills[dir.Name] = skill
						break
					}
				}
			}
		}
	}

	return tmpl, nil
}

// listGitHubDir lists contents of a directory via GitHub API
func listGitHubDir(ctx context.Context, path, branch string) ([]GitHubContentsResponse, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s",
		defaultOwner, defaultRepo, path, branch)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "agen-cli")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var contents []GitHubContentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&contents); err != nil {
		return nil, err
	}

	return contents, nil
}

// downloadFile downloads a single file from a URL
func downloadFile(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// CacheTemplates saves fetched templates to local cache for offline use.
func CacheTemplates(tmpl *Templates, cacheDir string) error {
	templatesDir := filepath.Join(cacheDir, "templates")

	// create directories
	dirs := []string{
		filepath.Join(templatesDir, "agents"),
		filepath.Join(templatesDir, "skills"),
		filepath.Join(templatesDir, "workflows"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Write agents
	for name, agent := range tmpl.Agents {
		file := filepath.Join(templatesDir, "agents", name+".md")
		if err := os.WriteFile(file, []byte(agent.Content), 0644); err != nil {
			return err
		}
	}

	// Write skills
	for name, skill := range tmpl.Skills {
		skillDir := filepath.Join(templatesDir, "skills", name)
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			return err
		}
		file := filepath.Join(skillDir, "SKILL.md")
		if err := os.WriteFile(file, []byte(skill.Content), 0644); err != nil {
			return err
		}
	}

	// Write workflows
	for name, workflow := range tmpl.Workflows {
		file := filepath.Join(templatesDir, "workflows", name+".md")
		if err := os.WriteFile(file, []byte(workflow.Content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// LoadFromCache loads templates from local cache.
func LoadFromCache(cacheDir string) (*Templates, error) {
	templatesDir := filepath.Join(cacheDir, "templates")

	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("no cached templates found")
	}

	tmpl := &Templates{
		Version:   CurrentVersion,
		Agents:    make(map[string]Agent),
		Skills:    make(map[string]Skill),
		Workflows: make(map[string]Workflow),
	}

	// Load agents
	agentsDir := filepath.Join(templatesDir, "agents")
	if entries, err := os.ReadDir(agentsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				content, err := os.ReadFile(filepath.Join(agentsDir, entry.Name()))
				if err != nil {
					continue
				}
				name := strings.TrimSuffix(entry.Name(), ".md")
				agent := parseAgentFile(string(content))
				agent.Name = name
				tmpl.Agents[name] = agent
			}
		}
	}

	// Load skills
	skillsDir := filepath.Join(templatesDir, "skills")
	if entries, err := os.ReadDir(skillsDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				skillFile := filepath.Join(skillsDir, entry.Name(), "SKILL.md")
				content, err := os.ReadFile(skillFile)
				if err != nil {
					continue
				}
				skill := parseSkillFile(string(content))
				skill.Name = entry.Name()
				tmpl.Skills[entry.Name()] = skill
			}
		}
	}

	// Load workflows
	workflowsDir := filepath.Join(templatesDir, "workflows")
	if entries, err := os.ReadDir(workflowsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				content, err := os.ReadFile(filepath.Join(workflowsDir, entry.Name()))
				if err != nil {
					continue
				}
				name := strings.TrimSuffix(entry.Name(), ".md")
				workflow := parseWorkflowFile(string(content))
				workflow.Name = name
				tmpl.Workflows[name] = workflow
			}
		}
	}

	return tmpl, nil
}
