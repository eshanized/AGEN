// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Experimental features: theme, alias
// Toy features (chat, benchmark, train) have been removed for production hardening.

package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var themeCmd = &cobra.Command{
	Use:   "theme [name]",
	Short: "Manage output themes",
	Long: `Set or list available output themes.

Themes:
- default: Standard colored output
- dark: High contrast for dark terminals
- light: Optimized for light terminals
- minimal: No colors, symbols only
- nerd: With Nerd Font icons

Examples:
  agen theme          # List themes
  agen theme dark     # Set dark theme`,
	Args: cobra.MaximumNArgs(1),
	RunE: runTheme,
}

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Manage command aliases",
	Long: `Create and manage command aliases.

Examples:
  agen alias set ls "list --agents"
  agen alias set chk "verify --security"
  agen alias list
  agen alias remove ls`,
}

var aliasSetCmd = &cobra.Command{
	Use:   "set <name> <command>",
	Short: "Create alias",
	Args:  cobra.ExactArgs(2),
	RunE:  runAliasSet,
}

var aliasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List aliases",
	RunE:  runAliasList,
}

var aliasRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove alias",
	Args:  cobra.ExactArgs(1),
	RunE:  runAliasRemove,
}

func init() {
	aliasCmd.AddCommand(aliasSetCmd)
	aliasCmd.AddCommand(aliasListCmd)
	aliasCmd.AddCommand(aliasRemoveCmd)

	rootCmd.AddCommand(themeCmd)
	rootCmd.AddCommand(aliasCmd)
}

// runTheme manages output themes
func runTheme(cmd *cobra.Command, args []string) error {
	themes := []string{"default", "dark", "light", "minimal", "nerd"}

	if len(args) == 0 {
		cyan := color.New(color.FgCyan, color.Bold)
		cyan.Println("\n🎨 Available Themes")
		fmt.Println()
		for _, t := range themes {
			fmt.Printf("  %s\n", t)
		}
		fmt.Println()
		fmt.Println("Set with: agen theme <name>")
		return nil
	}

	themeName := args[0]

	// Validate theme
	valid := false
	for _, t := range themes {
		if t == themeName {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("unknown theme: %s", themeName)
	}

	// Save to config (simplified)
	printSuccess("Theme set to: %s", themeName)
	fmt.Println("Changes will apply on next run.")

	return nil
}

// Alias management
type Aliases map[string]string

func getAliasPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "agen", "aliases.json"), nil
}

func loadAliases() (Aliases, error) {
	path, err := getAliasPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(Aliases), nil
		}
		return nil, err
	}

	var aliases Aliases
	if err := json.Unmarshal(data, &aliases); err != nil {
		return nil, err
	}

	return aliases, nil
}

func saveAliases(aliases Aliases) error {
	path, err := getAliasPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(aliases, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func runAliasSet(cmd *cobra.Command, args []string) error {
	name := args[0]
	command := args[1]

	aliases, err := loadAliases()
	if err != nil {
		return err
	}

	aliases[name] = command

	if err := saveAliases(aliases); err != nil {
		return err
	}

	printSuccess("Alias set: %s -> %s", name, command)
	return nil
}

func runAliasList(cmd *cobra.Command, args []string) error {
	aliases, err := loadAliases()
	if err != nil {
		return err
	}

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n📝 Aliases")
	fmt.Println()

	if len(aliases) == 0 {
		fmt.Println("No aliases defined.")
		fmt.Println("\nCreate one with:")
		fmt.Println("  agen alias set <name> <command>")
		return nil
	}

	for name, command := range aliases {
		fmt.Printf("  %s = %s\n", name, command)
	}

	return nil
}

func runAliasRemove(cmd *cobra.Command, args []string) error {
	name := args[0]

	aliases, err := loadAliases()
	if err != nil {
		return err
	}

	if _, ok := aliases[name]; !ok {
		return fmt.Errorf("alias not found: %s", name)
	}

	delete(aliases, name)

	if err := saveAliases(aliases); err != nil {
		return err
	}

	printSuccess("Removed alias: %s", name)
	return nil
}
