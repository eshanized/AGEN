// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// A cross-platform CLI tool for managing AI agent templates

package cli

import (
	"fmt"
	"path/filepath"

	"github.com/eshanized/agen/internal/ide"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// statusCmd checks and displays the current installation status.
// helps users understand what's installed and if it's up to date.
var statusCmd = &cobra.Command{
	Use:   "status [path]",
	Short: "Check installation status",
	Long: `Check the current status of AGEN templates in a project.

Shows:
- Detected IDE and configuration
- Installed agents, skills, and workflows
- Version information
- Update availability

Examples:
  agen status              # Check current directory
  agen status /path/to/proj # Check specific directory`,
	Args: cobra.MaximumNArgs(1),
	RunE: runStatus,
}

// runStatus is the main logic for the status command.
//
// How it works:
// 1. figure out which directory to check
// 2. Try to detect the IDE from existing files
// 3. Check what templates are installed
// 4. Compare with latest available version
// 5. Print a nice summary
//
// NOTE: this doesn't make network requests by default. use --check-updates for that
func runStatus(cmd *cobra.Command, args []string) error {
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

	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	cyan.Println("\n📊 AGEN Status")
	fmt.Printf("Directory: %s\n\n", absPath)

	// Step 1: Detect IDE
	ideAdapter := ide.Detect(absPath)
	if ideAdapter == nil {
		yellow.Println("⚠ No IDE detected")
		fmt.Println("  Run 'agen init' to set up templates")
		return nil
	}

	green.Printf("✓ IDE: %s\n", ideAdapter.Name())
	fmt.Printf("  Config: %s\n", ideAdapter.GetRulesPath())

	// Step 2: Check installed templates
	installed, err := ide.GetInstalledInfo(absPath, ideAdapter)
	if err != nil {
		if verbose {
			printWarning("Could not read installed info: %v", err)
		}
	}

	if installed != nil {
		fmt.Printf("\n📦 Installed:\n")
		fmt.Printf("  Agents:    %d\n", installed.AgentCount)
		fmt.Printf("  Skills:    %d\n", installed.SkillCount)
		fmt.Printf("  Workflows: %d\n", installed.WorkflowCount)
		fmt.Printf("  Version:   %s\n", installed.Version)

		if installed.ModifiedFiles > 0 {
			yellow.Printf("  ⚠ %d file(s) modified locally\n", installed.ModifiedFiles)
		}
	} else {
		fmt.Println("\n📦 No templates installed")
		fmt.Println("  Run 'agen init' to install templates")
	}

	// Step 3: show quick actions
	fmt.Println("\n💡 Actions:")
	fmt.Println("  agen list      - See available agents")
	fmt.Println("  agen update    - Update templates")
	fmt.Println("  agen verify    - Run verification scripts")
	fmt.Println()

	return nil
}
