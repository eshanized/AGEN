// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Unit tests for Antigravity adapter

package ide

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAntigravityAdapter_Name(t *testing.T) {
	adapter := &AntigravityAdapter{}
	if got := adapter.Name(); got != "Antigravity" {
		t.Errorf("AntigravityAdapter.Name() = %q, want %q", got, "Antigravity")
	}
}

func TestAntigravityAdapter_GetRulesPath(t *testing.T) {
	adapter := &AntigravityAdapter{}
	if got := adapter.GetRulesPath(); got != ".agent/rules/GEMINI.md" {
		t.Errorf("AntigravityAdapter.GetRulesPath() = %q, want %q", got, ".agent/rules/GEMINI.md")
	}
}

func TestAntigravityAdapter_Detect(t *testing.T) {
	adapter := &AntigravityAdapter{}

	tests := []struct {
		name     string
		setup    func(dir string) error
		expected bool
	}{
		{
			name: "detect .agent directory",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, ".agent"), 0755)
			},
			expected: true,
		},
		{
			name: "detect GEMINI.md file",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "GEMINI.md"), []byte("# gemini"), 0644)
			},
			expected: true,
		},
		{
			name:     "no antigravity config",
			setup:    func(dir string) error { return nil },
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
				t.Errorf("AntigravityAdapter.Detect() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAntigravityAdapter_Install(t *testing.T) {
	adapter := &AntigravityAdapter{}
	tmpl := createMockTemplates()

	t.Run("install creates .agent directory structure", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		if err := adapter.Install(tmpl, opts); err != nil {
			t.Fatalf("Install() error = %v", err)
		}

		agentDir := filepath.Join(tmpDir, ".agent")
		if _, err := os.Stat(agentDir); os.IsNotExist(err) {
			t.Error(".agent directory was not created")
		}

		// Check subdirectories
		subdirs := []string{"agents", "skills", "workflows", "rules"}
		for _, subdir := range subdirs {
			path := filepath.Join(agentDir, subdir)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf(".agent/%s directory was not created", subdir)
			}
		}
	})

	t.Run("install errors without force when exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		agentDir := filepath.Join(tmpDir, ".agent")
		os.MkdirAll(agentDir, 0755)

		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		if err := adapter.Install(tmpl, opts); err == nil {
			t.Error("Install() should error when .agent exists without --force")
		}
	})

	t.Run("install succeeds with force", func(t *testing.T) {
		tmpDir := t.TempDir()
		agentDir := filepath.Join(tmpDir, ".agent")
		os.MkdirAll(agentDir, 0755)

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

		agentDir := filepath.Join(tmpDir, ".agent")
		if _, err := os.Stat(agentDir); !os.IsNotExist(err) {
			t.Error("dry run should not create files")
		}
	})
}

func TestAntigravityAdapter_Update(t *testing.T) {
	adapter := &AntigravityAdapter{}
	tmpl := createMockTemplates()

	t.Run("update adds new files", func(t *testing.T) {
		tmpDir := t.TempDir()
		// Create .agent structure first
		agentDir := filepath.Join(tmpDir, ".agent", "agents")
		os.MkdirAll(agentDir, 0755)
		os.MkdirAll(filepath.Join(tmpDir, ".agent", "skills"), 0755)

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

	t.Run("update returns changes object", func(t *testing.T) {
		tmpDir := t.TempDir()

		opts := UpdateOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		changes, err := adapter.Update(tmpl, opts)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		if changes == nil {
			t.Error("Update() should return changes object")
		}
	})

	t.Run("update with force option", func(t *testing.T) {
		tmpDir := t.TempDir()
		// Create .agent structure with existing agent
		agentDir := filepath.Join(tmpDir, ".agent", "agents")
		os.MkdirAll(agentDir, 0755)
		// Create an agent with different content
		os.WriteFile(filepath.Join(agentDir, "test-agent.md"), []byte("different content"), 0644)

		opts := UpdateOptions{
			TargetDir: tmpDir,
			Force:     true,
		}

		changes, err := adapter.Update(tmpl, opts)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		// Should have either updated or added files
		totalChanges := len(changes.Added) + len(changes.Updated)
		if totalChanges == 0 {
			t.Error("Update() with force should report changed files")
		}
	})
}

func TestAntigravityAdapter_InstallCreatesTemplateStructure(t *testing.T) {
	adapter := &AntigravityAdapter{}
	tmpl := createMockTemplates()

	tmpDir := t.TempDir()
	opts := InstallOptions{
		TargetDir: tmpDir,
		Force:     true,
	}

	if err := adapter.Install(tmpl, opts); err != nil {
		t.Fatalf("Install() error = %v", err)
	}

	// Verify agents were installed
	agentsDir := filepath.Join(tmpDir, ".agent", "agents")
	entries, err := os.ReadDir(agentsDir)
	if err != nil {
		t.Fatalf("Failed to read agents directory: %v", err)
	}
	if len(entries) == 0 {
		t.Error("No agents were installed")
	}
}
