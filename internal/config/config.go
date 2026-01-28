// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// A cross-platform CLI tool for managing AI agent templates

package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds global AGEN configuration
type Config struct {
	// Analytics settings
	AnalyticsEnabled bool   `json:"analytics_enabled"`
	AnalyticsID      string `json:"analytics_id,omitempty"`

	// Update settings
	AutoCheckUpdates bool   `json:"auto_check_updates"`
	UpdateChannel    string `json:"update_channel"` // "stable" or "beta"

	// Default settings
	DefaultIDE    string `json:"default_ide,omitempty"`
	DefaultBranch string `json:"default_branch"`

	// Cache settings
	CacheDir     string `json:"cache_dir,omitempty"`
	CacheTTLDays int    `json:"cache_ttl_days"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		AnalyticsEnabled: false,
		AutoCheckUpdates: true,
		UpdateChannel:    "stable",
		DefaultBranch:    "main",
		CacheTTLDays:     7,
	}
}

// GetConfigDir returns the path to the AGEN config directory.
// Uses XDG config dir on Linux, ~/Library/Application Support on macOS,
// and %APPDATA% on Windows.
func GetConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "agen"), nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// Load reads the config file or returns defaults if it doesn't exist.
//
// How it works:
// 1. Try to find config file in standard location
// 2. If exists, parse JSON and return
// 3. If not exists, return default config (don't create file yet)
//
// This way we don't pollute the user's system until they explicitly
// change a setting.
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return DefaultConfig(), nil
	}
	if err != nil {
		return nil, err
	}

	config := DefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// Save writes the config to disk
func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Make sure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// GetCacheDir returns the cache directory, creating it if needed
func (c *Config) GetCacheDir() (string, error) {
	if c.CacheDir != "" {
		return c.CacheDir, nil
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	agenCache := filepath.Join(cacheDir, "agen")
	if err := os.MkdirAll(agenCache, 0755); err != nil {
		return "", err
	}

	return agenCache, nil
}
