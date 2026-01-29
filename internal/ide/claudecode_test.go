// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Unit tests for Claude Code adapter

package ide

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestClaudeCodeAdapter_Name(t *testing.T) {
	adapter := &ClaudeCodeAdapter{}
	if got := adapter.Name(); got != "ClaudeCode" {
		t.Errorf("ClaudeCodeAdapter.Name() = %q, want %q", got, "ClaudeCode")
	}
}

func TestClaudeCodeAdapter_GetRulesPath(t *testing.T) {
	adapter := &ClaudeCodeAdapter{}
	if got := adapter.GetRulesPath(); got != "CLAUDE.md" {
		t.Errorf("ClaudeCodeAdapter.GetRulesPath() = %q, want %q", got, "CLAUDE.md")
	}
}

func TestClaudeCodeAdapter_Detect(t *testing.T) {
	adapter := &ClaudeCodeAdapter{}

	tests := []struct {
		name     string
		setup    func(dir string) error
		expected bool
	}{
		{
			name: "detect CLAUDE.md file (uppercase)",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# claude"), 0644)
			},
			expected: true,
		},
		{
			name: "detect claude.md file (lowercase)",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "claude.md"), []byte("# claude"), 0644)
			},
			expected: true,
		},
		{
			name:     "no claude config",
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
				t.Errorf("ClaudeCodeAdapter.Detect() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestClaudeCodeAdapter_Install(t *testing.T) {
	adapter := &ClaudeCodeAdapter{}
	tmpl := createMockTemplates()

	t.Run("install creates CLAUDE.md", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		if err := adapter.Install(tmpl, opts); err != nil {
			t.Fatalf("Install() error = %v", err)
		}

		claudeFile := filepath.Join(tmpDir, "CLAUDE.md")
		if _, err := os.Stat(claudeFile); os.IsNotExist(err) {
			t.Error("CLAUDE.md was not created")
		}
	})

	t.Run("install errors without force when exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.WriteFile(filepath.Join(tmpDir, "CLAUDE.md"), []byte("existing"), 0644)

		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		if err := adapter.Install(tmpl, opts); err == nil {
			t.Error("Install() should error when CLAUDE.md exists without --force")
		}
	})

	t.Run("install succeeds with force", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.WriteFile(filepath.Join(tmpDir, "CLAUDE.md"), []byte("existing"), 0644)

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

		claudeFile := filepath.Join(tmpDir, "CLAUDE.md")
		if _, err := os.Stat(claudeFile); !os.IsNotExist(err) {
			t.Error("dry run should not create files")
		}
	})
}

func TestClaudeCodeAdapter_Update(t *testing.T) {
	adapter := &ClaudeCodeAdapter{}
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
		os.WriteFile(filepath.Join(tmpDir, "CLAUDE.md"), []byte("existing"), 0644)

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
		os.WriteFile(filepath.Join(tmpDir, "CLAUDE.md"), []byte("existing"), 0644)

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

func TestClaudeCodeAdapter_BuildContent(t *testing.T) {
	adapter := &ClaudeCodeAdapter{}
	tmpl := createMockTemplates()

	content := adapter.buildContent(tmpl)

	if !strings.Contains(content, "CLAUDE.md") {
		t.Error("Content should contain 'CLAUDE.md' header")
	}
	if !strings.Contains(content, "AGEN") {
		t.Error("Content should mention AGEN")
	}
	if !strings.Contains(content, "Available Agents") {
		t.Error("Content should list available agents")
	}
}
