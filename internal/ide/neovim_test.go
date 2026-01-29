// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Unit tests for Neovim adapter

package ide

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNeovimAdapter_Name(t *testing.T) {
	adapter := &NeovimAdapter{}
	if got := adapter.Name(); got != "Neovim" {
		t.Errorf("NeovimAdapter.Name() = %q, want %q", got, "Neovim")
	}
}

func TestNeovimAdapter_GetRulesPath(t *testing.T) {
	adapter := &NeovimAdapter{}
	if got := adapter.GetRulesPath(); got != ".nvim/ai-rules.md" {
		t.Errorf("NeovimAdapter.GetRulesPath() = %q, want %q", got, ".nvim/ai-rules.md")
	}
}

func TestNeovimAdapter_Detect(t *testing.T) {
	adapter := &NeovimAdapter{}

	tests := []struct {
		name     string
		setup    func(dir string) error
		expected bool
	}{
		{
			name: "detect .nvim directory",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, ".nvim"), 0755)
			},
			expected: true,
		},
		{
			name: "detect .nvim.lua file",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, ".nvim.lua"), []byte("-- config"), 0644)
			},
			expected: true,
		},
		{
			name: "detect .exrc file",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, ".exrc"), []byte("\" config"), 0644)
			},
			expected: true,
		},
		{
			name:     "no neovim config",
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
				t.Errorf("NeovimAdapter.Detect() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNeovimAdapter_Install(t *testing.T) {
	adapter := &NeovimAdapter{}
	tmpl := createMockTemplates()

	t.Run("install creates files", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		if err := adapter.Install(tmpl, opts); err != nil {
			t.Fatalf("Install() error = %v", err)
		}

		// Check .nvim/ai-rules.md exists
		rulesFile := filepath.Join(tmpDir, ".nvim", "ai-rules.md")
		if _, err := os.Stat(rulesFile); os.IsNotExist(err) {
			t.Error(".nvim/ai-rules.md was not created")
		}

		// Check .nvim.lua exists
		nvimLua := filepath.Join(tmpDir, ".nvim.lua")
		if _, err := os.Stat(nvimLua); os.IsNotExist(err) {
			t.Error(".nvim.lua was not created")
		}
	})

	t.Run("install errors without force when exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.MkdirAll(filepath.Join(tmpDir, ".nvim"), 0755)

		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		if err := adapter.Install(tmpl, opts); err == nil {
			t.Error("Install() should error when directory exists without --force")
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

		nvimDir := filepath.Join(tmpDir, ".nvim")
		if _, err := os.Stat(nvimDir); !os.IsNotExist(err) {
			t.Error("dry run should not create files")
		}
	})
}

func TestNeovimAdapter_Update(t *testing.T) {
	adapter := &NeovimAdapter{}
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
		nvimDir := filepath.Join(tmpDir, ".nvim")
		os.MkdirAll(nvimDir, 0755)
		os.WriteFile(filepath.Join(nvimDir, "ai-rules.md"), []byte("existing"), 0644)

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
}

func TestNeovimAdapter_BuildRulesContent(t *testing.T) {
	adapter := &NeovimAdapter{}
	tmpl := createMockTemplates()

	content := adapter.buildRulesContent(tmpl)

	if !strings.Contains(content, "Neovim AI Rules") {
		t.Error("Content should contain 'Neovim AI Rules' header")
	}
	if !strings.Contains(content, "AGEN") {
		t.Error("Content should mention AGEN")
	}
}

func TestNeovimAdapter_BuildNvimLua(t *testing.T) {
	adapter := &NeovimAdapter{}
	content := adapter.buildNvimLua()

	if !strings.Contains(content, "AGEN") {
		t.Error("Content should mention AGEN")
	}
	if !strings.Contains(content, "ai_rules") {
		t.Error("Content should reference ai_rules")
	}
}
