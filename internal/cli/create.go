// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Create command implementation

package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/eshanized/agen/internal/templates"
	"github.com/eshanized/agen/internal/tui"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Interactively create a new agent",
	Long: `Launch an interactive wizard to create a new agent template.

Collects:
- Agent Name
- Role Description
- Required Skills

Generates a standard Markdown file with frontmatter headers.

Examples:
  agen create
  agen create --output ./my-agents/`,
	RunE: runCreate,
}

func init() {
	createCmd.Flags().StringP("output", "o", ".", "Directory to save the agent file")
	rootCmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
	outputDir, _ := cmd.Flags().GetString("output")

	// 1. Load templates (needed for skill list)
	tmpl, err := templates.LoadEmbedded()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	// 2. Run TUI Wizard
	result, err := tui.RunCreator(tmpl)
	if err != nil {
		return fmt.Errorf("wizard failed: %w", err)
	}

	if result.Cancelled {
		color.Yellow("Operation cancelled.")
		return nil
	}

	// 3. Generate Content
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("name: %s\n", result.Name))
	sb.WriteString(fmt.Sprintf("description: %s\n", result.Description))

	if len(result.Skills) > 0 {
		sb.WriteString("skills: " + strings.Join(result.Skills, ", ") + "\n")
	}
	sb.WriteString("---\n\n")

	sb.WriteString(fmt.Sprintf("# %s\n\n", result.Name))
	sb.WriteString(fmt.Sprintf("> %s\n\n", result.Description))
	sb.WriteString("## Core Responsibilities\n\n")
	sb.WriteString("- [ ] responsibility 1\n")
	sb.WriteString("- [ ] responsibility 2\n\n")
	sb.WriteString("## Guidelines\n\n")
	sb.WriteString("1. First guideline\n")
	sb.WriteString("2. Second guideline\n")

	// 4. Save to file
	filename := strings.ToLower(strings.ReplaceAll(result.Name, " ", "-")) + ".md"
	if !strings.HasSuffix(filename, ".md") {
		filename += ".md"
	}

	// Ensure output dir exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	path := fmt.Sprintf("%s/%s", strings.TrimRight(outputDir, "/"), filename)
	if err := os.WriteFile(path, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	printSuccess("Created agent: %s", path)
	return nil
}
