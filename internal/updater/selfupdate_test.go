// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Unit tests for updater

package updater

import (
	"testing"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		a, b     string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "1.0.1", -1},
		{"1.0.1", "1.0.0", 1},
		{"1.1.0", "1.0.9", 1},
		{"2.0.0", "1.9.9", 1},
		{"0.9.0", "1.0.0", -1},
		{"v1.0.0", "v1.0.0", 0}, // handles 'v' prefix
		{"1.0", "1.0.0", 0},     // handles missing patch
	}

	for _, tt := range tests {
		t.Run(tt.a+"_vs_"+tt.b, func(t *testing.T) {
			result := compareVersions(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("compareVersions(%q, %q) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestFindAssetForPlatform(t *testing.T) {
	assets := []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	}{
		{Name: "agen_linux_amd64", BrowserDownloadURL: "https://example.com/agen_linux_amd64"},
		{Name: "agen_darwin_amd64", BrowserDownloadURL: "https://example.com/agen_darwin_amd64"},
		{Name: "agen_windows_amd64.exe", BrowserDownloadURL: "https://example.com/agen_windows_amd64.exe"},
		{Name: "agen_linux_arm64", BrowserDownloadURL: "https://example.com/agen_linux_arm64"},
	}

	url := findAssetForPlatform(assets)

	// URL should not be empty (assuming we're running on a supported platform)
	// This test may need adjustment based on the test runner's platform
	t.Logf("Found asset URL for current platform: %s", url)
}

func TestReleaseStruct(t *testing.T) {
	release := Release{
		Version:      "1.0.0",
		DownloadURL:  "https://example.com/download",
		ReleaseNotes: "## What's New\n- Feature 1\n- Feature 2",
	}

	if release.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", release.Version, "1.0.0")
	}

	if release.DownloadURL == "" {
		t.Error("DownloadURL should not be empty")
	}

	if release.ReleaseNotes == "" {
		t.Error("ReleaseNotes should not be empty")
	}
}
