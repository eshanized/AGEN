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
	"github.com/sahilm/fuzzy"
	"github.com/spf13/cobra"
)

// searchCmd provides fuzzy search across agents, skills, and workflows.
// super helpful when you kinda know what you want but not the exact name.
var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search agents, skills, and workflows",
	Long: `Fuzzy search across all available agents, skills, and workflows.

Examples:
  agen search security     # Find security-related items
  agen search front        # Find frontend-related items
  agen search test         # Find testing-related items`,
	Args: cobra.ExactArgs(1),
	RunE: runSearch,
}

func init() {
	searchCmd.Flags().IntP("limit", "n", 10, "maximum number of results to show")
}

// searchItem represents something that can be searched
type searchItem struct {
	Name        string
	Description string
	Type        string // "agent", "skill", or "workflow"
}

// runSearch implements fuzzy search across all templates.
//
// How it works:
// 1. Load all agents, skills, and workflows
// 2. Create a searchable list with names + descriptions
// 3. Use fuzzy matching to find relevant items
// 4. Sort by match score (best matches first)
// 5. Display with highlighting
//
// Why fuzzy search? Users often don't remember exact names.
// "sec" should find "security-auditor", "frontend" should find
// "frontend-specialist", etc. Makes the tool much more usable.
func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]
	limit, _ := cmd.Flags().GetInt("limit")

	// Load templates
	tmpl, err := templates.LoadEmbedded()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	// Build searchable items list
	var items []searchItem

	for name, agent := range tmpl.Agents {
		items = append(items, searchItem{
			Name:        name,
			Description: agent.Description,
			Type:        "agent",
		})
	}

	for name, skill := range tmpl.Skills {
		items = append(items, searchItem{
			Name:        name,
			Description: skill.Description,
			Type:        "skill",
		})
	}

	for name, workflow := range tmpl.Workflows {
		displayName := name
		if !strings.HasPrefix(name, "/") {
			displayName = "/" + name
		}
		items = append(items, searchItem{
			Name:        displayName,
			Description: workflow.Description,
			Type:        "workflow",
		})
	}

	// Create search strings (name + description for better matching)
	searchStrings := make([]string, len(items))
	for i, item := range items {
		searchStrings[i] = item.Name + " " + item.Description
	}

	// Do fuzzy matching
	matches := fuzzy.Find(query, searchStrings)

	if len(matches) == 0 {
		printWarning("No matches found for '%s'", query)
		fmt.Println("\nTry:")
		fmt.Println("  agen list --agents    # See all available agents")
		fmt.Println("  agen list --skills    # See all available skills")
		return nil
	}

	// Sort by score (higher is better)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})

	// Limit results
	if len(matches) > limit {
		matches = matches[:limit]
	}

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Printf("\n🔍 Search results for '%s'\n\n", query)

	for _, match := range matches {
		item := items[match.Index]

		// Color based on type
		var typeColor *color.Color
		var typeIcon string
		switch item.Type {
		case "agent":
			typeColor = color.New(color.FgGreen)
			typeIcon = "📦"
		case "skill":
			typeColor = color.New(color.FgBlue)
			typeIcon = "🧩"
		case "workflow":
			typeColor = color.New(color.FgMagenta)
			typeIcon = "🔄"
		}

		typeStr := typeColor.Sprintf("[%s]", item.Type)
		fmt.Printf("%s %s %s\n", typeIcon, typeStr, color.New(color.Bold).Sprint(item.Name))
		fmt.Printf("   %s\n\n", color.New(color.Faint).Sprint(item.Description))
	}

	if len(matches) == limit {
		printInfo("Showing top %d results. Use -n to show more.", limit)
	}

	return nil
}
