// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Unit tests for verification runner

package verify

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSecurityScan(t *testing.T) {
	// Create temp project with test files
	tmpDir := t.TempDir()

	// Create a file with a potential secret
	jsContent := `
const apiKey = "ghp_1234567890abcdefghij1234567890abcdef";
console.log(apiKey);
`
	jsFile := filepath.Join(tmpDir, "test.js")
	if err := os.WriteFile(jsFile, []byte(jsContent), 0644); err != nil {
		t.Fatal(err)
	}

	runner := NewRunner(tmpDir, RunnerOptions{})
	result := runner.RunSecurity()

	if result.Name != "Security Scan" {
		t.Errorf("Expected name 'Security Scan', got %q", result.Name)
	}

	// Should detect the GitHub token pattern
	if len(result.Issues) == 0 {
		t.Error("Should have detected at least one issue (GitHub token)")
	}

	foundGitHubToken := false
	for _, issue := range result.Issues {
		if issue.Rule == "security/no-secrets" {
			foundGitHubToken = true
		}
	}

	if !foundGitHubToken {
		t.Error("Should have detected GitHub token as secret")
	}
}

func TestSecurityScanClean(t *testing.T) {
	// Create temp project with clean files
	tmpDir := t.TempDir()

	cleanContent := `
function hello() {
    return "Hello, World!";
}
`
	jsFile := filepath.Join(tmpDir, "clean.js")
	if err := os.WriteFile(jsFile, []byte(cleanContent), 0644); err != nil {
		t.Fatal(err)
	}

	runner := NewRunner(tmpDir, RunnerOptions{})
	result := runner.RunSecurity()

	if !result.Passed {
		t.Errorf("Clean project should pass security scan, got %d issues", len(result.Issues))
	}
}

func TestLintCheck(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file with console.log
	jsContent := `
console.log("debug message");
function test() { return true; }
`
	jsFile := filepath.Join(tmpDir, "app.js")
	if err := os.WriteFile(jsFile, []byte(jsContent), 0644); err != nil {
		t.Fatal(err)
	}

	runner := NewRunner(tmpDir, RunnerOptions{})
	result := runner.RunLint()

	if result.Name != "Lint Check" {
		t.Errorf("Expected name 'Lint Check', got %q", result.Name)
	}

	// Should warn about console.log
	foundConsoleLog := false
	for _, issue := range result.Issues {
		if issue.Rule == "lint/no-console" {
			foundConsoleLog = true
		}
	}

	if !foundConsoleLog {
		t.Error("Should have warned about console.log usage")
	}
}

func TestUXAudit(t *testing.T) {
	tmpDir := t.TempDir()

	// Create HTML with accessibility issues
	htmlContent := `
<!DOCTYPE html>
<html>
<body>
    <img src="photo.jpg">
    <input type="text">
</body>
</html>
`
	htmlFile := filepath.Join(tmpDir, "index.html")
	if err := os.WriteFile(htmlFile, []byte(htmlContent), 0644); err != nil {
		t.Fatal(err)
	}

	runner := NewRunner(tmpDir, RunnerOptions{})
	result := runner.RunUX()

	if result.Name != "UX Audit" {
		t.Errorf("Expected name 'UX Audit', got %q", result.Name)
	}

	// Should detect missing alt on img
	foundImgAlt := false
	for _, issue := range result.Issues {
		if issue.Rule == "ux/img-alt" {
			foundImgAlt = true
		}
	}

	if !foundImgAlt {
		t.Error("Should have detected missing alt attribute on image")
	}
}

func TestSEOCheck(t *testing.T) {
	tmpDir := t.TempDir()

	// Create HTML without title and meta description
	htmlContent := `
<!DOCTYPE html>
<html>
<head></head>
<body>
    <p>Content without proper SEO</p>
</body>
</html>
`
	htmlFile := filepath.Join(tmpDir, "page.html")
	if err := os.WriteFile(htmlFile, []byte(htmlContent), 0644); err != nil {
		t.Fatal(err)
	}

	runner := NewRunner(tmpDir, RunnerOptions{})
	result := runner.RunSEO()

	if result.Name != "SEO Check" {
		t.Errorf("Expected name 'SEO Check', got %q", result.Name)
	}

	// Should detect missing title and meta description
	hasTitle := false
	hasDescription := false
	for _, issue := range result.Issues {
		if issue.Rule == "seo/title" {
			hasTitle = true
		}
		if issue.Rule == "seo/meta-description" {
			hasDescription = true
		}
	}

	if !hasTitle {
		t.Error("Should have detected missing title tag")
	}
	if !hasDescription {
		t.Error("Should have detected missing meta description")
	}
}

func TestResultStruct(t *testing.T) {
	result := Result{
		Name:          "Test",
		Passed:        false,
		HasCritical:   true,
		CriticalCount: 2,
		WarningCount:  3,
		Issues: []Issue{
			{Severity: "critical", Message: "test"},
		},
	}

	if result.Passed {
		t.Error("Result.Passed should be false")
	}

	if !result.HasCritical {
		t.Error("Result.HasCritical should be true")
	}

	if len(result.Issues) != 1 {
		t.Errorf("Expected 1 issue, got %d", len(result.Issues))
	}
}
