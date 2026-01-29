// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Unit tests for IDE detection

package ide

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAntigravityDetect(t *testing.T) {
	// Create temp directory with .agent folder
	tmpDir := t.TempDir()
	agentDir := filepath.Join(tmpDir, ".agent")
	if err := os.MkdirAll(agentDir, 0755); err != nil {
		t.Fatal(err)
	}

	adapter := &AntigravityAdapter{}

	if !adapter.Detect(tmpDir) {
		t.Error("AntigravityAdapter.Detect() should return true when .agent/ exists")
	}

	// Test without .agent
	emptyDir := t.TempDir()
	if adapter.Detect(emptyDir) {
		t.Error("AntigravityAdapter.Detect() should return false when .agent/ doesn't exist")
	}
}

func TestCursorDetect(t *testing.T) {
	// Create temp directory with .cursorrules file
	tmpDir := t.TempDir()
	rulesFile := filepath.Join(tmpDir, ".cursorrules")
	if err := os.WriteFile(rulesFile, []byte("# rules"), 0644); err != nil {
		t.Fatal(err)
	}

	adapter := &CursorAdapter{}

	if !adapter.Detect(tmpDir) {
		t.Error("CursorAdapter.Detect() should return true when .cursorrules exists")
	}

	// Test without .cursorrules
	emptyDir := t.TempDir()
	if adapter.Detect(emptyDir) {
		t.Error("CursorAdapter.Detect() should return false when .cursorrules doesn't exist")
	}
}

func TestWindsurfDetect(t *testing.T) {
	// Create temp directory with .windsurfrules file
	tmpDir := t.TempDir()
	rulesFile := filepath.Join(tmpDir, ".windsurfrules")
	if err := os.WriteFile(rulesFile, []byte("# rules"), 0644); err != nil {
		t.Fatal(err)
	}

	adapter := &WindsurfAdapter{}

	if !adapter.Detect(tmpDir) {
		t.Error("WindsurfAdapter.Detect() should return true when .windsurfrules exists")
	}
}

func TestZedDetect(t *testing.T) {
	// Create temp directory with .zed folder
	tmpDir := t.TempDir()
	zedDir := filepath.Join(tmpDir, ".zed")
	if err := os.MkdirAll(zedDir, 0755); err != nil {
		t.Fatal(err)
	}

	adapter := &ZedAdapter{}

	if !adapter.Detect(tmpDir) {
		t.Error("ZedAdapter.Detect() should return true when .zed/ exists")
	}
}

func TestDetectPriority(t *testing.T) {
	// Create directory with both .cursorrules and .agent
	// .cursorrules should be detected first (priority order)
	tmpDir := t.TempDir()

	// Create .agent folder
	agentDir := filepath.Join(tmpDir, ".agent")
	if err := os.MkdirAll(agentDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create .cursorrules
	rulesFile := filepath.Join(tmpDir, ".cursorrules")
	if err := os.WriteFile(rulesFile, []byte("# rules"), 0644); err != nil {
		t.Fatal(err)
	}

	detected := Detect(tmpDir)
	if detected == nil {
		t.Fatal("Detect() returned nil")
	}

	// Cursor should be detected first based on priority
	if detected.Name() != "Cursor" {
		t.Errorf("Expected Cursor to be detected (priority), got %s", detected.Name())
	}
}

func TestGetAdapter(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		// Original adapters
		{"antigravity", "Antigravity"},
		{"cursor", "Cursor"},
		{"windsurf", "Windsurf"},
		{"zed", "Zed"},
		// New VS Code extensions
		{"continue", "Continue"},
		{"cline", "Cline"},
		// Desktop IDEs
		{"jetbrains", "JetBrains"},
		{"neovim", "Neovim"},
		{"emacs", "Emacs"},
		// CLI tools
		{"aider", "Aider"},
		// Cloud/Platform
		{"claudecode", "ClaudeCode"},
		{"copilotworkspace", "CopilotWorkspace"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := GetAdapter(tt.name)
			if adapter == nil {
				t.Fatalf("GetAdapter(%q) returned nil", tt.name)
			}
			if adapter.Name() != tt.expected {
				t.Errorf("GetAdapter(%q).Name() = %q, want %q", tt.name, adapter.Name(), tt.expected)
			}
		})
	}

	// Test unknown adapter
	if adapter := GetAdapter("unknown"); adapter != nil {
		t.Error("GetAdapter('unknown') should return nil")
	}
}

func TestAdapterNames(t *testing.T) {
	adapters := []struct {
		adapter Adapter
		name    string
	}{
		// Original
		{&AntigravityAdapter{}, "Antigravity"},
		{&CursorAdapter{}, "Cursor"},
		{&WindsurfAdapter{}, "Windsurf"},
		{&ZedAdapter{}, "Zed"},
		// New
		{&ContinueAdapter{}, "Continue"},
		{&ClineAdapter{}, "Cline"},
		{&JetBrainsAdapter{}, "JetBrains"},
		{&NeovimAdapter{}, "Neovim"},
		{&EmacsAdapter{}, "Emacs"},
		{&AiderAdapter{}, "Aider"},
		{&ClaudeCodeAdapter{}, "ClaudeCode"},
		{&CopilotWorkspaceAdapter{}, "CopilotWorkspace"},
	}

	for _, tt := range adapters {
		if tt.adapter.Name() != tt.name {
			t.Errorf("%T.Name() = %q, want %q", tt.adapter, tt.adapter.Name(), tt.name)
		}
	}
}

func TestAdapterGetRulesPath(t *testing.T) {
	tests := []struct {
		adapter  Adapter
		expected string
	}{
		{&AntigravityAdapter{}, ".agent/rules/GEMINI.md"},
		{&CursorAdapter{}, ".cursorrules"},
		{&WindsurfAdapter{}, ".windsurfrules"},
		{&ZedAdapter{}, ".zed/settings.json"},
		{&ContinueAdapter{}, ".continuerules"},
		{&ClineAdapter{}, ".clinerules"},
		{&JetBrainsAdapter{}, ".jbrules.md"},
		{&NeovimAdapter{}, ".nvim/ai-rules.md"},
		{&EmacsAdapter{}, ".emacs-project/ai-context.md"},
		{&AiderAdapter{}, ".aider-context.md"},
		{&ClaudeCodeAdapter{}, "CLAUDE.md"},
		{&CopilotWorkspaceAdapter{}, ".github/copilot-instructions.md"},
	}

	for _, tt := range tests {
		t.Run(tt.adapter.Name(), func(t *testing.T) {
			if got := tt.adapter.GetRulesPath(); got != tt.expected {
				t.Errorf("%T.GetRulesPath() = %q, want %q", tt.adapter, got, tt.expected)
			}
		})
	}
}

func TestDetectNoIDE(t *testing.T) {
	tmpDir := t.TempDir()
	adapter := Detect(tmpDir)
	if adapter != nil {
		t.Errorf("Detect() on empty directory should return nil, got %s", adapter.Name())
	}
}

func TestDetectAllIDEs(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(dir string) error
		expected string
	}{
		{
			name: "detect cursor",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, ".cursorrules"), []byte("# rules"), 0644)
			},
			expected: "Cursor",
		},
		{
			name: "detect windsurf",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, ".windsurfrules"), []byte("# rules"), 0644)
			},
			expected: "Windsurf",
		},
		{
			name: "detect cline",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, ".clinerules"), []byte("# rules"), 0644)
			},
			expected: "Cline",
		},
		{
			name: "detect continue",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, ".continue"), 0755)
			},
			expected: "Continue",
		},
		{
			name: "detect claude code",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# claude"), 0644)
			},
			expected: "ClaudeCode",
		},
		{
			name: "detect copilot workspace",
			setup: func(dir string) error {
				githubDir := filepath.Join(dir, ".github")
				if err := os.MkdirAll(githubDir, 0755); err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(githubDir, "copilot-instructions.md"), []byte("# instructions"), 0644)
			},
			expected: "CopilotWorkspace",
		},
		{
			name: "detect aider",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, ".aider.conf.yml"), []byte("# config"), 0644)
			},
			expected: "Aider",
		},
		{
			name: "detect jetbrains",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, ".idea"), 0755)
			},
			expected: "JetBrains",
		},
		{
			name: "detect zed",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, ".zed"), 0755)
			},
			expected: "Zed",
		},
		{
			name: "detect neovim",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, ".nvim"), 0755)
			},
			expected: "Neovim",
		},
		{
			name: "detect emacs",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, ".dir-locals.el"), []byte(";;; config"), 0644)
			},
			expected: "Emacs",
		},
		{
			name: "detect antigravity",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, ".agent"), 0755)
			},
			expected: "Antigravity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			if err := tt.setup(tmpDir); err != nil {
				t.Fatal(err)
			}
			adapter := Detect(tmpDir)
			if adapter == nil {
				t.Fatalf("Detect() returned nil, expected %s", tt.expected)
			}
			if adapter.Name() != tt.expected {
				t.Errorf("Detect() = %s, want %s", adapter.Name(), tt.expected)
			}
		})
	}
}

func TestRegisterAdapter(t *testing.T) {
	// Test registering a custom adapter
	customAdapter := &CursorAdapter{} // Using CursorAdapter as a stand-in
	RegisterAdapter("custom-test", customAdapter)

	retrieved := GetAdapter("custom-test")
	if retrieved == nil {
		t.Fatal("RegisterAdapter() should register the adapter")
	}
	if retrieved != customAdapter {
		t.Error("GetAdapter() should return the registered adapter")
	}
}
