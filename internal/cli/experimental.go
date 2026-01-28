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

// Theme features removed for production hardening

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

	rootCmd.AddCommand(aliasCmd)
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
