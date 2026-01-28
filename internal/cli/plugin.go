// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Plugin CLI commands

package cli

import (
	"fmt"

	"github.com/eshanized/agen/internal/plugin"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// pluginCmd is the parent command for plugin management
var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage plugins",
	Long: `Manage AGEN plugins for custom agents, skills, and workflows.

Plugins extend AGEN with community-created or custom templates.

Examples:
  agen plugin list                          # List installed plugins
  agen plugin install github.com/user/repo  # Install from GitHub
  agen plugin create my-agent --type agent  # Create new plugin`,
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install <source>",
	Short: "Install a plugin",
	Long: `Install a plugin from various sources.

Supported sources:
- GitHub: github.com/user/repo[@version]
- Local: ./path/to/plugin
- URL: https://example.com/plugin.zip

Examples:
  agen plugin install github.com/eshanized/agen-plugins
  agen plugin install ./my-local-plugin
  agen plugin install github.com/user/repo@v1.0.0`,
	Args: cobra.ExactArgs(1),
	RunE: runPluginInstall,
}

var pluginUninstallCmd = &cobra.Command{
	Use:   "uninstall <name>",
	Short: "Uninstall a plugin",
	Args:  cobra.ExactArgs(1),
	RunE:  runPluginUninstall,
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed plugins",
	RunE:  runPluginList,
}

var pluginCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new plugin",
	Long: `Create a new plugin project with boilerplate files.

Types:
- agent: Single agent plugin
- skill: Single skill plugin
- bundle: Multiple agents/skills/workflows

Examples:
  agen plugin create my-agent --type agent
  agen plugin create my-toolkit --type bundle`,
	Args: cobra.ExactArgs(1),
	RunE: runPluginCreate,
}

var pluginInfoCmd = &cobra.Command{
	Use:   "info <name>",
	Short: "Show plugin details",
	Args:  cobra.ExactArgs(1),
	RunE:  runPluginInfo,
}

func init() {
	pluginCreateCmd.Flags().String("type", "bundle", "plugin type (agent, skill, workflow, bundle)")

	pluginCmd.AddCommand(pluginInstallCmd)
	pluginCmd.AddCommand(pluginUninstallCmd)
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginCreateCmd)
	pluginCmd.AddCommand(pluginInfoCmd)

	rootCmd.AddCommand(pluginCmd)
}

func runPluginInstall(cmd *cobra.Command, args []string) error {
	source := args[0]

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n🔌 Installing Plugin")
	fmt.Printf("Source: %s\n\n", source)

	manager, err := plugin.NewManager()
	if err != nil {
		return fmt.Errorf("failed to initialize plugin manager: %w", err)
	}

	p, err := manager.Install(source)
	if err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	printSuccess("Installed: %s v%s", p.Name, p.Version)
	if len(p.Agents) > 0 {
		fmt.Printf("  Agents: %v\n", p.Agents)
	}
	if len(p.Skills) > 0 {
		fmt.Printf("  Skills: %v\n", p.Skills)
	}
	if len(p.Workflows) > 0 {
		fmt.Printf("  Workflows: %v\n", p.Workflows)
	}

	return nil
}

func runPluginUninstall(cmd *cobra.Command, args []string) error {
	name := args[0]

	manager, err := plugin.NewManager()
	if err != nil {
		return err
	}

	if err := manager.Uninstall(name); err != nil {
		return err
	}

	printSuccess("Uninstalled: %s", name)
	return nil
}

func runPluginList(cmd *cobra.Command, args []string) error {
	manager, err := plugin.NewManager()
	if err != nil {
		return err
	}

	plugins := manager.List()

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n🔌 Installed Plugins")
	fmt.Println()

	if len(plugins) == 0 {
		fmt.Println("No plugins installed.")
		fmt.Println("\nInstall a plugin with:")
		fmt.Println("  agen plugin install github.com/user/repo")
		return nil
	}

	for _, p := range plugins {
		fmt.Printf("  %s v%s\n", p.Name, p.Version)
		fmt.Printf("    Type: %s\n", p.Type)
		if p.Description != "" {
			fmt.Printf("    %s\n", p.Description)
		}
	}

	return nil
}

func runPluginCreate(cmd *cobra.Command, args []string) error {
	name := args[0]
	pluginType, _ := cmd.Flags().GetString("type")

	manager, err := plugin.NewManager()
	if err != nil {
		return err
	}

	path, err := manager.Create(name, pluginType)
	if err != nil {
		return err
	}

	printSuccess("Created plugin: %s", path)
	fmt.Println("\nNext steps:")
	fmt.Printf("  1. cd %s\n", path)
	fmt.Println("  2. Edit the template files")
	fmt.Printf("  3. agen plugin install ./%s\n", name)

	return nil
}

func runPluginInfo(cmd *cobra.Command, args []string) error {
	name := args[0]

	manager, err := plugin.NewManager()
	if err != nil {
		return err
	}

	p, err := manager.Get(name)
	if err != nil {
		return err
	}

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Printf("\n🔌 %s\n\n", p.Name)

	fmt.Printf("Version:     %s\n", p.Version)
	fmt.Printf("Type:        %s\n", p.Type)
	fmt.Printf("Author:      %s\n", p.Author)
	fmt.Printf("Source:      %s\n", p.Source)

	if p.Description != "" {
		fmt.Printf("\n%s\n", p.Description)
	}

	if len(p.Agents) > 0 {
		fmt.Printf("\nAgents:\n")
		for _, a := range p.Agents {
			fmt.Printf("  - %s\n", a)
		}
	}

	if len(p.Skills) > 0 {
		fmt.Printf("\nSkills:\n")
		for _, s := range p.Skills {
			fmt.Printf("  - %s\n", s)
		}
	}

	return nil
}
