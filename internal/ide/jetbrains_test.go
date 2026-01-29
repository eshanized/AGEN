// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Unit tests for JetBrains adapter

package ide

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestJetBrainsAdapter_Name(t *testing.T) {
	adapter := &JetBrainsAdapter{}
	if got := adapter.Name(); got != "JetBrains" {
		t.Errorf("JetBrainsAdapter.Name() = %q, want %q", got, "JetBrains")
	}
}

func TestJetBrainsAdapter_GetRulesPath(t *testing.T) {
	adapter := &JetBrainsAdapter{}
	if got := adapter.GetRulesPath(); got != ".jbrules.md" {
		t.Errorf("JetBrainsAdapter.GetRulesPath() = %q, want %q", got, ".jbrules.md")
	}
}

func TestJetBrainsAdapter_Detect(t *testing.T) {
	adapter := &JetBrainsAdapter{}

	tests := []struct {
		name     string
		setup    func(dir string) error
		expected bool
	}{
		{
			name: "detect .idea directory",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, ".idea"), 0755)
			},
			expected: true,
		},
		{
			name:     "no jetbrains config",
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
				t.Errorf("JetBrainsAdapter.Detect() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestJetBrainsAdapter_Install(t *testing.T) {
	adapter := &JetBrainsAdapter{}
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

		// Check .idea/ai-assistant.xml exists
		aiConfig := filepath.Join(tmpDir, ".idea", "ai-assistant.xml")
		if _, err := os.Stat(aiConfig); os.IsNotExist(err) {
			t.Error(".idea/ai-assistant.xml was not created")
		}

		// Check .jbrules.md exists
		rulesFile := filepath.Join(tmpDir, ".jbrules.md")
		if _, err := os.Stat(rulesFile); os.IsNotExist(err) {
			t.Error(".jbrules.md was not created")
		}
	})

	t.Run("install errors without force when exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		ideaDir := filepath.Join(tmpDir, ".idea")
		os.MkdirAll(ideaDir, 0755)
		os.WriteFile(filepath.Join(ideaDir, "ai-assistant.xml"), []byte("existing"), 0644)

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

		rulesFile := filepath.Join(tmpDir, ".jbrules.md")
		if _, err := os.Stat(rulesFile); !os.IsNotExist(err) {
			t.Error("dry run should not create files")
		}
	})
}

func TestJetBrainsAdapter_Update(t *testing.T) {
	adapter := &JetBrainsAdapter{}
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
		os.WriteFile(filepath.Join(tmpDir, ".jbrules.md"), []byte("existing"), 0644)

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

func TestJetBrainsAdapter_BuildAIConfig(t *testing.T) {
	adapter := &JetBrainsAdapter{}
	content := adapter.buildAIConfig()

	if !strings.Contains(content, "AIAssistantSettings") {
		t.Error("Config should contain AIAssistantSettings")
	}
	if !strings.Contains(content, "enableCodeCompletion") {
		t.Error("Config should contain enableCodeCompletion")
	}
}

func TestJetBrainsAdapter_BuildRulesContent(t *testing.T) {
	adapter := &JetBrainsAdapter{}
	tmpl := createMockTemplates()

	content := adapter.buildRulesContent(tmpl)

	if !strings.Contains(content, "Project Rules") {
		t.Error("Content should contain 'Project Rules' header")
	}
	if !strings.Contains(content, "AGEN") {
		t.Error("Content should mention AGEN")
	}
}
