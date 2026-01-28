// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// AI CLI commands: suggest, explain, compose

package cli

import (
	"fmt"
	"os"

	"github.com/eshanized/agen/internal/ai"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var suggestCmd = &cobra.Command{
	Use:   "suggest [path]",
	Short: "Suggest agents for project",
	Long: `Analyze the project and suggest appropriate agents and skills.

How it works:
- Scans project files and structure
- Detects languages, frameworks, and tools
- Recommends agents based on project type
- Suggests skills for common patterns

Examples:
  agen suggest           # Analyze current directory
  agen suggest ./myapp   # Analyze specific path`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSuggest,
}

var explainCmd = &cobra.Command{
	Use:   "explain <name>",
	Short: "Explain agent or skill",
	Long: `Show detailed explanation of an agent or skill.

Provides:
- Full description
- Associated skills (for agents)
- Complete content/rules
- Usage recommendations

Examples:
  agen explain frontend-specialist
  agen explain clean-code`,
	Args: cobra.ExactArgs(1),
	RunE: runExplain,
}

var composeCmd = &cobra.Command{
	Use:   "compose <name>",
	Short: "Compose custom agent",
	Long: `Create a custom agent by combining existing ones.

Generates a new agent that merges skills and rules
from multiple base agents.

Examples:
  agen compose fullstack --from frontend-specialist,backend-specialist
  agen compose security-dev --from security-auditor,debugger`,
	Args: cobra.ExactArgs(1),
	RunE: runCompose,
}

func init() {
	suggestCmd.Flags().Bool("json", false, "output as JSON")
	suggestCmd.Flags().Int("top", 10, "number of suggestions to show")

	explainCmd.Flags().String("type", "auto", "type (agent, skill, auto)")

	composeCmd.Flags().StringSlice("from", []string{}, "base agents to compose from")
	composeCmd.Flags().StringP("description", "d", "", "agent description")
	composeCmd.Flags().StringP("output", "o", "", "output file")

	rootCmd.AddCommand(suggestCmd)
	rootCmd.AddCommand(explainCmd)
	rootCmd.AddCommand(composeCmd)
}

func runSuggest(cmd *cobra.Command, args []string) error {
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	top, _ := cmd.Flags().GetInt("top")

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n🤖 AGEN Suggest")
	fmt.Printf("Analyzing: %s\n\n", targetDir)

	suggester, err := ai.NewSuggester()
	if err != nil {
		return err
	}

	suggestions, err := suggester.Suggest(targetDir)
	if err != nil {
		return err
	}

	if len(suggestions) == 0 {
		fmt.Println("No specific suggestions for this project.")
		return nil
	}

	// Limit to top N
	if len(suggestions) > top {
		suggestions = suggestions[:top]
	}

	fmt.Println("Recommended:")
	fmt.Println()

	green := color.New(color.FgGreen)
	dim := color.New(color.Faint)

	for i, s := range suggestions {
		scoreColor := green
		if s.Score < 0.7 {
			scoreColor = color.New(color.FgYellow)
		}

		fmt.Printf("  %d. ", i+1)
		scoreColor.Printf("%.0f%% ", s.Score*100)
		fmt.Printf("%s ", s.Name)
		dim.Printf("(%s)\n", s.Type)
		dim.Printf("     %s\n", s.Reason)
		fmt.Println()
	}

	fmt.Println("Install all recommended:")
	fmt.Printf("  agen init --agents %s\n",
		joinNames(suggestions[:min(5, len(suggestions))]))

	return nil
}

func runExplain(cmd *cobra.Command, args []string) error {
	name := args[0]
	typeHint, _ := cmd.Flags().GetString("type")

	suggester, err := ai.NewSuggester()
	if err != nil {
		return err
	}

	var explanation string

	// Try to detect type if auto
	if typeHint == "auto" {
		explanation, err = suggester.ExplainAgent(name)
		if err != nil {
			explanation, err = suggester.ExplainSkill(name)
		}
	} else if typeHint == "agent" {
		explanation, err = suggester.ExplainAgent(name)
	} else {
		explanation, err = suggester.ExplainSkill(name)
	}

	if err != nil {
		return fmt.Errorf("not found: %s", name)
	}

	fmt.Println(explanation)
	return nil
}

func runCompose(cmd *cobra.Command, args []string) error {
	name := args[0]
	baseAgents, _ := cmd.Flags().GetStringSlice("from")
	description, _ := cmd.Flags().GetString("description")
	output, _ := cmd.Flags().GetString("output")

	if len(baseAgents) == 0 {
		return fmt.Errorf("--from flag required (e.g., --from agent1,agent2)")
	}

	if description == "" {
		description = fmt.Sprintf("Custom agent composed from %v", baseAgents)
	}

	suggester, err := ai.NewSuggester()
	if err != nil {
		return err
	}

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n🔨 Composing Agent")
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("From: %v\n\n", baseAgents)

	composed, err := suggester.Compose(name, description, baseAgents)
	if err != nil {
		return err
	}

	if output != "" {
		if err := os.WriteFile(output, []byte(composed.Content), 0644); err != nil {
			return err
		}
		printSuccess("Created: %s", output)
	} else {
		fmt.Println("--- Generated Agent ---")
		fmt.Println(composed.Content)
		fmt.Println("--- End ---")
		fmt.Println()
		fmt.Println("Save with:")
		fmt.Printf("  agen compose %s --from %s -o %s.md\n",
			name, joinStrings(baseAgents), name)
	}

	return nil
}

// Helper functions

func joinNames(suggestions []ai.Suggestion) string {
	names := make([]string, len(suggestions))
	for i, s := range suggestions {
		names[i] = s.Name
	}
	return joinStrings(names)
}

func joinStrings(s []string) string {
	result := ""
	for i, str := range s {
		if i > 0 {
			result += ","
		}
		result += str
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
