// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// A cross-platform CLI tool for managing AI agent templates

package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Version info - will be set by goreleaser at build time
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// rootCmd is the base command when called without any subcommands.
// Think of it as the "agen" command itself before any subcommand is added.
var rootCmd = &cobra.Command{
	Use:   "agen",
	Short: "AI Agent Template Manager",
	Long: `AGEN - AI Agent Template Manager

A cross-platform CLI tool for managing AI agent templates.
Supports multiple IDEs: Antigravity, Cursor, Windsurf, and Zed.

Quick Start:
  agen init          Initialize templates in current project
  agen list          List available agents and skills
  agen status        Check installation status
  agen verify        Run verification scripts
  agen upgrade       Update agen to latest version

For more info: https://github.com/eshanized/agen`,
	// Don't show usage on errors - it's too verbose
	SilenceUsage: true,
	// We handle errors ourselves
	SilenceErrors: true,
}

// Execute runs the root command. This is the main entry point called from main.go
//
// How it works:
// 1. Cobra parses the command line args
// 2. Finds the right subcommand (or uses root if none specified)
// 3. Runs the command's Run function
// 4. Returns any error that occurred
//
// NOTE: we don't print the error here because some commands handle their own output
func Execute() error {
	return rootCmd.Execute()
}

// init registers all subcommands and sets up global flags.
// This runs automatically when the package is imported.
func init() {
	// Global flags that work on all commands
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable colored output")

	// Add version flag manually for better control
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate(getVersionTemplate())

	// Register all subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(verifyCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(profileCmd)
	// rootCmd.AddCommand(playgroundCmd) // Playground removed for production
}

// getVersionTemplate returns a nicely formatted version string.
// we include build info for debugging purposes
func getVersionTemplate() string {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	return fmt.Sprintf(`%s

Version:    %s
Commit:     %s
Built:      %s
`, cyan("AGEN - AI Agent Template Manager"), Version, Commit, BuildDate)
}

// checkVerbose is a helper to check if verbose mode is enabled.
// Used throughout commands to show extra debug info
func checkVerbose(cmd *cobra.Command) bool {
	verbose, _ := cmd.Flags().GetBool("verbose")
	return verbose
}

// printSuccess prints a success message in green
func printSuccess(format string, args ...interface{}) {
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	fmt.Printf("%s %s\n", green("✓"), fmt.Sprintf(format, args...))
}

// printError prints an error message in red to stderr
func printError(format string, args ...interface{}) {
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	fmt.Fprintf(os.Stderr, "%s %s\n", red("✗"), fmt.Sprintf(format, args...))
}

// printWarning prints a warning message in yellow
func printWarning(format string, args ...interface{}) {
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	fmt.Printf("%s %s\n", yellow("⚠"), fmt.Sprintf(format, args...))
}

// printInfo prints an info message in blue
func printInfo(format string, args ...interface{}) {
	blue := color.New(color.FgBlue).SprintFunc()
	fmt.Printf("%s %s\n", blue("ℹ"), fmt.Sprintf(format, args...))
}
