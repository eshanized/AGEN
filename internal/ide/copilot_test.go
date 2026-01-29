// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Unit tests for GitHub Copilot Workspace adapter

package ide

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopilotWorkspaceAdapter_Name(t *testing.T) {
	adapter := &CopilotWorkspaceAdapter{}
	if got := adapter.Name(); got != "CopilotWorkspace" {
		t.Errorf("CopilotWorkspaceAdapter.Name() = %q, want %q", got, "CopilotWorkspace")
	}
}

func TestCopilotWorkspaceAdapter_GetRulesPath(t *testing.T) {
	adapter := &CopilotWorkspaceAdapter{}
	if got := adapter.GetRulesPath(); got != ".github/copilot-instructions.md" {
		t.Errorf("CopilotWorkspaceAdapter.GetRulesPath() = %q, want %q", got, ".github/copilot-instructions.md")
	}
}

func TestCopilotWorkspaceAdapter_Detect(t *testing.T) {
	adapter := &CopilotWorkspaceAdapter{}

	tests := []struct {
		name     string
		setup    func(dir string) error
		expected bool
	}{
		{
			name: "detect .github/copilot-instructions.md file",
			setup: func(dir string) error {
				githubDir := filepath.Join(dir, ".github")
				if err := os.MkdirAll(githubDir, 0755); err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(githubDir, "copilot-instructions.md"), []byte("# instructions"), 0644)
			},
			expected: true,
		},
		{
			name:     "no copilot config",
			setup:    func(dir string) error { return nil },
			expected: false,
		},
		{
			name: ".github exists but no copilot-instructions.md",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, ".github"), 0755)
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			if err := tt.setup(tmpDir); err != nil {
				t.Fatal(err)
			}
			if got := adapter.Detect(tmpDir); got != tt.expected {
				t.Errorf("CopilotWorkspaceAdapter.Detect() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCopilotWorkspaceAdapter_Install(t *testing.T) {
	adapter := &CopilotWorkspaceAdapter{}
	tmpl := createMockTemplates()

	t.Run("install creates copilot-instructions.md", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		if err := adapter.Install(tmpl, opts); err != nil {
			t.Fatalf("Install() error = %v", err)
		}

		instructionsFile := filepath.Join(tmpDir, ".github", "copilot-instructions.md")
		if _, err := os.Stat(instructionsFile); os.IsNotExist(err) {
			t.Error(".github/copilot-instructions.md was not created")
		}
	})

	t.Run("install creates .github directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		if err := adapter.Install(tmpl, opts); err != nil {
			t.Fatalf("Install() error = %v", err)
		}

		githubDir := filepath.Join(tmpDir, ".github")
		if info, err := os.Stat(githubDir); err != nil || !info.IsDir() {
			t.Error(".github directory was not created")
		}
	})

	t.Run("install errors without force when exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		githubDir := filepath.Join(tmpDir, ".github")
		os.MkdirAll(githubDir, 0755)
		os.WriteFile(filepath.Join(githubDir, "copilot-instructions.md"), []byte("existing"), 0644)

		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		if err := adapter.Install(tmpl, opts); err == nil {
			t.Error("Install() should error when copilot-instructions.md exists without --force")
		}
	})

	t.Run("install succeeds with force", func(t *testing.T) {
		tmpDir := t.TempDir()
		githubDir := filepath.Join(tmpDir, ".github")
		os.MkdirAll(githubDir, 0755)
		os.WriteFile(filepath.Join(githubDir, "copilot-instructions.md"), []byte("existing"), 0644)

		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     true,
		}

		if err := adapter.Install(tmpl, opts); err != nil {
			t.Errorf("Install() with --force error = %v", err)
		}
	})

	t.Run("dry run does not create files", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts := InstallOptions{
			TargetDir: tmpDir,
			DryRun:    true,
		}

		if err := adapter.Install(tmpl, opts); err != nil {
			t.Fatalf("Install() dry run error = %v", err)
		}

		instructionsFile := filepath.Join(tmpDir, ".github", "copilot-instructions.md")
		if _, err := os.Stat(instructionsFile); !os.IsNotExist(err) {
			t.Error("dry run should not create files")
		}
	})
}

func TestCopilotWorkspaceAdapter_Update(t *testing.T) {
	adapter := &CopilotWorkspaceAdapter{}
	tmpl := createMockTemplates()

	t.Run("update adds new files", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts := UpdateOptions{
			TargetDir: tmpDir,
			Force:     true,
		}

		changes, err := adapter.Update(tmpl, opts)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		if len(changes.Added) == 0 {
			t.Error("Update() should report added files")
		}
	})

	t.Run("update skips existing without force", func(t *testing.T) {
		tmpDir := t.TempDir()
		githubDir := filepath.Join(tmpDir, ".github")
		os.MkdirAll(githubDir, 0755)
		os.WriteFile(filepath.Join(githubDir, "copilot-instructions.md"), []byte("existing"), 0644)

		opts := UpdateOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		changes, err := adapter.Update(tmpl, opts)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		if len(changes.Skipped) == 0 {
			t.Error("Update() should report skipped files")
		}
	})

	t.Run("update force overwrites existing", func(t *testing.T) {
		tmpDir := t.TempDir()
		githubDir := filepath.Join(tmpDir, ".github")
		os.MkdirAll(githubDir, 0755)
		os.WriteFile(filepath.Join(githubDir, "copilot-instructions.md"), []byte("existing"), 0644)

		opts := UpdateOptions{
			TargetDir: tmpDir,
			Force:     true,
		}

		changes, err := adapter.Update(tmpl, opts)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		if len(changes.Updated) == 0 {
			t.Error("Update() with force should report updated files")
		}
	})
}

func TestCopilotWorkspaceAdapter_BuildContent(t *testing.T) {
	adapter := &CopilotWorkspaceAdapter{}
	tmpl := createMockTemplates()

	content := adapter.buildContent(tmpl)

	if !strings.Contains(content, "Copilot Instructions") {
		t.Error("Content should contain 'Copilot Instructions' header")
	}
	if !strings.Contains(content, "AGEN") {
		t.Error("Content should mention AGEN")
	}
	if !strings.Contains(content, "Available Personas") {
		t.Error("Content should list available personas")
	}
	if !strings.Contains(content, "Workflow Commands") {
		t.Error("Content should list workflow commands")
	}
}
