// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// A cross-platform CLI tool for managing AI agent templates

package cli

import (
	"fmt"
	"runtime"

	"github.com/eshanized/agen/internal/updater"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// upgradeCmd handles self-updating the agen binary.
// downloads the latest version from GitHub releases and replaces itself.
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade agen to the latest version",
	Long: `Upgrade the agen binary to the latest version from GitHub releases.

The binary will be replaced in-place. On Windows, a helper script
may be created to complete the upgrade after the command exits.

Examples:
  agen upgrade        # Upgrade to latest version
  agen upgrade --check  # Just check if update is available`,
	RunE: runUpgrade,
}

func init() {
	upgradeCmd.Flags().Bool("check", false, "only check for updates, don't install")
	upgradeCmd.Flags().Bool("force", false, "upgrade even if already on latest version")
}

// runUpgrade is the main logic for the upgrade command.
//
// How it works (the tricky part):
//  1. Query GitHub releases API for latest version
//  2. Compare with current version (Version variable from root.go)
//  3. If newer version available:
//     a. Download the binary for current OS/arch
//     b. Verify checksum (if provided)
//     c. Replace current binary with new one
//  4. On Windows, we can't replace a running binary, so we create
//     a batch script that runs after this process exits
//
// Why self-update? Makes it super easy for users to stay current.
// No need to mess with package managers or manual downloads.
func runUpgrade(cmd *cobra.Command, args []string) error {
	checkOnly, _ := cmd.Flags().GetBool("check")
	force, _ := cmd.Flags().GetBool("force")

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n🚀 AGEN Upgrade")
	fmt.Printf("Current version: %s\n", Version)
	fmt.Printf("Platform: %s/%s\n\n", runtime.GOOS, runtime.GOARCH)

	// Step 1: Check for updates
	printInfo("Checking for updates...")
	release, err := updater.CheckForUpdate(Version)
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if release == nil {
		printSuccess("You're already on the latest version!")
		if !force {
			return nil
		}
		printInfo("Force flag set, proceeding anyway...")
	} else {
		fmt.Printf("\n📦 New version available: %s\n", color.GreenString(release.Version))
		if release.ReleaseNotes != "" {
			fmt.Println("\nRelease notes:")
			fmt.Println(release.ReleaseNotes)
		}
	}

	// If just checking, stop here
	if checkOnly {
		if release != nil {
			fmt.Println("\nRun 'agen upgrade' to install the update")
		}
		return nil
	}

	// Step 2: Download and install
	if release == nil {
		// force mode with same version - redownload current
		release = &updater.Release{Version: Version}
	}

	printInfo("Downloading %s for %s/%s...", release.Version, runtime.GOOS, runtime.GOARCH)
	if err := updater.DownloadAndReplace(release); err != nil {
		return fmt.Errorf("upgrade failed: %w", err)
	}

	green := color.New(color.FgGreen, color.Bold)
	green.Printf("\n✨ Successfully upgraded to %s!\n", release.Version)

	// on windows, remind user about the pending operation
	if runtime.GOOS == "windows" {
		printInfo("Note: The upgrade will complete when you close this window")
	}

	return nil
}
