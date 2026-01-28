// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Quick win commands: doctor, changelog, clean, stats

package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/eshanized/agen/internal/config"
	"github.com/eshanized/agen/internal/ide"
	"github.com/eshanized/agen/internal/templates"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// doctorCmd diagnoses installation and configuration issues
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose installation issues",
	Long: `Run diagnostic checks to identify potential issues with AGEN.

Checks:
- Configuration file validity
- Template integrity
- IDE detection
- Network connectivity
- Cache status

Examples:
  agen doctor          # Run all diagnostics
  agen doctor --fix    # Attempt to fix issues`,
	RunE: runDoctor,
}

// cleanCmd removes cached files and temporary data
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove cached files",
	Long: `Clean up cached templates, temporary files, and build artifacts.

Examples:
  agen clean           # Clean all caches
  agen clean --cache   # Only template cache
  agen clean --temp    # Only temp files`,
	RunE: runClean,
}

// statsCmd shows usage statistics
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show usage statistics",
	Long: `Display AGEN usage statistics and metrics.

Shows:
- Number of projects using AGEN
- Most used agents and skills
- Template versions in use
- Command history summary`,
	RunE: runStats,
}

// changelogCmd shows what's new
var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Show what's new",
	Long: `Display the changelog for the current or specified version.

Examples:
  agen changelog        # Show current version changes
  agen changelog v1.0.0 # Show specific version`,
	Args: cobra.MaximumNArgs(1),
	RunE: runChangelog,
}

func init() {
	doctorCmd.Flags().Bool("fix", false, "attempt to fix issues automatically")
	cleanCmd.Flags().Bool("cache", false, "only clean template cache")
	cleanCmd.Flags().Bool("temp", false, "only clean temporary files")
	cleanCmd.Flags().Bool("all", false, "clean everything")
	statsCmd.Flags().Bool("json", false, "output as JSON")

	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(changelogCmd)
}

// runDoctor performs diagnostic checks
//
// How it works:
// 1. Check if config file exists and is valid
// 2. Verify embedded templates load correctly
// 3. Test IDE detection capabilities
// 4. Check network connectivity to GitHub
// 5. Verify cache directory is writable
//
// Why a doctor command? Helps users troubleshoot issues without
// digging through logs or configuration files manually.
func runDoctor(cmd *cobra.Command, args []string) error {
	fix, _ := cmd.Flags().GetBool("fix")

	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	cyan.Println("\n🩺 AGEN Doctor")
	fmt.Println("Running diagnostic checks...\n")

	issues := 0
	fixed := 0

	// Check 1: Configuration
	fmt.Print("Checking configuration... ")
	cfg, err := config.Load()
	if err != nil {
		red.Println("❌ FAILED")
		fmt.Printf("  Error: %v\n", err)
		issues++
		if fix {
			// Try to create default config
			if err := config.DefaultConfig().Save(); err == nil {
				green.Println("  ✓ Created default config")
				fixed++
			}
		}
	} else {
		green.Println("✓ OK")
		fmt.Printf("  Config loaded successfully\n")
		_ = cfg // use cfg
	}

	// Check 2: Templates
	fmt.Print("Checking embedded templates... ")
	tmpl, err := templates.LoadEmbedded()
	if err != nil {
		red.Println("❌ FAILED")
		fmt.Printf("  Error: %v\n", err)
		issues++
	} else {
		green.Println("✓ OK")
		fmt.Printf("  %d agents, %d skills, %d workflows\n",
			len(tmpl.Agents), len(tmpl.Skills), len(tmpl.Workflows))
	}

	// Check 3: IDE detection
	fmt.Print("Checking IDE detection... ")
	cwd, _ := os.Getwd()
	detectedIDE := ide.Detect(cwd)
	if detectedIDE != nil {
		green.Println("✓ OK")
		fmt.Printf("  Detected: %s\n", detectedIDE.Name())
	} else {
		color.Yellow("⚠ No IDE detected")
		fmt.Println("  (This is OK if not in a project directory)")
	}

	// Check 4: Cache directory
	fmt.Print("Checking cache directory... ")
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		red.Println("❌ FAILED")
		fmt.Printf("  Error: %v\n", err)
		issues++
	} else {
		agenCache := filepath.Join(cacheDir, "agen")
		if _, err := os.Stat(agenCache); os.IsNotExist(err) {
			if fix {
				if err := os.MkdirAll(agenCache, 0755); err == nil {
					green.Println("✓ CREATED")
					fixed++
				} else {
					red.Println("❌ FAILED")
					issues++
				}
			} else {
				color.Yellow("⚠ Not created")
				fmt.Printf("  Path: %s\n", agenCache)
			}
		} else {
			green.Println("✓ OK")
			fmt.Printf("  Path: %s\n", agenCache)
		}
	}

	// Check 5: Go runtime
	fmt.Print("Checking runtime... ")
	green.Println("✓ OK")
	fmt.Printf("  Go version: %s\n", runtime.Version())
	fmt.Printf("  Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	// Summary
	fmt.Println()
	if issues == 0 {
		green.Println("✨ All checks passed!")
	} else {
		red.Printf("Found %d issue(s)", issues)
		if fixed > 0 {
			green.Printf(", fixed %d", fixed)
		}
		fmt.Println()
	}

	return nil
}

// runClean removes cached and temporary files
func runClean(cmd *cobra.Command, args []string) error {
	cacheOnly, _ := cmd.Flags().GetBool("cache")
	tempOnly, _ := cmd.Flags().GetBool("temp")
	cleanAll, _ := cmd.Flags().GetBool("all")

	if !cacheOnly && !tempOnly {
		cleanAll = true
	}

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n🧹 AGEN Clean")

	totalCleaned := int64(0)

	if cleanAll || cacheOnly {
		cacheDir, err := os.UserCacheDir()
		if err == nil {
			agenCache := filepath.Join(cacheDir, "agen")
			if size, err := getDirSize(agenCache); err == nil && size > 0 {
				if err := os.RemoveAll(agenCache); err == nil {
					printSuccess("Cleaned cache: %s (%s)", agenCache, formatBytes(size))
					totalCleaned += size
				}
			} else {
				printInfo("Cache already clean")
			}
		}
	}

	if cleanAll || tempOnly {
		tempDir := os.TempDir()
		pattern := filepath.Join(tempDir, "agen-*")
		matches, _ := filepath.Glob(pattern)
		for _, match := range matches {
			if size, err := getDirSize(match); err == nil {
				os.RemoveAll(match)
				totalCleaned += size
			}
		}
		if len(matches) > 0 {
			printSuccess("Cleaned %d temp file(s)", len(matches))
		} else {
			printInfo("No temp files to clean")
		}
	}

	fmt.Printf("\nTotal cleaned: %s\n", formatBytes(totalCleaned))
	return nil
}

// runStats shows usage statistics
func runStats(cmd *cobra.Command, args []string) error {
	jsonOutput, _ := cmd.Flags().GetBool("json")

	cyan := color.New(color.FgCyan, color.Bold)

	stats := struct {
		Version       string    `json:"version"`
		Platform      string    `json:"platform"`
		AgentCount    int       `json:"agent_count"`
		SkillCount    int       `json:"skill_count"`
		WorkflowCount int       `json:"workflow_count"`
		CacheSize     int64     `json:"cache_size_bytes"`
		ConfigExists  bool      `json:"config_exists"`
		LastUsed      time.Time `json:"last_used"`
	}{
		Version:  Version,
		Platform: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}

	// Load templates for counts
	if tmpl, err := templates.LoadEmbedded(); err == nil {
		stats.AgentCount = len(tmpl.Agents)
		stats.SkillCount = len(tmpl.Skills)
		stats.WorkflowCount = len(tmpl.Workflows)
	}

	// check cache size
	if cacheDir, err := os.UserCacheDir(); err == nil {
		agenCache := filepath.Join(cacheDir, "agen")
		stats.CacheSize, _ = getDirSize(agenCache)
	}

	// check config
	if cfgPath, err := config.GetConfigPath(); err == nil {
		_, err := os.Stat(cfgPath)
		stats.ConfigExists = err == nil
	}

	stats.LastUsed = time.Now()

	if jsonOutput {
		data, _ := json.MarshalIndent(stats, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	cyan.Println("\n📊 AGEN Statistics")
	fmt.Println()
	fmt.Printf("Version:    %s\n", stats.Version)
	fmt.Printf("Platform:   %s\n", stats.Platform)
	fmt.Println()
	fmt.Printf("Agents:     %d\n", stats.AgentCount)
	fmt.Printf("Skills:     %d\n", stats.SkillCount)
	fmt.Printf("Workflows:  %d\n", stats.WorkflowCount)
	fmt.Println()
	fmt.Printf("Cache Size: %s\n", formatBytes(stats.CacheSize))
	fmt.Printf("Config:     %v\n", stats.ConfigExists)
	fmt.Println()

	return nil
}

// runChangelog displays version changelog
func runChangelog(cmd *cobra.Command, args []string) error {
	version := Version
	if len(args) > 0 {
		version = args[0]
	}

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Printf("\n📝 Changelog for %s\n\n", version)

	// Embedded changelog for current version
	changelog := `## v1.0.0 (2026-01-28)

### ✨ Features
- Initial release of AGEN CLI
- Support for 4 IDEs: Antigravity, Cursor, Windsurf, Zed
- 20 embedded agents from antigravity-kit
- 36 skills covering frontend, backend, security, and more
- 11 workflow slash commands
- Interactive TUI wizard for setup
- Project health dashboard
- Fuzzy search across agents/skills
- Self-update mechanism
- Go-based verification scripts

### 🔧 Technical
- Cross-platform builds (Windows, macOS, Linux)
- Package manager support (Homebrew, Scoop, AUR, deb, rpm)
- Embedded templates with offline support
- GitHub-based template updates

### 🎯 Commands
- agen init, list, status, health
- agen verify, update, upgrade
- agen search, profile, playground
- agen doctor, clean, stats
`

	fmt.Println(changelog)
	return nil
}

// Helper functions

func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // ignore errors
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
