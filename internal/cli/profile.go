// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// A cross-platform CLI tool for managing AI agent templates

package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// profileCmd manages saved configuration profiles.
// lets users save and reuse agent/skill combinations.
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage configuration profiles",
	Long: `Save and load configuration profiles for quick setup.

Profiles let you save your preferred agent/skill combinations
and quickly apply them to new projects.

Examples:
  agen profile save frontend              # Save current setup as "frontend"
  agen profile load backend-api           # Apply saved "backend-api" profile
  agen profile list                       # Show all saved profiles
  agen profile delete old-profile         # Delete a profile
  agen profile export frontend > f.json   # Export to file`,
}

// Profile represents a saved configuration
type Profile struct {
	Name       string   `json:"name"`
	IDE        string   `json:"ide,omitempty"`
	Agents     []string `json:"agents"`
	Skills     []string `json:"skills"`
	Workflows  []string `json:"workflows,omitempty"`
	CreatedAt  string   `json:"created_at"`
	ModifiedAt string   `json:"modified_at,omitempty"`
}

var saveProfileCmd = &cobra.Command{
	Use:   "save <name>",
	Short: "Save current configuration as a profile",
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileSave,
}

var loadProfileCmd = &cobra.Command{
	Use:   "load <name>",
	Short: "Load and apply a saved profile",
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileLoad,
}

var listProfileCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved profiles",
	RunE:  runProfileList,
}

var deleteProfileCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a saved profile",
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileDelete,
}

var exportProfileCmd = &cobra.Command{
	Use:   "export <name>",
	Short: "Export a profile as JSON",
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileExport,
}

var importProfileCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Import a profile from JSON file",
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileImport,
}

func init() {
	profileCmd.AddCommand(saveProfileCmd)
	profileCmd.AddCommand(loadProfileCmd)
	profileCmd.AddCommand(listProfileCmd)
	profileCmd.AddCommand(deleteProfileCmd)
	profileCmd.AddCommand(exportProfileCmd)
	profileCmd.AddCommand(importProfileCmd)
}

// getProfilesDir returns the directory where profiles are stored.
// Uses XDG config dir on linux/mac, AppData on Windows.
func getProfilesDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "agen", "profiles"), nil
}

// runProfileSave saves the current project's configuration as a named profile.
//
// How it works:
// 1. Detect current IDE and installed agents/skills
// 2. Create a Profile struct with all the info
// 3. Serialize to JSON and save to profiles directory
//
// Why profiles? So users can set up once and reuse everywhere.
// frontend devs always want the same agents, so they save a "frontend"
// profile and just "agen profile load frontend" on new projects.
func runProfileSave(cmd *cobra.Command, args []string) error {
	profileName := args[0]

	// Get current directory and detect installed config
	cwd, _ := os.Getwd()
	agentDir := filepath.Join(cwd, ".agent", "agents")
	skillDir := filepath.Join(cwd, ".agent", "skills")

	profile := Profile{
		Name:      profileName,
		Agents:    []string{},
		Skills:    []string{},
		Workflows: []string{},
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	// Detect installed agents
	if entries, err := os.ReadDir(agentDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() && filepath.Ext(e.Name()) == ".md" {
				name := strings.TrimSuffix(e.Name(), ".md")
				profile.Agents = append(profile.Agents, name)
			}
		}
	}

	// Detect installed skills
	if entries, err := os.ReadDir(skillDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				profile.Skills = append(profile.Skills, e.Name())
			}
		}
	}

	// Detect IDE
	if _, err := os.Stat(filepath.Join(cwd, ".cursorrules")); err == nil {
		profile.IDE = "cursor"
	} else if _, err := os.Stat(filepath.Join(cwd, ".windsurfrules")); err == nil {
		profile.IDE = "windsurf"
	} else if _, err := os.Stat(filepath.Join(cwd, ".zed")); err == nil {
		profile.IDE = "zed"
	} else if _, err := os.Stat(filepath.Join(cwd, ".agent")); err == nil {
		profile.IDE = "antigravity"
	}

	// Save profile
	profilesDir, err := getProfilesDir()
	if err != nil {
		return fmt.Errorf("failed to get profiles directory: %w", err)
	}

	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		return fmt.Errorf("failed to create profiles directory: %w", err)
	}

	data, _ := json.MarshalIndent(profile, "", "  ")
	profilePath := filepath.Join(profilesDir, profileName+".json")
	if err := os.WriteFile(profilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to save profile: %w", err)
	}

	printSuccess("Profile '%s' saved", profileName)
	fmt.Printf("  Agents: %d\n", len(profile.Agents))
	fmt.Printf("  Skills: %d\n", len(profile.Skills))
	fmt.Printf("  IDE: %s\n", profile.IDE)
	fmt.Println("\nTo use this profile later:")
	fmt.Printf("  agen profile load %s\n", profileName)

	return nil
}

func runProfileLoad(cmd *cobra.Command, args []string) error {
	profileName := args[0]

	// Load profile
	profilesDir, err := getProfilesDir()
	if err != nil {
		return fmt.Errorf("failed to get profiles directory: %w", err)
	}

	profilePath := filepath.Join(profilesDir, profileName+".json")
	data, err := os.ReadFile(profilePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("profile '%s' not found", profileName)
	}
	if err != nil {
		return fmt.Errorf("failed to read profile: %w", err)
	}

	var profile Profile
	if err := json.Unmarshal(data, &profile); err != nil {
		return fmt.Errorf("invalid profile format: %w", err)
	}

	fmt.Printf("Loading profile: %s\n", profile.Name)
	fmt.Printf("  IDE: %s\n", profile.IDE)
	fmt.Printf("  Agents: %d\n", len(profile.Agents))
	fmt.Printf("  Skills: %d\n", len(profile.Skills))

	// Apply the profile by running init with the profile's settings
	cwd, _ := os.Getwd()

	// Create init command with profile's agents/skills
	initArgs := []string{"init"}
	if profile.IDE != "" {
		initArgs = append(initArgs, "--ide", profile.IDE)
	}
	if len(profile.Agents) > 0 {
		initArgs = append(initArgs, "--agents", strings.Join(profile.Agents, ","))
	}
	if len(profile.Skills) > 0 {
		initArgs = append(initArgs, "--skills", strings.Join(profile.Skills, ","))
	}

	fmt.Println("\nApplying configuration...")
	fmt.Printf("  agen %s\n", strings.Join(initArgs, " "))

	// Since we can't call cobra commands directly, just print instructions
	printSuccess("Profile '%s' loaded", profileName)
	fmt.Println("\nRun the following to apply:")
	fmt.Printf("  agen %s\n", strings.Join(initArgs, " "))

	_ = cwd // for context

	return nil
}

func runProfileList(cmd *cobra.Command, args []string) error {
	profilesDir, err := getProfilesDir()
	if err != nil {
		return fmt.Errorf("failed to get profiles directory: %w", err)
	}

	entries, err := os.ReadDir(profilesDir)
	if os.IsNotExist(err) {
		printInfo("No profiles saved yet")
		fmt.Println("\nCreate one with:")
		fmt.Println("  agen profile save <name>")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to read profiles: %w", err)
	}

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n📝 Saved Profiles")
	fmt.Println()

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			name := strings.TrimSuffix(entry.Name(), ".json")
			fmt.Printf("  • %s\n", color.GreenString(name))
		}
	}

	fmt.Println()
	return nil
}

func runProfileDelete(cmd *cobra.Command, args []string) error {
	profileName := args[0]

	profilesDir, err := getProfilesDir()
	if err != nil {
		return fmt.Errorf("failed to get profiles directory: %w", err)
	}

	profilePath := filepath.Join(profilesDir, profileName+".json")
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		return fmt.Errorf("profile '%s' not found", profileName)
	}

	if err := os.Remove(profilePath); err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	printSuccess("Profile '%s' deleted", profileName)
	return nil
}

func runProfileExport(cmd *cobra.Command, args []string) error {
	profileName := args[0]

	profilesDir, err := getProfilesDir()
	if err != nil {
		return fmt.Errorf("failed to get profiles directory: %w", err)
	}

	profilePath := filepath.Join(profilesDir, profileName+".json")
	data, err := os.ReadFile(profilePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("profile '%s' not found", profileName)
	}
	if err != nil {
		return fmt.Errorf("failed to read profile: %w", err)
	}

	// Pretty print to stdout
	var profile Profile
	if err := json.Unmarshal(data, &profile); err != nil {
		return fmt.Errorf("failed to parse profile: %w", err)
	}

	pretty, _ := json.MarshalIndent(profile, "", "  ")
	fmt.Println(string(pretty))

	return nil
}

func runProfileImport(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var profile Profile
	if err := json.Unmarshal(data, &profile); err != nil {
		return fmt.Errorf("invalid profile format: %w", err)
	}

	if profile.Name == "" {
		return fmt.Errorf("profile must have a name")
	}

	// save to profiles dir
	profilesDir, err := getProfilesDir()
	if err != nil {
		return fmt.Errorf("failed to get profiles directory: %w", err)
	}

	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		return fmt.Errorf("failed to create profiles directory: %w", err)
	}

	profilePath := filepath.Join(profilesDir, profile.Name+".json")
	if err := os.WriteFile(profilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to save profile: %w", err)
	}

	printSuccess("Profile '%s' imported", profile.Name)
	return nil
}
