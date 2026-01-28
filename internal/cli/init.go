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
	"github.com/eshanized/agen/internal/tui"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// initCmd handles the "agen init" command which installs agent templates
// into the current project directory.
var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize agent templates in a project",
	Long: `Initialize agent templates in the current or specified project directory.

By default, this will:
1. Auto-detect which IDE you're using (Cursor, Windsurf, Zed, or Antigravity)
2. Install the appropriate configuration files
3. Copy agent templates to your project

If no IDE is detected, you'll be prompted to choose one.

Examples:
  agen init                           # Initialize in current directory
  agen init /path/to/project          # Initialize in specific directory
  agen init --ide cursor              # Force Cursor format
  agen init --agents frontend,backend # Only install specific agents`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

func init() {
	// Flags for the init command
	initCmd.Flags().StringP("ide", "i", "", "force specific IDE (antigravity, cursor, windsurf, zed)")
	initCmd.Flags().StringSliceP("agents", "a", nil, "comma-separated list of agents to install")
	initCmd.Flags().StringSliceP("skills", "s", nil, "comma-separated list of skills to install")
	initCmd.Flags().BoolP("force", "f", false, "overwrite existing files without prompting")
	initCmd.Flags().Bool("dry-run", false, "show what would be done without making changes")
	initCmd.Flags().Bool("no-wizard", false, "skip interactive wizard even if no flags provided")
}

// runInit is the main logic for the init command.
//
// How it works:
// 1. Figure out the target directory (current dir or arg)
// 2. Detect the IDE being used (or use --ide flag)
// 3. If interactive mode, launch the TUI wizard
// 4. Load templates (embedded or from network)
// 5. Filter templates based on --agents and --skills flags
// 6. Install templates using the appropriate IDE adapter
//
// Why detect IDE first? Different IDEs need different file formats.
// Cursor uses .cursorrules, Windsurf uses .windsurfrules, etc.
func runInit(cmd *cobra.Command, args []string) error {
	verbose := checkVerbose(cmd)

	// Step 1: determine target directory
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// make sure the directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", absPath)
	}

	if verbose {
		printInfo("Target directory: %s", absPath)
	}

	// Step 2: detect or get IDE
	ideName, _ := cmd.Flags().GetString("ide")
	var ideAdapter ide.Adapter

	if ideName != "" {
		// User specified IDE explicitly
		ideAdapter = ide.GetAdapter(ideName)
		if ideAdapter == nil {
			return fmt.Errorf("unknown IDE: %s (supported: antigravity, cursor, windsurf, zed)", ideName)
		}
		printInfo("Using IDE: %s (specified via --ide)", ideAdapter.Name())
	} else {
		// Try to auto-detect
		ideAdapter = ide.Detect(absPath)
		if ideAdapter != nil {
			printInfo("Detected IDE: %s", ideAdapter.Name())
		}
	}

	// Step 3: Check if we should launch interactive wizard
	noWizard, _ := cmd.Flags().GetBool("no-wizard")
	agents, _ := cmd.Flags().GetStringSlice("agents")
	skills, _ := cmd.Flags().GetStringSlice("skills")

	// Launch wizard if: no IDE detected AND no flags provided AND not disabled
	if ideAdapter == nil && len(agents) == 0 && len(skills) == 0 && !noWizard {
		// Launch interactive wizard
		// We need to load templates first for the wizard
		tmpl, err := templates.LoadEmbedded()
		if err != nil {
			return fmt.Errorf("failed to load templates for wizard: %w", err)
		}

		result, err := tui.RunWizard(tmpl)
		if err != nil {
			return fmt.Errorf("wizard failed: %w", err)
		}

		if result.Cancelled {
			color.Yellow("Setup cancelled.")
			return nil
		}

		// Apply selection
		if result.IDE != "" {
			ideAdapter = ide.GetAdapter(result.IDE)
		}

		agents = result.Agents
		skills = result.Skills

		printInfo("Selected from wizard: IDE=%s, Agents=%d, Skills=%d",
			result.IDE, len(agents), len(skills))
	} else if ideAdapter == nil && len(agents) == 0 && len(skills) == 0 {
		// Only show warning if wizard was explicitly disabled or not applicable
		printWarning("No IDE detected. Use --ide flag or run without --no-wizard for interactive mode.")
		fmt.Println("\nSupported IDEs:")
		fmt.Println("  antigravity  - Claude Code / Antigravity (full .agent/ folder)")
		fmt.Println("  cursor       - Cursor IDE (.cursorrules file)")
		fmt.Println("  windsurf     - Windsurf IDE (.windsurfrules file)")
		fmt.Println("  zed          - Zed Editor (.zed/ folder)")
		return nil
	}

	// If still no IDE, default to Antigravity
	if ideAdapter == nil {
		ideAdapter = ide.GetAdapter("antigravity")
		printInfo("Defaulting to Antigravity format")
	}

	// Step 4: Load templates
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	force, _ := cmd.Flags().GetBool("force")

	if dryRun {
		printWarning("DRY RUN: No changes will be made")
	}

	// load from embedded templates first
	tmpl, err := templates.LoadEmbedded()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	if verbose {
		printInfo("Loaded %d agents, %d skills, %d workflows",
			len(tmpl.Agents), len(tmpl.Skills), len(tmpl.Workflows))
	}

	// Step 5: Filter if specific agents/skills requested
	if len(agents) > 0 || len(skills) > 0 {
		tmpl = tmpl.Filter(agents, skills)
		if verbose {
			printInfo("Filtered to %d agents, %d skills",
				len(tmpl.Agents), len(tmpl.Skills))
		}
	}

	// Step 6: Install!
	opts := ide.InstallOptions{
		TargetDir: absPath,
		DryRun:    dryRun,
		Force:     force,
		Verbose:   verbose,
	}

	if err := ideAdapter.Install(tmpl, opts); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	// Success!
	if !dryRun {
		green := color.New(color.FgGreen, color.Bold)
		green.Println("\n✨ AGEN initialized successfully!")
		fmt.Printf("\nInstalled for: %s\n", ideAdapter.Name())
		fmt.Printf("Location: %s\n", absPath)
		fmt.Println("\nNext steps:")
		fmt.Println("  agen status   - Check installation status")
		fmt.Println("  agen list     - See available agents and skills")
		fmt.Println("  agen verify   - Run verification scripts")
	}

	return nil
}
