// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Unit tests for Zed adapter

package ide

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eshanized/agen/internal/templates"
)

func TestZedAdapter_Name(t *testing.T) {
	adapter := &ZedAdapter{}
	if got := adapter.Name(); got != "Zed" {
		t.Errorf("ZedAdapter.Name() = %q, want %q", got, "Zed")
	}
}

func TestZedAdapter_GetRulesPath(t *testing.T) {
	adapter := &ZedAdapter{}
	if got := adapter.GetRulesPath(); got != ".zed/settings.json" {
		t.Errorf("ZedAdapter.GetRulesPath() = %q, want %q", got, ".zed/settings.json")
	}
}

func TestZedAdapter_Detect(t *testing.T) {
	adapter := &ZedAdapter{}

	tests := []struct {
		name     string
		setup    func(dir string) error
		expected bool
	}{
		{
			name: "detect .zed directory",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, ".zed"), 0755)
			},
			expected: true,
		},
		{
			name:     "no zed config",
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
				t.Errorf("ZedAdapter.Detect() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestZedAdapter_Install(t *testing.T) {
	adapter := &ZedAdapter{}
	tmpl := createMockTemplates()

	t.Run("install creates .zed directory and files", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		if err := adapter.Install(tmpl, opts); err != nil {
			t.Fatalf("Install() error = %v", err)
		}

		zedDir := filepath.Join(tmpDir, ".zed")
		if _, err := os.Stat(zedDir); os.IsNotExist(err) {
			t.Error(".zed directory was not created")
		}

		settingsFile := filepath.Join(zedDir, "settings.json")
		if _, err := os.Stat(settingsFile); os.IsNotExist(err) {
			t.Error(".zed/settings.json was not created")
		}
	})

	t.Run("install overwrites existing directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		zedDir := filepath.Join(tmpDir, ".zed")
		os.MkdirAll(zedDir, 0755)
		os.WriteFile(filepath.Join(zedDir, "settings.json"), []byte("{}"), 0644)

		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     true,
		}

		if err := adapter.Install(tmpl, opts); err != nil {
			t.Errorf("Install() with force error = %v", err)
		}
	})

	t.Run("install creates prompts directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		if err := adapter.Install(tmpl, opts); err != nil {
			t.Fatalf("Install() error = %v", err)
		}

		promptsDir := filepath.Join(tmpDir, ".zed", "prompts")
		if _, err := os.Stat(promptsDir); os.IsNotExist(err) {
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

		zedDir := filepath.Join(tmpDir, ".zed")
		if _, err := os.Stat(zedDir); !os.IsNotExist(err) {
			t.Error("dry run should not create files")
		}
	})
}

func TestZedAdapter_Update(t *testing.T) {
	adapter := &ZedAdapter{}
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
		zedDir := filepath.Join(tmpDir, ".zed")
		os.MkdirAll(zedDir, 0755)
		os.WriteFile(filepath.Join(zedDir, "settings.json"), []byte("{}"), 0644)

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
		zedDir := filepath.Join(tmpDir, ".zed")
		os.MkdirAll(zedDir, 0755)
		os.WriteFile(filepath.Join(zedDir, "settings.json"), []byte("{}"), 0644)

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

func TestZedAdapter_BuildPromptContent(t *testing.T) {
	adapter := &ZedAdapter{}
	agent := templates.Agent{
		Name:        "test-agent",
		Description: "A test agent",
		Content:     "Test content",
	}

	content := adapter.buildPromptContent("test-agent", agent)

	if !strings.Contains(content, "test-agent") {
		t.Error("Content should contain agent name")
	}
}

func TestZedAdapter_BuildMainPrompt(t *testing.T) {
	adapter := &ZedAdapter{}
	tmpl := createMockTemplates()

	content := adapter.buildMainPrompt(tmpl)

	if content == "" {
		t.Error("buildMainPrompt() should return non-empty content")
	}
	if !strings.Contains(content, "AGEN") {
		t.Error("Content should mention AGEN")
	}
}
