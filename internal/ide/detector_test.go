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
		{"antigravity", "Antigravity"},
		{"cursor", "Cursor"},
		{"windsurf", "Windsurf"},
		{"zed", "Zed"},
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
		{&AntigravityAdapter{}, "Antigravity"},
		{&CursorAdapter{}, "Cursor"},
		{&WindsurfAdapter{}, "Windsurf"},
		{&ZedAdapter{}, "Zed"},
	}

	for _, tt := range adapters {
		if tt.adapter.Name() != tt.name {
			t.Errorf("%T.Name() = %q, want %q", tt.adapter, tt.adapter.Name(), tt.name)
		}
	}
}
