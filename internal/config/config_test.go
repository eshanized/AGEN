// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Unit tests for configuration

package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	// Check defaults
	if cfg.AnalyticsEnabled {
		t.Error("Analytics should be disabled by default")
	}

	if !cfg.AutoCheckUpdates {
		t.Error("AutoCheckUpdates should be enabled by default")
	}

	if cfg.UpdateChannel != "stable" {
		t.Errorf("UpdateChannel = %q, want 'stable'", cfg.UpdateChannel)
	}

	if cfg.DefaultBranch != "main" {
		t.Errorf("DefaultBranch = %q, want 'main'", cfg.DefaultBranch)
	}

	if cfg.CacheTTLDays != 7 {
		t.Errorf("CacheTTLDays = %d, want 7", cfg.CacheTTLDays)
	}
}

func TestGetConfigDir(t *testing.T) {
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir() failed: %v", err)
	}

	if dir == "" {
		t.Error("GetConfigDir() returned empty string")
	}

	t.Logf("Config directory: %s", dir)
}

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() failed: %v", err)
	}

	if path == "" {
		t.Error("GetConfigPath() returned empty string")
	}

	t.Logf("Config path: %s", path)
}

func TestLoadReturnsDefaultWhenNoFile(t *testing.T) {
	// This should return default config when no file exists
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	// Should have default values
	if cfg.UpdateChannel != "stable" {
		t.Errorf("Should have default UpdateChannel, got %q", cfg.UpdateChannel)
	}
}
