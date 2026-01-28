// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Team CLI commands

package cli

import (
	"fmt"
	"os"

	"github.com/eshanized/agen/internal/team"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Team collaboration",
	Long: `Manage team configurations for shared agent templates.

Team configs let you:
- Require specific agents/skills across the team
- Lock template versions for consistency
- Sync configurations across team members

Examples:
  agen team init my-team  # Initialize team config
  agen team sync          # Sync with team requirements
  agen team validate      # Check if project meets requirements`,
}

var teamInitCmd = &cobra.Command{
	Use:   "init <name>",
	Short: "Initialize team config",
	Args:  cobra.ExactArgs(1),
	RunE:  runTeamInit,
}

var teamSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync with team requirements",
	RunE:  runTeamSync,
}

var teamValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate against team config",
	RunE:  runTeamValidate,
}

var teamAddCmd = &cobra.Command{
	Use:   "add <type> <name>",
	Short: "Add required agent/skill",
	Long: `Add a required agent or skill to team config.

Types: agent, skill

Examples:
  agen team add agent security-auditor
  agen team add skill clean-code`,
	Args: cobra.ExactArgs(2),
	RunE: runTeamAdd,
}

var teamRemoveCmd = &cobra.Command{
	Use:   "remove <type> <name>",
	Short: "Remove required agent/skill",
	Args:  cobra.ExactArgs(2),
	RunE:  runTeamRemove,
}

var teamLockCmd = &cobra.Command{
	Use:   "lock <name> <version>",
	Short: "Lock a template version",
	Args:  cobra.ExactArgs(2),
	RunE:  runTeamLock,
}

var teamInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show team config info",
	RunE:  runTeamInfo,
}

func init() {
	teamCmd.AddCommand(teamInitCmd)
	teamCmd.AddCommand(teamSyncCmd)
	teamCmd.AddCommand(teamValidateCmd)
	teamCmd.AddCommand(teamAddCmd)
	teamCmd.AddCommand(teamRemoveCmd)
	teamCmd.AddCommand(teamLockCmd)
	teamCmd.AddCommand(teamInfoCmd)

	rootCmd.AddCommand(teamCmd)
}

func runTeamInit(cmd *cobra.Command, args []string) error {
	name := args[0]
	cwd, _ := os.Getwd()

	config, err := team.InitTeam(cwd, name)
	if err != nil {
		return err
	}

	printSuccess("Initialized team config: %s", name)
	fmt.Printf("Config file: %s/.agen-team.json\n", cwd)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Add required agents: agen team add agent <name>")
	fmt.Println("  2. Add required skills: agen team add skill <name>")
	fmt.Println("  3. Commit .agen-team.json to version control")

	_ = config
	return nil
}

func runTeamSync(cmd *cobra.Command, args []string) error {
	cwd, _ := os.Getwd()

	config, err := team.LoadTeamConfig(cwd)
	if err != nil {
		return err
	}

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Printf("\n🔄 Syncing with team: %s\n\n", config.Name)

	result, err := config.Sync(cwd)
	if err != nil {
		return err
	}

	if len(result.Added) > 0 {
		fmt.Println("Added:")
		for _, item := range result.Added {
			printSuccess("  + %s", item)
		}
	}

	if len(result.Updated) > 0 {
		fmt.Println("Updated:")
		for _, item := range result.Updated {
			printInfo("  ~ %s", item)
		}
	}

	if len(result.Errors) > 0 {
		fmt.Println("Errors:")
		for _, err := range result.Errors {
			printError("  ! %s", err)
		}
	}

	if len(result.Added) == 0 && len(result.Updated) == 0 {
		printSuccess("Already in sync!")
	}

	return nil
}

func runTeamValidate(cmd *cobra.Command, args []string) error {
	cwd, _ := os.Getwd()

	config, err := team.LoadTeamConfig(cwd)
	if err != nil {
		return err
	}

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Printf("\n✓ Validating against: %s\n\n", config.Name)

	result := config.Validate(cwd)

	if result.Valid {
		color.New(color.FgGreen, color.Bold).Println("✨ Validation passed!")
	} else {
		color.New(color.FgRed, color.Bold).Println("❌ Validation failed!")
	}

	if len(result.Missing) > 0 {
		fmt.Println("\nMissing:")
		for _, item := range result.Missing {
			printError("  ✗ %s", item)
		}
	}

	if len(result.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		for _, warn := range result.Warnings {
			printWarning("  ⚠ %s", warn)
		}
	}

	return nil
}

func runTeamAdd(cmd *cobra.Command, args []string) error {
	itemType := args[0]
	name := args[1]

	cwd, _ := os.Getwd()
	config, err := team.LoadTeamConfig(cwd)
	if err != nil {
		return err
	}

	if err := config.AddRequired(itemType, name); err != nil {
		return err
	}

	if err := config.Save(cwd); err != nil {
		return err
	}

	printSuccess("Added required %s: %s", itemType, name)
	return nil
}

func runTeamRemove(cmd *cobra.Command, args []string) error {
	itemType := args[0]
	name := args[1]

	cwd, _ := os.Getwd()
	config, err := team.LoadTeamConfig(cwd)
	if err != nil {
		return err
	}

	if err := config.RemoveRequired(itemType, name); err != nil {
		return err
	}

	if err := config.Save(cwd); err != nil {
		return err
	}

	printSuccess("Removed %s: %s", itemType, name)
	return nil
}

func runTeamLock(cmd *cobra.Command, args []string) error {
	name := args[0]
	version := args[1]

	cwd, _ := os.Getwd()
	config, err := team.LoadTeamConfig(cwd)
	if err != nil {
		return err
	}

	config.LockVersion(name, version)

	if err := config.Save(cwd); err != nil {
		return err
	}

	printSuccess("Locked %s to version %s", name, version)
	return nil
}

func runTeamInfo(cmd *cobra.Command, args []string) error {
	cwd, _ := os.Getwd()
	config, err := team.LoadTeamConfig(cwd)
	if err != nil {
		return err
	}

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Printf("\n👥 Team: %s\n\n", config.Name)

	fmt.Printf("Version:     %s\n", config.Version)
	fmt.Printf("Created:     %s\n", config.CreatedAt.Format("2006-01-02"))
	fmt.Printf("Updated:     %s\n", config.UpdatedAt.Format("2006-01-02"))

	if len(config.RequiredAgents) > 0 {
		fmt.Println("\nRequired Agents:")
		for _, a := range config.RequiredAgents {
			fmt.Printf("  - %s\n", a)
		}
	}

	if len(config.RequiredSkills) > 0 {
		fmt.Println("\nRequired Skills:")
		for _, s := range config.RequiredSkills {
			fmt.Printf("  - %s\n", s)
		}
	}

	if len(config.LockedVersions) > 0 {
		fmt.Println("\nLocked Versions:")
		for name, version := range config.LockedVersions {
			fmt.Printf("  %s: %s\n", name, version)
		}
	}

	return nil
}
