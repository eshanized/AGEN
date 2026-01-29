// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Unit tests for Cursor adapter

package ide

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCursorAdapter_Name(t *testing.T) {
	adapter := &CursorAdapter{}
	if got := adapter.Name(); got != "Cursor" {
		t.Errorf("CursorAdapter.Name() = %q, want %q", got, "Cursor")
	}
}

func TestCursorAdapter_GetRulesPath(t *testing.T) {
	adapter := &CursorAdapter{}
	if got := adapter.GetRulesPath(); got != ".cursorrules" {
		t.Errorf("CursorAdapter.GetRulesPath() = %q, want %q", got, ".cursorrules")
	}
}

func TestCursorAdapter_Detect(t *testing.T) {
	adapter := &CursorAdapter{}

	tests := []struct {
		name     string
		setup    func(dir string) error
		expected bool
	}{
		{
			name: "detect .cursorrules file",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, ".cursorrules"), []byte("# rules"), 0644)
			},
			expected: true,
		},
		{
			name:     "no cursor config",
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
				t.Errorf("CursorAdapter.Detect() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCursorAdapter_Install(t *testing.T) {
	adapter := &CursorAdapter{}
	tmpl := createMockTemplates()

	t.Run("install creates .cursorrules", func(t *testing.T) {
		tmpDir := t.TempDir()
		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		if err := adapter.Install(tmpl, opts); err != nil {
			t.Fatalf("Install() error = %v", err)
		}

		rulesFile := filepath.Join(tmpDir, ".cursorrules")
		if _, err := os.Stat(rulesFile); os.IsNotExist(err) {
			t.Error(".cursorrules was not created")
		}
	})

	t.Run("install errors without force when exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.WriteFile(filepath.Join(tmpDir, ".cursorrules"), []byte("existing"), 0644)

		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		if err := adapter.Install(tmpl, opts); err == nil {
			t.Error("Install() should error when .cursorrules exists without --force")
		}
	})

	t.Run("install succeeds with force", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.WriteFile(filepath.Join(tmpDir, ".cursorrules"), []byte("existing"), 0644)

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

		rulesFile := filepath.Join(tmpDir, ".cursorrules")
		if _, err := os.Stat(rulesFile); !os.IsNotExist(err) {
			t.Error("dry run should not create files")
		}
	})
}

func TestCursorAdapter_Update(t *testing.T) {
	adapter := &CursorAdapter{}
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
		os.WriteFile(filepath.Join(tmpDir, ".cursorrules"), []byte("existing"), 0644)

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
		os.WriteFile(filepath.Join(tmpDir, ".cursorrules"), []byte("existing"), 0644)

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

func TestCursorAdapter_BuildRulesContent(t *testing.T) {
	adapter := &CursorAdapter{}
	tmpl := createMockTemplates()

	content := adapter.buildRulesContent(tmpl)

	if !strings.Contains(content, "Cursor Rules") {
		t.Error("Content should contain 'Cursor Rules' header")
	}
	if !strings.Contains(content, "AGEN") {
		t.Error("Content should mention AGEN")
	}
}
