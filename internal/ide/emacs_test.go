// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Unit tests for Emacs adapter

package ide

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEmacsAdapter_Name(t *testing.T) {
	adapter := &EmacsAdapter{}
	if got := adapter.Name(); got != "Emacs" {
		t.Errorf("EmacsAdapter.Name() = %q, want %q", got, "Emacs")
	}
}

func TestEmacsAdapter_GetRulesPath(t *testing.T) {
	adapter := &EmacsAdapter{}
	if got := adapter.GetRulesPath(); got != ".emacs-project/ai-context.md" {
		t.Errorf("EmacsAdapter.GetRulesPath() = %q, want %q", got, ".emacs-project/ai-context.md")
	}
}

func TestEmacsAdapter_Detect(t *testing.T) {
	adapter := &EmacsAdapter{}

	tests := []struct {
		name     string
		setup    func(dir string) error
		expected bool
	}{
		{
			name: "detect .dir-locals.el file",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, ".dir-locals.el"), []byte(";;; config"), 0644)
			},
			expected: true,
		},
		{
			name: "detect .emacs.d directory",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, ".emacs.d"), 0755)
			},
			expected: true,
		},
		{
			name:     "no emacs config",
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
				t.Errorf("EmacsAdapter.Detect() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEmacsAdapter_Install(t *testing.T) {
	adapter := &EmacsAdapter{}
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

		// Check .emacs-project/ai-context.md exists
		contextFile := filepath.Join(tmpDir, ".emacs-project", "ai-context.md")
		if _, err := os.Stat(contextFile); os.IsNotExist(err) {
			t.Error(".emacs-project/ai-context.md was not created")
		}

		// Check .dir-locals.el exists
		dirLocals := filepath.Join(tmpDir, ".dir-locals.el")
		if _, err := os.Stat(dirLocals); os.IsNotExist(err) {
			t.Error(".dir-locals.el was not created")
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

		emacsDir := filepath.Join(tmpDir, ".emacs-project")
		if _, err := os.Stat(emacsDir); !os.IsNotExist(err) {
			t.Error("dry run should not create files")
		}
	})
}

func TestEmacsAdapter_Update(t *testing.T) {
	adapter := &EmacsAdapter{}
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
}

func TestEmacsAdapter_BuildRulesContent(t *testing.T) {
	adapter := &EmacsAdapter{}
	tmpl := createMockTemplates()

	content := adapter.buildRulesContent(tmpl)

	if !strings.Contains(content, "Emacs AI Context") {
		t.Error("Content should contain 'Emacs AI Context' header")
	}
	if !strings.Contains(content, "AGEN") {
		t.Error("Content should mention AGEN")
	}
}

func TestEmacsAdapter_BuildDirLocals(t *testing.T) {
	adapter := &EmacsAdapter{}
	content := adapter.buildDirLocals()

	if !strings.Contains(content, ".dir-locals.el") {
		t.Error("Content should reference .dir-locals.el")
	}
	if !strings.Contains(content, "gptel-context-file") {
		t.Error("Content should reference gptel")
	}
}
