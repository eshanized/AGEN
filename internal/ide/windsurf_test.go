// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Unit tests for Windsurf adapter

package ide

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWindsurfAdapter_Name(t *testing.T) {
	adapter := &WindsurfAdapter{}
	if got := adapter.Name(); got != "Windsurf" {
		t.Errorf("WindsurfAdapter.Name() = %q, want %q", got, "Windsurf")
	}
}

func TestWindsurfAdapter_GetRulesPath(t *testing.T) {
	adapter := &WindsurfAdapter{}
	if got := adapter.GetRulesPath(); got != ".windsurfrules" {
		t.Errorf("WindsurfAdapter.GetRulesPath() = %q, want %q", got, ".windsurfrules")
	}
}

func TestWindsurfAdapter_Detect(t *testing.T) {
	adapter := &WindsurfAdapter{}

	tests := []struct {
		name     string
		setup    func(dir string) error
		expected bool
	}{
		{
			name: "detect .windsurfrules file",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, ".windsurfrules"), []byte("# rules"), 0644)
			},
			expected: true,
		},
		{
			name:     "no windsurf config",
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
				t.Errorf("WindsurfAdapter.Detect() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestWindsurfAdapter_Install(t *testing.T) {
	adapter := &WindsurfAdapter{}
	tmpl := createMockTemplates()

	t.Run("install creates .windsurfrules", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		if err := adapter.Install(tmpl, opts); err != nil {
			t.Fatalf("Install() error = %v", err)
		}

		rulesFile := filepath.Join(tmpDir, ".windsurfrules")
		if _, err := os.Stat(rulesFile); os.IsNotExist(err) {
			t.Error(".windsurfrules was not created")
		}
	})

	t.Run("install overwrites when file exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.WriteFile(filepath.Join(tmpDir, ".windsurfrules"), []byte("existing"), 0644)

		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     true,
		}

		if err := adapter.Install(tmpl, opts); err != nil {
			t.Errorf("Install() with force error = %v", err)
		}
	})

	t.Run("file content is updated with force", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.WriteFile(filepath.Join(tmpDir, ".windsurfrules"), []byte("existing"), 0644)

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

		rulesFile := filepath.Join(tmpDir, ".windsurfrules")
		if _, err := os.Stat(rulesFile); !os.IsNotExist(err) {
			t.Error("dry run should not create files")
		}
	})
}

func TestWindsurfAdapter_Update(t *testing.T) {
	adapter := &WindsurfAdapter{}
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
		os.WriteFile(filepath.Join(tmpDir, ".windsurfrules"), []byte("existing"), 0644)

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
		os.WriteFile(filepath.Join(tmpDir, ".windsurfrules"), []byte("existing"), 0644)

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

func TestWindsurfAdapter_BuildRulesContent(t *testing.T) {
	adapter := &WindsurfAdapter{}
	tmpl := createMockTemplates()

	content := adapter.buildRulesContent(tmpl)

	if !strings.Contains(content, "Windsurf Rules") {
		t.Error("Content should contain 'Windsurf Rules' header")
	}
	if !strings.Contains(content, "AGEN") {
		t.Error("Content should mention AGEN")
	}
}
