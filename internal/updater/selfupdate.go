// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// A cross-platform CLI tool for managing AI agent templates

package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Release represents a GitHub release
type Release struct {
	Version      string
	DownloadURL  string
	ReleaseNotes string
	PublishedAt  time.Time
}

// GitHubRelease is the API response structure
type GitHubRelease struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	PublishedAt string `json:"published_at"`
	Assets      []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

const (
	repoOwner = "eshanized"
	repoName  = "agen"
	apiBase   = "https://api.github.com"
)

// CheckForUpdate checks if a newer version is available.
//
// How it works:
// 1. Query GitHub API for latest release
// 2. Compare version strings (semantic versioning)
// 3. Return release info if newer, nil if current is latest
//
// Why GitHub API? It's the standard way to distribute Go binaries.
// We also support prerelease versions for early adopters.
func CheckForUpdate(currentVersion string) (*Release, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", apiBase, repoOwner, repoName)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "agen-updater")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		// No releases yet
		return nil, nil
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var ghRelease GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&ghRelease); err != nil {
		return nil, fmt.Errorf("failed to parse release info: %w", err)
	}

	latestVersion := strings.TrimPrefix(ghRelease.TagName, "v")
	current := strings.TrimPrefix(currentVersion, "v")

	// simple version comparison (assumes semver)
	if compareVersions(current, latestVersion) >= 0 {
		return nil, nil // already on latest
	}

	// Find download URL for current platform
	downloadURL := findAssetForPlatform(ghRelease.Assets)

	if downloadURL == "" {
		return nil, fmt.Errorf("no binary available for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	return &Release{
		Version:      latestVersion,
		DownloadURL:  downloadURL,
		ReleaseNotes: ghRelease.Body,
	}, nil
}

// DownloadAndReplace downloads the new binary and replaces the current one.
//
// How it works (the tricky part):
// 1. Download new binary to temp file
// 2. Verify it's executable
// 3. Get path of current binary
// 4. On Unix: rename current to .old, rename new to current
// 5. On Windows: create a batch script to do the swap after exit
//
// Why so complex? You can't replace a running binary on Windows.
// On Unix it technically works but we do atomic swap for safety.
func DownloadAndReplace(release *Release) error {
	if release.DownloadURL == "" {
		return fmt.Errorf("no download URL provided")
	}

	// Download to temp file
	tmpFile, err := os.CreateTemp("", "agen-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	resp, err := http.Get(release.DownloadURL)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return fmt.Errorf("failed to save update: %w", err)
	}
	tmpFile.Close()

	// Make executable
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	// Get current binary path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("failed to resolve executable path: %w", err)
	}

	if runtime.GOOS == "windows" {
		return windowsUpdate(tmpFile.Name(), execPath)
	}

	return unixUpdate(tmpFile.Name(), execPath)
}

// unixUpdate does atomic replacement on Unix systems
func unixUpdate(newPath, currentPath string) error {
	// Backup current
	backupPath := currentPath + ".old"
	os.Remove(backupPath) // ignore error if doesn't exist

	// Rename current to backup
	if err := os.Rename(currentPath, backupPath); err != nil {
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	// Move new to current
	if err := os.Rename(newPath, currentPath); err != nil {
		// try to restore backup
		os.Rename(backupPath, currentPath)
		return fmt.Errorf("failed to install new binary: %w", err)
	}

	// Clean up backup
	os.Remove(backupPath)

	return nil
}

// windowsUpdate creates a batch script to replace the binary after exit
func windowsUpdate(newPath, currentPath string) error {
	// Windows can't replace running executables, so we create a batch script
	// that runs after this process exits

	batchScript := fmt.Sprintf(`@echo off
ping 127.0.0.1 -n 2 > nul
move /y "%s" "%s"
del "%s"
`, newPath, currentPath, "%~f0")

	batchPath := filepath.Join(os.TempDir(), "agen-update.bat")
	if err := os.WriteFile(batchPath, []byte(batchScript), 0755); err != nil {
		return fmt.Errorf("failed to create update script: %w", err)
	}

	// Start the batch script detached
	cmd := exec.Command("cmd.exe", "/C", batchPath)
	// cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // Windows specific, breaks on Linux
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start update script: %w", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// findAssetForPlatform finds the right binary for the current OS/arch
func findAssetForPlatform(assets []struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}) string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	// Expected naming convention: agen_linux_amd64, agen_darwin_arm64, etc.
	expectedName := fmt.Sprintf("agen_%s_%s", os, arch)
	if os == "windows" {
		expectedName += ".exe"
	}

	for _, asset := range assets {
		if strings.Contains(strings.ToLower(asset.Name), strings.ToLower(expectedName)) {
			return asset.BrowserDownloadURL
		}
	}

	return ""
}

// compareVersions compares two semantic versions.
// Returns: -1 if a < b, 0 if a == b, 1 if a > b
func compareVersions(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	for i := 0; i < 3; i++ {
		var aNum, bNum int
		if i < len(aParts) {
			fmt.Sscanf(aParts[i], "%d", &aNum)
		}
		if i < len(bParts) {
			fmt.Sscanf(bParts[i], "%d", &bNum)
		}

		if aNum < bNum {
			return -1
		}
		if aNum > bNum {
			return 1
		}
	}

	return 0
}
