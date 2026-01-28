// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Remote repository management and hub commands

package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// RemoteRepo represents a custom template source
type RemoteRepo struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Type   string `json:"type"` // git, http
	Branch string `json:"branch,omitempty"`
}

// remoteCmd manages remote template repositories
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Manage remote repositories",
	Long: `Manage remote template repositories.

Add custom sources for agents, skills, and workflows.

Examples:
  agen remote add company https://github.com/company/agents
  agen remote list
  agen remote remove company`,
}

var remoteAddCmd = &cobra.Command{
	Use:   "add <name> <url>",
	Short: "Add remote repository",
	Args:  cobra.ExactArgs(2),
	RunE:  runRemoteAdd,
}

var remoteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List remote repositories",
	RunE:  runRemoteList,
}

var remoteRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove remote repository",
	Args:  cobra.ExactArgs(1),
	RunE:  runRemoteRemove,
}

// Hub commands removed due to being mock/stubs

func init() {
	remoteAddCmd.Flags().String("branch", "main", "branch to use")
	remoteAddCmd.Flags().String("type", "git", "repository type (git, http)")

	remoteCmd.AddCommand(remoteAddCmd)
	remoteCmd.AddCommand(remoteListCmd)
	remoteCmd.AddCommand(remoteRemoveCmd)

	rootCmd.AddCommand(remoteCmd)
}

func getRemotesPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "agen", "remotes.json"), nil
}

func loadRemotes() ([]RemoteRepo, error) {
	path, err := getRemotesPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []RemoteRepo{}, nil
		}
		return nil, err
	}

	var remotes []RemoteRepo
	if err := json.Unmarshal(data, &remotes); err != nil {
		return nil, err
	}

	return remotes, nil
}

func saveRemotes(remotes []RemoteRepo) error {
	path, err := getRemotesPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(remotes, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func runRemoteAdd(cmd *cobra.Command, args []string) error {
	name := args[0]
	url := args[1]
	branch, _ := cmd.Flags().GetString("branch")
	repoType, _ := cmd.Flags().GetString("type")

	remotes, err := loadRemotes()
	if err != nil {
		return err
	}

	// Check if exists
	for _, r := range remotes {
		if r.Name == name {
			return fmt.Errorf("remote already exists: %s", name)
		}
	}

	remotes = append(remotes, RemoteRepo{
		Name:   name,
		URL:    url,
		Type:   repoType,
		Branch: branch,
	})

	if err := saveRemotes(remotes); err != nil {
		return err
	}

	printSuccess("Added remote: %s -> %s", name, url)
	return nil
}

func runRemoteList(cmd *cobra.Command, args []string) error {
	remotes, err := loadRemotes()
	if err != nil {
		return err
	}

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n🌐 Remote Repositories")
	fmt.Println()

	if len(remotes) == 0 {
		fmt.Println("No remotes configured.")
		fmt.Println("\nAdd a remote with:")
		fmt.Println("  agen remote add <name> <url>")
		return nil
	}

	for _, r := range remotes {
		fmt.Printf("  %s\n", r.Name)
		fmt.Printf("    URL: %s\n", r.URL)
		if r.Branch != "" && r.Branch != "main" {
			fmt.Printf("    Branch: %s\n", r.Branch)
		}
	}

	return nil
}

func runRemoteRemove(cmd *cobra.Command, args []string) error {
	name := args[0]

	remotes, err := loadRemotes()
	if err != nil {
		return err
	}

	found := false
	filtered := make([]RemoteRepo, 0)
	for _, r := range remotes {
		if r.Name == name {
			found = true
		} else {
			filtered = append(filtered, r)
		}
	}

	if !found {
		return fmt.Errorf("remote not found: %s", name)
	}

	if err := saveRemotes(filtered); err != nil {
		return err
	}

	printSuccess("Removed remote: %s", name)
	return nil
}
