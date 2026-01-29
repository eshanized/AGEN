// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Unit tests for Continue adapter

package ide

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestContinueAdapter_Name(t *testing.T) {
	adapter := &ContinueAdapter{}
	if got := adapter.Name(); got != "Continue" {
		t.Errorf("ContinueAdapter.Name() = %q, want %q", got, "Continue")
	}
}

func TestContinueAdapter_GetRulesPath(t *testing.T) {
	adapter := &ContinueAdapter{}
	if got := adapter.GetRulesPath(); got != ".continuerules" {
		t.Errorf("ContinueAdapter.GetRulesPath() = %q, want %q", got, ".continuerules")
	}
}

func TestContinueAdapter_Detect(t *testing.T) {
	adapter := &ContinueAdapter{}

	tests := []struct {
		name     string
		setup    func(dir string) error
		expected bool
	}{
		{
			name: "detect .continue directory",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, ".continue"), 0755)
			},
			expected: true,
		},
		{
			name: "detect .continuerules file",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, ".continuerules"), []byte("# rules"), 0644)
			},
			expected: true,
		},
		{
			name:     "no continue config",
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
				t.Errorf("ContinueAdapter.Detect() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestContinueAdapter_Install(t *testing.T) {
	adapter := &ContinueAdapter{}
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

		// Check .continuerules exists
		rulesFile := filepath.Join(tmpDir, ".continuerules")
		if _, err := os.Stat(rulesFile); os.IsNotExist(err) {
			t.Error(".continuerules was not created")
		}

		// Check .continue/config.json exists
		configFile := filepath.Join(tmpDir, ".continue", "config.json")
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			t.Error(".continue/config.json was not created")
		}
	})

	t.Run("install errors without force when exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.MkdirAll(filepath.Join(tmpDir, ".continue"), 0755)

		opts := InstallOptions{
			TargetDir: tmpDir,
			Force:     false,
		}

		if err := adapter.Install(tmpl, opts); err == nil {
			t.Error("Install() should error when directory exists without --force")
		}
	})

	t.Run("install succeeds with force", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.MkdirAll(filepath.Join(tmpDir, ".continue"), 0755)

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

		rulesFile := filepath.Join(tmpDir, ".continuerules")
		if _, err := os.Stat(rulesFile); !os.IsNotExist(err) {
			t.Error("dry run should not create files")
		}
	})
}

func TestContinueAdapter_Update(t *testing.T) {
	adapter := &ContinueAdapter{}
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
		rulesFile := filepath.Join(tmpDir, ".continuerules")
		os.WriteFile(rulesFile, []byte("existing"), 0644)

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
		rulesFile := filepath.Join(tmpDir, ".continuerules")
		os.WriteFile(rulesFile, []byte("existing"), 0644)

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

func TestContinueAdapter_BuildRulesContent(t *testing.T) {
	adapter := &ContinueAdapter{}
	tmpl := createMockTemplates()

	content := adapter.buildRulesContent(tmpl)

	if !strings.Contains(content, "Continue Rules") {
		t.Error("Content should contain 'Continue Rules' header")
	}
	if !strings.Contains(content, "AGEN") {
		t.Error("Content should mention AGEN")
	}
	if !strings.Contains(content, "test-agent") {
		t.Error("Content should include test agent")
	}
	if !strings.Contains(content, "test-skill") {
		t.Error("Content should include test skill")
	}
}

func TestContinueAdapter_BuildConfigContent(t *testing.T) {
	adapter := &ContinueAdapter{}
	content := adapter.buildConfigContent()

	if !strings.Contains(content, "continue.dev") {
		t.Error("Config should contain continue.dev schema reference")
	}
	if !strings.Contains(content, "contextProviders") {
		t.Error("Config should contain contextProviders")
	}
}
