// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// A cross-platform CLI tool for managing AI agent templates

package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/eshanized/agen/internal/templates"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// listCmd shows available agents, skills, and workflows.
// Super useful when you want to know what's available before running init.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available agents, skills, and workflows",
	Long: `List all available agents, skills, and workflows that can be installed.

By default, shows everything. Use flags to filter the output.

Examples:
  agen list              # List everything
  agen list --agents     # Only show agents
  agen list --skills     # Only show skills
  agen list --workflows  # Only show workflows`,
	RunE: runList,
}

func init() {
	listCmd.Flags().BoolP("agents", "a", false, "only show agents")
	listCmd.Flags().BoolP("skills", "s", false, "only show skills")
	listCmd.Flags().BoolP("workflows", "w", false, "only show workflows")
	listCmd.Flags().Bool("json", false, "output as JSON")
}

// runList is the main logic for the list command.
//
// How it works:
// 1. Load templates from embedded storage
// 2. Check which flags are set to filter output
// 3. If no flags, show everything
// 4. Print in a nice formatted table-ish layout
//
// why load embedded? because we want to show what's AVAILABLE to install,
// not what's already installed. For installed stuff, use "agen status"
func runList(cmd *cobra.Command, args []string) error {
	// load templates
	tmpl, err := templates.LoadEmbedded()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	showAgents, _ := cmd.Flags().GetBool("agents")
	showSkills, _ := cmd.Flags().GetBool("skills")
	showWorkflows, _ := cmd.Flags().GetBool("workflows")

	// if no specific flags, show everything
	showAll := !showAgents && !showSkills && !showWorkflows

	cyan := color.New(color.FgCyan, color.Bold)
	dim := color.New(color.Faint)

	if showAll || showAgents {
		fmt.Println()
		cyan.Println("📦 AGENTS")
		dim.Println("Specialist AI personas for different domains")
		fmt.Println()

		// Sort agents by name for consistent output
		names := make([]string, 0, len(tmpl.Agents))
		for name := range tmpl.Agents {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			agent := tmpl.Agents[name]
			fmt.Printf("  %-25s %s\n", color.GreenString(name), agent.Description)
		}
		fmt.Printf("\n  Total: %d agents\n", len(tmpl.Agents))
	}

	if showAll || showSkills {
		fmt.Println()
		cyan.Println("🧩 SKILLS")
		dim.Println("Domain-specific knowledge modules")
		fmt.Println()

		names := make([]string, 0, len(tmpl.Skills))
		for name := range tmpl.Skills {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			skill := tmpl.Skills[name]
			fmt.Printf("  %-25s %s\n", color.BlueString(name), skill.Description)
		}
		fmt.Printf("\n  Total: %d skills\n", len(tmpl.Skills))
	}

	if showAll || showWorkflows {
		fmt.Println()
		cyan.Println("🔄 WORKFLOWS")
		dim.Println("Slash command procedures")
		fmt.Println()

		names := make([]string, 0, len(tmpl.Workflows))
		for name := range tmpl.Workflows {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			workflow := tmpl.Workflows[name]
			// workflows usually have a leading slash in their name
			displayName := name
			if !strings.HasPrefix(name, "/") {
				displayName = "/" + name
			}
			fmt.Printf("  %-25s %s\n", color.MagentaString(displayName), workflow.Description)
		}
		fmt.Printf("\n  Total: %d workflows\n", len(tmpl.Workflows))
	}

	fmt.Println()
	return nil
}
