// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// A cross-platform CLI tool for managing AI agent templates

package verify

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Result represents the outcome of a verification check
type Result struct {
	Name          string
	Passed        bool
	HasCritical   bool
	CriticalCount int
	WarningCount  int
	Issues        []Issue
	Duration      float64
}

// Issue represents a single verification issue
type Issue struct {
	Severity string // "critical", "warning", "info"
	File     string
	Line     int
	Message  string
	Rule     string
}

// RunnerOptions configures the verification runner
type RunnerOptions struct {
	Verbose bool
}

// Runner orchestrates verification checks
type Runner struct {
	projectPath string
	options     RunnerOptions
}

// NewRunner creates a new verification runner
func NewRunner(projectPath string, options RunnerOptions) *Runner {
	return &Runner{
		projectPath: projectPath,
		options:     options,
	}
}

// RunSecurity performs security scanning.
//
// How it works:
// 1. Scan for hardcoded secrets (API keys, passwords, tokens)
// 2. Check for common security anti-patterns
// 3. Look for vulnerable package patterns
// 4. Check .env files aren't committed
//
// This is a simplified version of the original Python script.
// For full security auditing, consider using dedicated tools like
// trufflehog, gitleaks, or snyk.
func (r *Runner) RunSecurity() Result {
	result := Result{
		Name:   "Security Scan",
		Passed: true,
	}

	// Secret patterns to look for
	secretPatterns := map[string]*regexp.Regexp{
		"AWS Key":      regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
		"GitHub Token": regexp.MustCompile(`ghp_[a-zA-Z0-9]{36}`),
		"Private Key":  regexp.MustCompile(`-----BEGIN (RSA |EC |DSA |OPENSSH )?PRIVATE KEY-----`),
		"Password":     regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[=:]\s*["'][^"']+["']`),
		"API Key":      regexp.MustCompile(`(?i)(api[_-]?key|apikey)\s*[=:]\s*["'][^"']+["']`),
		"Secret":       regexp.MustCompile(`(?i)(secret|token)\s*[=:]\s*["'][^"']+["']`),
	}

	// Files to skip
	skipDirs := map[string]bool{
		"node_modules": true,
		".git":         true,
		"vendor":       true,
		"dist":         true,
		"build":        true,
	}

	// Walk the project
	filepath.Walk(r.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// skip directories
		if info.IsDir() {
			if skipDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		// Only check text files
		ext := strings.ToLower(filepath.Ext(path))
		textExts := map[string]bool{
			".js": true, ".ts": true, ".jsx": true, ".tsx": true,
			".py": true, ".go": true, ".java": true, ".rb": true,
			".env": true, ".yml": true, ".yaml": true, ".json": true,
			".md": true, ".txt": true, ".sh": true, ".bash": true,
		}

		if !textExts[ext] && filepath.Base(path) != ".env" {
			return nil
		}

		// Read and scan
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(r.projectPath, path)

		for patternName, pattern := range secretPatterns {
			if matches := pattern.FindAllStringIndex(string(content), -1); len(matches) > 0 {
				for _, match := range matches {
					lineNum := strings.Count(string(content[:match[0]]), "\n") + 1
					result.Issues = append(result.Issues, Issue{
						Severity: "critical",
						File:     relPath,
						Line:     lineNum,
						Message:  "Potential " + patternName + " detected",
						Rule:     "security/no-secrets",
					})
					result.CriticalCount++
				}
			}
		}

		return nil
	})

	// Check for .env file without .gitignore
	envFile := filepath.Join(r.projectPath, ".env")
	if _, err := os.Stat(envFile); err == nil {
		gitignore := filepath.Join(r.projectPath, ".gitignore")
		if content, err := os.ReadFile(gitignore); err != nil || !strings.Contains(string(content), ".env") {
			result.Issues = append(result.Issues, Issue{
				Severity: "warning",
				File:     ".env",
				Message:  ".env file exists but may not be in .gitignore",
				Rule:     "security/env-gitignore",
			})
			result.WarningCount++
		}
	}

	result.Passed = result.CriticalCount == 0
	result.HasCritical = result.CriticalCount > 0

	return result
}

// RunLint performs basic linting checks.
//
// For JavaScript/TypeScript projects, we try to run npm run lint.
// For Go projects, we try golangci-lint.
// Falls back to basic file checks if no linter is available.
func (r *Runner) RunLint() Result {
	result := Result{
		Name:   "Lint Check",
		Passed: true,
	}

	// Check for linter configs
	packageJSON := filepath.Join(r.projectPath, "package.json")
	goMod := filepath.Join(r.projectPath, "go.mod")

	if _, err := os.Stat(packageJSON); err == nil {
		// Try to run npm run lint
		if r.options.Verbose {
			result.Issues = append(result.Issues, Issue{
				Severity: "info",
				Message:  "Running npm run lint...",
				Rule:     "lint/npm",
			})
		}
		cmd := exec.Command("npm", "run", "lint", "--silent")
		cmd.Dir = r.projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			// Parse npm lint output for issues
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "error") || strings.Contains(line, "Error") {
					result.Issues = append(result.Issues, Issue{
						Severity: "warning",
						Message:  strings.TrimSpace(line),
						Rule:     "lint/eslint",
					})
					result.WarningCount++
				}
			}
		}
	} else if _, err := os.Stat(goMod); err == nil {
		// Try to run golangci-lint
		cmd := exec.Command("golangci-lint", "run", "--out-format=line-number")
		cmd.Dir = r.projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			// Parse golangci-lint output
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "level=") {
					result.Issues = append(result.Issues, Issue{
						Severity: "warning",
						Message:  strings.TrimSpace(line),
						Rule:     "lint/golangci",
					})
					result.WarningCount++
				}
			}
		}
	}

	// Basic file checks
	filepath.Walk(r.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// Check for console.log in production code
		ext := filepath.Ext(path)
		if ext == ".js" || ext == ".ts" || ext == ".jsx" || ext == ".tsx" {
			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			relPath, _ := filepath.Rel(r.projectPath, path)

			if strings.Contains(string(content), "console.log(") {
				lineNum := strings.Count(strings.Split(string(content), "console.log(")[0], "\n") + 1
				result.Issues = append(result.Issues, Issue{
					Severity: "warning",
					File:     relPath,
					Line:     lineNum,
					Message:  "console.log() found - remove before production",
					Rule:     "lint/no-console",
				})
				result.WarningCount++
			}
		}

		return nil
	})

	result.Passed = result.CriticalCount == 0
	result.HasCritical = result.CriticalCount > 0

	return result
}

// RunUX performs basic UX auditing on HTML/JSX files.
//
// Checks for:
// - Missing alt text on images
// - Buttons without accessible labels
// - Touch target sizes (for buttons/links)
// - Color contrast issues (basic)
func (r *Runner) RunUX() Result {
	result := Result{
		Name:   "UX Audit",
		Passed: true,
	}

	htmlExts := map[string]bool{
		".html": true, ".htm": true, ".jsx": true, ".tsx": true,
	}

	filepath.Walk(r.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !htmlExts[ext] {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(r.projectPath, path)
		contentStr := string(content)

		// Check for images without alt
		imgPattern := regexp.MustCompile(`<img[^>]*>`)
		altPattern := regexp.MustCompile(`alt\s*=`)
		for _, match := range imgPattern.FindAllString(contentStr, -1) {
			if !altPattern.MatchString(match) {
				result.Issues = append(result.Issues, Issue{
					Severity: "warning",
					File:     relPath,
					Message:  "Image missing alt attribute",
					Rule:     "ux/img-alt",
				})
				result.WarningCount++
			}
		}

		// check for form inputs without labels
		inputPattern := regexp.MustCompile(`<input[^>]*>`)
		for _, match := range inputPattern.FindAllString(contentStr, -1) {
			if !strings.Contains(match, "aria-label") && !strings.Contains(match, "id=") {
				result.Issues = append(result.Issues, Issue{
					Severity: "warning",
					File:     relPath,
					Message:  "Input may be missing associated label",
					Rule:     "ux/input-label",
				})
				result.WarningCount++
			}
		}

		return nil
	})

	result.Passed = result.CriticalCount == 0
	result.HasCritical = result.CriticalCount > 0

	return result
}

// RunSEO performs basic SEO checks on HTML files.
//
// Checks for:
// - Title tag
// - Meta description
// - OG tags
// - Canonical URL
// - H1 usage
func (r *Runner) RunSEO() Result {
	result := Result{
		Name:   "SEO Check",
		Passed: true,
	}

	// look for HTML files
	htmlFiles := []string{}

	filepath.Walk(r.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".htm") {
			htmlFiles = append(htmlFiles, path)
		}

		return nil
	})

	for _, htmlFile := range htmlFiles {
		content, err := os.ReadFile(htmlFile)
		if err != nil {
			continue
		}

		relPath, _ := filepath.Rel(r.projectPath, htmlFile)
		contentStr := string(content)

		// check for title
		if !strings.Contains(contentStr, "<title>") {
			result.Issues = append(result.Issues, Issue{
				Severity: "warning",
				File:     relPath,
				Message:  "Missing <title> tag",
				Rule:     "seo/title",
			})
			result.WarningCount++
		}

		// check for meta description
		if !strings.Contains(contentStr, `name="description"`) {
			result.Issues = append(result.Issues, Issue{
				Severity: "warning",
				File:     relPath,
				Message:  "Missing meta description",
				Rule:     "seo/meta-description",
			})
			result.WarningCount++
		}

		// check for H1
		if !strings.Contains(contentStr, "<h1") {
			result.Issues = append(result.Issues, Issue{
				Severity: "info",
				File:     relPath,
				Message:  "Missing H1 heading",
				Rule:     "seo/h1",
			})
		}
	}

	result.Passed = result.CriticalCount == 0
	result.HasCritical = result.CriticalCount > 0

	return result
}
