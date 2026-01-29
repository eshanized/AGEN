// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Unit tests for Cline adapter

package ide

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestClineAdapter_Name(t *testing.T) {
	adapter := &ClineAdapter{}
	if got := adapter.Name(); got != "Cline" {
		t.Errorf("ClineAdapter.Name() = %q, want %q", got, "Cline")
	}
}

func TestClineAdapter_GetRulesPath(t *testing.T) {
	adapter := &ClineAdapter{}
	if got := adapter.GetRulesPath(); got != ".clinerules" {
		t.Errorf("ClineAdapter.GetRulesPath() = %q, want %q", got, ".clinerules")
	}
}

func TestClineAdapter_Detect(t *testing.T) {
	adapter := &ClineAdapter{}

	tests := []struct {
		name     string
		setup    func(dir string) error
		expected bool
	}{
		{
			name: "detect .clinerules file",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, ".clinerules"), []byte("# rules"), 0644)
			},
			expected: true,
		},
		{
			name: "detect .cline directory",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, ".cline"), 0755)
			},
			expected: true,
		},
		{
			name:     "no cline config",
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
				t.Errorf("ClineAdapter.Detect() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestClineAdapter_Install(t *testing.T) {
	adapter := &ClineAdapter{}
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

		rulesFile := filepath.Join(tmpDir, ".clinerules")
		if _, err := os.Stat(rulesFile); os.IsNotExist(err) {
			t.Error(".clinerules was not created")
		}
	})

	t.Run("install errors without force when exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.WriteFile(filepath.Join(tmpDir, ".clinerules"), []byte("existing"), 0644)

		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		if err := adapter.Install(tmpl, opts); err == nil {
			t.Error("Install() should error when file exists without --force")
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

		rulesFile := filepath.Join(tmpDir, ".clinerules")
		if _, err := os.Stat(rulesFile); !os.IsNotExist(err) {
			t.Error("dry run should not create files")
		}
	})
}

func TestClineAdapter_Update(t *testing.T) {
	adapter := &ClineAdapter{}
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
		os.WriteFile(filepath.Join(tmpDir, ".clinerules"), []byte("existing"), 0644)

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

func TestClineAdapter_BuildRulesContent(t *testing.T) {
	adapter := &ClineAdapter{}
	tmpl := createMockTemplates()

	content := adapter.buildRulesContent(tmpl)

	if !strings.Contains(content, "Cline Rules") {
		t.Error("Content should contain 'Cline Rules' header")
	}
	if !strings.Contains(content, "AGEN") {
		t.Error("Content should mention AGEN")
	}
}
