// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// A cross-platform CLI tool for managing AI agent templates

package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/eshanized/agen/internal/ide"
	"github.com/eshanized/agen/internal/templates"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// updateCmd updates templates from GitHub to the latest version.
// handles smart conflict resolution when user has modified files.
var updateCmd = &cobra.Command{
	Use:   "update [path]",
	Short: "Update templates to latest version",
	Long: `Update agent templates to the latest version from GitHub.

Features:
- Downloads latest templates from the repository
- Detects user modifications and offers merge options
- Creates backups before overwriting
- Supports specific branch selection

Examples:
  agen update                # Update current directory
  agen update --branch dev   # Update from dev branch
  agen update --force        # Overwrite without prompting`,
	Args: cobra.MaximumNArgs(1),
	RunE: runUpdate,
}

func init() {
	updateCmd.Flags().String("branch", "main", "git branch to fetch from")
	updateCmd.Flags().BoolP("force", "f", false, "overwrite modified files without prompting")
	updateCmd.Flags().Bool("dry-run", false, "show what would be updated without making changes")
	updateCmd.Flags().Bool("no-backup", false, "don't create backups of modified files")
}

// runUpdate is the main logic for the update command.
//
// How it works:
// 1. Check if templates are already installed
// 2. Fetch latest templates from GitHub
// 3. Compare with installed templates to find changes
// 4. For each changed file:
//   - If user hasn't modified it → update silently
//   - If user HAS modified it → prompt for action (merge/overwrite/skip)
//
// 5. Create backups unless --no-backup is specified
//
// Why smart conflict resolution? Users often customize their templates.
// We don't want to blow away their changes on every update. This approach
// lets them keep their customizations while still getting new features.
func runUpdate(cmd *cobra.Command, args []string) error {
	verbose := checkVerbose(cmd)

	// determine target directory
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// check directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", absPath)
	}

	// Get flags
	branch, _ := cmd.Flags().GetString("branch")
	force, _ := cmd.Flags().GetBool("force")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n🔄 AGEN Update")
	fmt.Printf("Directory: %s\n", absPath)
	fmt.Printf("Branch: %s\n\n", branch)

	if dryRun {
		printWarning("DRY RUN: No changes will be made")
	}

	// Step 1: Detect IDE
	ideAdapter := ide.Detect(absPath)
	if ideAdapter == nil {
		return fmt.Errorf("no AGEN installation found. Run 'agen init' first")
	}

	printInfo("Detected IDE: %s", ideAdapter.Name())

	// Step 2: Fetch latest templates
	printInfo("Fetching latest templates from GitHub...")
	latest, err := templates.FetchFromGitHub(branch)
	if err != nil {
		// Fall back to embedded if network fails
		printWarning("Network fetch failed, using embedded templates: %v", err)
		latest, err = templates.LoadEmbedded()
		if err != nil {
			return fmt.Errorf("failed to load templates: %w", err)
		}
	}

	if verbose {
		printInfo("Fetched %d agents, %d skills, %d workflows",
			len(latest.Agents), len(latest.Skills), len(latest.Workflows))
	}

	// Step 3: Compare and update
	opts := ide.UpdateOptions{
		TargetDir: absPath,
		DryRun:    dryRun,
		Force:     force,
		Verbose:   verbose,
	}

	changes, err := ideAdapter.Update(latest, opts)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	// Step 4: Print summary
	if len(changes.Updated) == 0 && len(changes.Added) == 0 {
		printSuccess("Already up to date!")
	} else {
		if len(changes.Added) > 0 {
			fmt.Printf("\n📦 Added %d new file(s):\n", len(changes.Added))
			for _, f := range changes.Added {
				fmt.Printf("  + %s\n", color.GreenString(f))
			}
		}

		if len(changes.Updated) > 0 {
			fmt.Printf("\n🔄 Updated %d file(s):\n", len(changes.Updated))
			for _, f := range changes.Updated {
				fmt.Printf("  ~ %s\n", color.YellowString(f))
			}
		}

		if len(changes.Skipped) > 0 {
			fmt.Printf("\n⏭ Skipped %d file(s) (user modified):\n", len(changes.Skipped))
			for _, f := range changes.Skipped {
				fmt.Printf("  - %s\n", color.CyanString(f))
			}
		}

		if !dryRun {
			green := color.New(color.FgGreen, color.Bold)
			green.Println("\n✨ Update complete!")
		}
	}

	return nil
}
