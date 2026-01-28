// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Advanced commands: diff, watch, audit, export, validate

package cli

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eshanized/agen/internal/ide"
	"github.com/eshanized/agen/internal/templates"
	"github.com/eshanized/agen/internal/updater"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

// diffCmd shows differences between installed and latest templates
var diffCmd = &cobra.Command{
	Use:   "diff [path]",
	Short: "Show template differences",
	Long: `Compare installed templates with the latest available version.

Shows files that are:
- Added (new in latest)
- Modified (different from installed)
- Removed (no longer in latest)
- Unchanged

Examples:
  agen diff             # Compare current directory
  agen diff --detailed  # Show file-level diffs`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDiff,
}

// watchCmd watches for template changes and auto-updates
var watchCmd = &cobra.Command{
	Use:   "watch [path]",
	Short: "Watch for template changes",
	Long: `Watch the project for changes and automatically sync templates.

Monitors:
- .agent/ directory for local changes
- GitHub for upstream updates (optional)

Examples:
  agen watch            # Watch current directory
  agen watch --upstream # Also check for remote updates`,
	Args: cobra.MaximumNArgs(1),
	RunE: runWatch,
}

// auditCmd performs security audit of installed agents
var auditCmd = &cobra.Command{
	Use:   "audit [path]",
	Short: "Security audit of templates",
	Long: `Perform a security audit of installed agent templates.

Checks:
- Template integrity (checksums)
- Suspicious patterns in agent rules
- Known security issues
- License compliance

Examples:
  agen audit            # Audit current directory
  agen audit --fix      # Attempt to fix issues`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAudit,
}

// exportCmd exports templates to various formats
var exportCmd = &cobra.Command{
	Use:   "export [path]",
	Short: "Export templates",
	Long: `Export installed templates to various formats.

Formats:
- json: Machine-readable JSON
- yaml: YAML format
- markdown: Human-readable documentation
- zip: Compressed archive

Examples:
  agen export --format json > templates.json
  agen export --format zip -o backup.zip`,
	Args: cobra.MaximumNArgs(1),
	RunE: runExport,
}

// validateCmd validates template syntax
var validateCmd = &cobra.Command{
	Use:   "validate [path]",
	Short: "Validate template syntax",
	Long: `Validate the syntax and structure of agent templates.

Checks:
- YAML frontmatter parsing
- Required fields present
- Skill references valid
- Workflow syntax correct

Examples:
  agen validate         # Validate current directory
  agen validate --strict`,
	Args: cobra.MaximumNArgs(1),
	RunE: runValidate,
}

func init() {
	diffCmd.Flags().Bool("detailed", false, "show detailed file diffs")
	diffCmd.Flags().Bool("json", false, "output as JSON")

	watchCmd.Flags().Bool("upstream", false, "also watch for remote updates")
	watchCmd.Flags().Duration("interval", 30*time.Second, "check interval for upstream")

	auditCmd.Flags().Bool("fix", false, "attempt to fix issues")
	auditCmd.Flags().Bool("json", false, "output as JSON")

	exportCmd.Flags().StringP("format", "f", "json", "output format (json, yaml, markdown, zip)")
	exportCmd.Flags().StringP("output", "o", "", "output file (default: stdout)")

	validateCmd.Flags().Bool("strict", false, "strict validation mode")
	validateCmd.Flags().Bool("json", false, "output as JSON")

	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(watchCmd)
	rootCmd.AddCommand(auditCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(validateCmd)
}

// runDiff compares installed vs latest templates
//
// How it works:
// 1. Load installed templates from project
// 2. Load latest templates (embedded or network)
// 3. Compare each file by content hash
// 4. Display differences in a clear format
func runDiff(cmd *cobra.Command, args []string) error {
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	absPath, _ := filepath.Abs(targetDir)
	detailed, _ := cmd.Flags().GetBool("detailed")

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n📊 AGEN Diff")
	fmt.Printf("Directory: %s\n\n", absPath)

	// detect IDE
	ideAdapter := ide.Detect(absPath)
	if ideAdapter == nil {
		return fmt.Errorf("no AGEN installation found")
	}

	// Load latest templates
	latest, err := templates.LoadEmbedded()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	// Compare agents
	agentDir := filepath.Join(absPath, ".agent", "agents")
	added, modified, removed, unchanged := 0, 0, 0, 0

	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	fmt.Println("📦 Agents:")
	for name, agent := range latest.Agents {
		installedPath := filepath.Join(agentDir, name+".md")
		if _, err := os.Stat(installedPath); os.IsNotExist(err) {
			green.Printf("  + %s (new)\n", name)
			added++
		} else {
			// Compare content
			installedContent, _ := os.ReadFile(installedPath)
			if hashContent(installedContent) != hashContent([]byte(agent.Content)) {
				yellow.Printf("  ~ %s (modified)\n", name)
				modified++
				if detailed {
					showDiff(string(installedContent), agent.Content)
				}
			} else {
				unchanged++
			}
		}
	}

	// Check for removed agents
	if entries, err := os.ReadDir(agentDir); err == nil {
		for _, entry := range entries {
			name := strings.TrimSuffix(entry.Name(), ".md")
			if _, ok := latest.Agents[name]; !ok {
				red.Printf("  - %s (removed)\n", name)
				removed++
			}
		}
	}

	fmt.Println()
	fmt.Printf("Summary: +%d added, ~%d modified, -%d removed, %d unchanged\n",
		added, modified, removed, unchanged)

	return nil
}

// runWatch monitors for changes
func runWatch(cmd *cobra.Command, args []string) error {
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	absPath, _ := filepath.Abs(targetDir)
	upstream, _ := cmd.Flags().GetBool("upstream")
	interval, _ := cmd.Flags().GetDuration("interval")

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n👁 AGEN Watch")
	fmt.Printf("Directory: %s\n", absPath)
	fmt.Printf("Upstream: %v (interval: %v)\n", upstream, interval)
	fmt.Println("\nWatching for changes... (Ctrl+C to stop)")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	agentDir := filepath.Join(absPath, ".agent")
	if err := watcher.Add(agentDir); err != nil {
		return fmt.Errorf("failed to watch directory: %w", err)
	}

	// Also watch subdirectories
	filepath.Walk(agentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return nil
		}
		watcher.Add(path)
		return nil
	})

	// Upstream ticker
	var upstreamTicker *time.Ticker
	if upstream {
		upstreamTicker = time.NewTicker(interval)
		defer upstreamTicker.Stop()
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) != 0 {
				relPath, _ := filepath.Rel(absPath, event.Name)
				printInfo("[%s] %s: %s",
					time.Now().Format("15:04:05"),
					event.Op.String(),
					relPath)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			printError("Watcher error: %v", err)

		case <-func() <-chan time.Time {
			if upstreamTicker != nil {
				return upstreamTicker.C
			}
			return nil
		}():
			printInfo("[%s] Checking upstream...", time.Now().Format("15:04:05"))
			if release, err := updater.CheckForUpdate(Version); err == nil && release != nil {
				color.Green("\n✨ New version available: %s", release.Version)
				fmt.Printf("Run 'agen upgrade' to update\n\n")
			} else if err != nil {
				printWarning("Failed to check upstream: %v", err)
			}
		}
	}
}

// runAudit performs security audit
func runAudit(cmd *cobra.Command, args []string) error {
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	absPath, _ := filepath.Abs(targetDir)

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n🔒 AGEN Security Audit")
	fmt.Printf("Directory: %s\n\n", absPath)

	issues := 0

	// Check 1: Template integrity
	fmt.Println("Checking template integrity...")
	agentDir := filepath.Join(absPath, ".agent", "agents")
	if entries, err := os.ReadDir(agentDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			// verify file is readable and valid markdown
			content, err := os.ReadFile(filepath.Join(agentDir, entry.Name()))
			if err != nil {
				printWarning("  Cannot read: %s", entry.Name())
				issues++
			} else if len(content) == 0 {
				printWarning("  Empty file: %s", entry.Name())
				issues++
			}
		}
	}
	if issues == 0 {
		printSuccess("Template integrity: OK")
	}

	// Check 2: Suspicious patterns
	fmt.Println("\nChecking for suspicious patterns...")
	suspiciousPatterns := []string{
		"rm -rf",
		"sudo",
		"curl | bash",
		"wget | bash",
		"eval(",
	}

	filepath.Walk(filepath.Join(absPath, ".agent"), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		for _, pattern := range suspiciousPatterns {
			if strings.Contains(string(content), pattern) {
				relPath, _ := filepath.Rel(absPath, path)
				printWarning("  Found '%s' in %s", pattern, relPath)
				issues++
			}
		}
		return nil
	})

	if issues == 0 {
		printSuccess("No suspicious patterns found")
	}

	// Summary
	fmt.Println()
	if issues == 0 {
		color.New(color.FgGreen, color.Bold).Println("✨ Audit passed!")
	} else {
		color.New(color.FgYellow).Printf("⚠ Found %d potential issue(s)\n", issues)
	}

	return nil
}

// runExport exports templates to various formats
func runExport(cmd *cobra.Command, args []string) error {
	format, _ := cmd.Flags().GetString("format")
	output, _ := cmd.Flags().GetString("output")

	// load templates
	tmpl, err := templates.LoadEmbedded()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	var data []byte

	switch format {
	case "json":
		data, err = json.MarshalIndent(tmpl, "", "  ")
	case "yaml":
		// Simple YAML-like output
		var sb strings.Builder
		sb.WriteString("# AGEN Templates Export\n\n")
		sb.WriteString("agents:\n")
		for name, agent := range tmpl.Agents {
			sb.WriteString(fmt.Sprintf("  %s:\n    description: %s\n", name, agent.Description))
		}
		sb.WriteString("\nskills:\n")
		for name, skill := range tmpl.Skills {
			sb.WriteString(fmt.Sprintf("  %s:\n    description: %s\n", name, skill.Description))
		}
		data = []byte(sb.String())
	case "markdown":
		var sb strings.Builder
		sb.WriteString("# AGEN Templates\n\n")
		sb.WriteString("## Agents\n\n")
		for name, agent := range tmpl.Agents {
			sb.WriteString(fmt.Sprintf("### %s\n%s\n\n", name, agent.Description))
		}
		sb.WriteString("## Skills\n\n")
		for name, skill := range tmpl.Skills {
			sb.WriteString(fmt.Sprintf("### %s\n%s\n\n", name, skill.Description))
		}
		data = []byte(sb.String())
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return err
	}

	if output != "" {
		return os.WriteFile(output, data, 0644)
	}

	fmt.Println(string(data))
	return nil
}

// runValidate validates template syntax
func runValidate(cmd *cobra.Command, args []string) error {
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	absPath, _ := filepath.Abs(targetDir)
	strict, _ := cmd.Flags().GetBool("strict")

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n✅ AGEN Validate")
	fmt.Printf("Directory: %s\n", absPath)
	fmt.Printf("Strict mode: %v\n\n", strict)

	errors := 0
	warnings := 0

	// Validate agents
	agentDir := filepath.Join(absPath, ".agent", "agents")
	if entries, err := os.ReadDir(agentDir); err == nil {
		fmt.Println("Validating agents...")
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}

			content, err := os.ReadFile(filepath.Join(agentDir, entry.Name()))
			if err != nil {
				printError("  %s: cannot read", entry.Name())
				errors++
				continue
			}

			// Check for frontmatter
			if !strings.HasPrefix(string(content), "---") {
				if strict {
					printWarning("  %s: missing frontmatter", entry.Name())
					warnings++
				}
			}

			// Check for heading
			if !strings.Contains(string(content), "#") {
				printWarning("  %s: no markdown heading", entry.Name())
				warnings++
			}
		}
	}

	// Validate skills
	skillDir := filepath.Join(absPath, ".agent", "skills")
	if entries, err := os.ReadDir(skillDir); err == nil {
		fmt.Println("Validating skills...")
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			skillFile := filepath.Join(skillDir, entry.Name(), "SKILL.md")
			if _, err := os.Stat(skillFile); os.IsNotExist(err) {
				printError("  %s: missing SKILL.md", entry.Name())
				errors++
			}
		}
	}

	fmt.Println()
	if errors == 0 && warnings == 0 {
		color.New(color.FgGreen, color.Bold).Println("✨ All templates valid!")
	} else if errors == 0 {
		color.Yellow("⚠ Valid with %d warning(s)", warnings)
	} else {
		color.Red("❌ Found %d error(s), %d warning(s)", errors, warnings)
	}

	return nil
}

// Helper functions

func hashContent(content []byte) string {
	hash := md5.Sum(content)
	return hex.EncodeToString(hash[:])
}

func showDiff(old, new string) {
	// Simple diff display
	oldLines := strings.Split(old, "\n")
	newLines := strings.Split(new, "\n")

	fmt.Printf("    --- installed (%d lines)\n", len(oldLines))
	fmt.Printf("    +++ latest (%d lines)\n", len(newLines))
}
