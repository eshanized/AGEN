// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Unit tests for Aider adapter

package ide

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAiderAdapter_Name(t *testing.T) {
	adapter := &AiderAdapter{}
	if got := adapter.Name(); got != "Aider" {
		t.Errorf("AiderAdapter.Name() = %q, want %q", got, "Aider")
	}
}

func TestAiderAdapter_GetRulesPath(t *testing.T) {
	adapter := &AiderAdapter{}
	if got := adapter.GetRulesPath(); got != ".aider-context.md" {
		t.Errorf("AiderAdapter.GetRulesPath() = %q, want %q", got, ".aider-context.md")
	}
}

func TestAiderAdapter_Detect(t *testing.T) {
	adapter := &AiderAdapter{}

	tests := []struct {
		name     string
		setup    func(dir string) error
		expected bool
	}{
		{
			name: "detect .aider.conf.yml file",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, ".aider.conf.yml"), []byte("# config"), 0644)
			},
			expected: true,
		},
		{
			name: "detect .aider directory",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, ".aider"), 0755)
			},
			expected: true,
		},
		{
			name:     "no aider config",
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
				t.Errorf("AiderAdapter.Detect() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAiderAdapter_Install(t *testing.T) {
	adapter := &AiderAdapter{}
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

		// Check .aider.conf.yml exists
		configFile := filepath.Join(tmpDir, ".aider.conf.yml")
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			t.Error(".aider.conf.yml was not created")
		}

		// Check .aider-context.md exists
		contextFile := filepath.Join(tmpDir, ".aider-context.md")
		if _, err := os.Stat(contextFile); os.IsNotExist(err) {
			t.Error(".aider-context.md was not created")
		}
	})

	t.Run("install errors without force when exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.WriteFile(filepath.Join(tmpDir, ".aider.conf.yml"), []byte("existing"), 0644)

		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		if err := adapter.Install(tmpl, opts); err == nil {
			t.Error("Install() should error when config exists without --force")
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

		configFile := filepath.Join(tmpDir, ".aider.conf.yml")
		if _, err := os.Stat(configFile); !os.IsNotExist(err) {
			t.Error("dry run should not create files")
		}
	})
}

func TestAiderAdapter_Update(t *testing.T) {
	adapter := &AiderAdapter{}
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
		os.WriteFile(filepath.Join(tmpDir, ".aider-context.md"), []byte("existing"), 0644)

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

func TestAiderAdapter_BuildConfigContent(t *testing.T) {
	adapter := &AiderAdapter{}
	content := adapter.buildConfigContent()

	if !strings.Contains(content, "Aider Configuration") {
		t.Error("Config should contain 'Aider Configuration' header")
	}
	if !strings.Contains(content, "auto-commits") {
		t.Error("Config should contain auto-commits setting")
	}
	if !strings.Contains(content, ".aider-context.md") {
		t.Error("Config should reference context file")
	}
}

func TestAiderAdapter_BuildContextContent(t *testing.T) {
	adapter := &AiderAdapter{}
	tmpl := createMockTemplates()

	content := adapter.buildContextContent(tmpl)

	if !strings.Contains(content, "Aider Context") {
		t.Error("Content should contain 'Aider Context' header")
	}
	if !strings.Contains(content, "AGEN") {
		t.Error("Content should mention AGEN")
	}
}
