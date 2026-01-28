// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// A cross-platform CLI tool for managing AI agent templates

package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/eshanized/agen/internal/ide"
	"github.com/eshanized/agen/internal/templates"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// healthCmd displays a project health dashboard.
// shows template status, compatibility score, and suggestions.
var healthCmd = &cobra.Command{
	Use:   "health [path]",
	Short: "Show project health dashboard",
	Long: `Display a comprehensive project health dashboard.

Shows:
- Template version status (up-to-date, outdated, modified)
- IDE compatibility score
- Missing recommended agents for the project type
- Quick fix suggestions

Examples:
  agen health              # Check current directory
  agen health /path/to/proj`,
	Args: cobra.MaximumNArgs(1),
	RunE: runHealth,
}

// runHealth shows the project health dashboard.
//
// How it works:
// 1. Detect IDE and installed templates
// 2. Analyze project type (web, mobile, backend, etc.)
// 3. Compare installed agents with recommended ones for that project type
// 4. Calculate a "health score" based on various factors
// 5. Provide actionable suggestions
//
// Why a health dashboard? Helps users understand if they're missing
// useful agents for their project. A backend project without security-auditor
// is probably not ideal, for example.
func runHealth(cmd *cobra.Command, args []string) error {
	verbose := checkVerbose(cmd)

	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", absPath)
	}

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n📊 AGEN Health Dashboard")
	fmt.Printf("Directory: %s\n\n", absPath)

	// Step 1: Detect IDE
	ideAdapter := ide.Detect(absPath)
	if ideAdapter == nil {
		printWarning("No AGEN installation found")
		fmt.Println("\nRun 'agen init' to set up agent templates")
		return nil
	}

	green := color.New(color.FgGreen)
	green.Printf("✓ IDE: %s\n\n", ideAdapter.Name())

	// Step 2: Get installed info
	installed, err := ide.GetInstalledInfo(absPath, ideAdapter)
	if err != nil && verbose {
		printWarning("Could not read installed info: %v", err)
	}

	// Step 3: Analyze project type
	projectType := analyzeProjectType(absPath)
	fmt.Printf("📁 Project Type: %s\n\n", color.YellowString(projectType))

	// Step 4: Calculate health metrics
	fmt.Println("📈 Health Metrics:")

	// Version status
	if installed != nil {
		latestVersion := templates.GetLatestVersion()
		if installed.Version == latestVersion {
			green.Println("  ✓ Templates: Up to date")
		} else {
			printWarning("  Templates: Outdated (have %s, latest %s)", installed.Version, latestVersion)
		}

		// Modified files
		if installed.ModifiedFiles > 0 {
			printInfo("  ℹ %d file(s) customized locally", installed.ModifiedFiles)
		} else {
			green.Println("  ✓ No local modifications")
		}
	}

	// Step 5: Agent recommendations based on project type
	fmt.Println("\n🎯 Agent Recommendations:")
	recommendations := getRecommendedAgents(projectType)

	for _, rec := range recommendations {
		// Check if agent is installed
		isInstalled := installed != nil && hasAgent(installed, rec.Name)

		if isInstalled {
			green.Printf("  ✓ %s\n", rec.Name)
		} else {
			if rec.Critical {
				printWarning("  ✗ %s (RECOMMENDED)", rec.Name)
				fmt.Printf("      %s\n", color.New(color.Faint).Sprint(rec.Reason))
			} else {
				printInfo("  ○ %s (optional)", rec.Name)
			}
		}
	}

	// Step 6: Calculate overall score
	score := calculateHealthScore(installed, recommendations)
	fmt.Printf("\n🏆 Health Score: ")
	if score >= 80 {
		green.Printf("%d/100", score)
		fmt.Println(" - Excellent!")
	} else if score >= 60 {
		color.Yellow("%d/100", score)
		fmt.Println(" - Good, some improvements possible")
	} else {
		color.Red("%d/100", score)
		fmt.Println(" - Needs attention")
	}

	// Step 7: Suggestions
	if score < 100 {
		fmt.Println("\n💡 Suggestions:")
		if installed == nil || latestOutdated(installed) {
			fmt.Println("  • Run 'agen update' to get latest templates")
		}
		for _, rec := range recommendations {
			if rec.Critical && (installed == nil || !hasAgent(installed, rec.Name)) {
				fmt.Printf("  • Add %s: agen init --agents %s\n", rec.Name, rec.Name)
			}
		}
	}

	fmt.Println()
	return nil
}

// AgentRecommendation represents a recommended agent for a project type
type AgentRecommendation struct {
	Name     string
	Critical bool
	Reason   string
}

// analyzeProjectType tries to figure out what kind of project this is.
// looks for telltale files like package.json, go.mod, requirements.txt etc.
func analyzeProjectType(path string) string {
	// Check for various project indicators
	indicators := map[string]string{
		"package.json":       "Node.js",
		"go.mod":             "Go",
		"requirements.txt":   "Python",
		"Cargo.toml":         "Rust",
		"pom.xml":            "Java (Maven)",
		"build.gradle":       "Java (Gradle)",
		"pubspec.yaml":       "Flutter",
		"Gemfile":            "Ruby",
		"composer.json":      "PHP",
		"next.config.js":     "Next.js",
		"next.config.mjs":    "Next.js",
		"next.config.ts":     "Next.js",
		"vite.config.js":     "Vite",
		"vite.config.ts":     "Vite",
		"nuxt.config.js":     "Nuxt.js",
		"angular.json":       "Angular",
		"app.json":           "React Native",
		"expo.json":          "Expo",
		"docker-compose.yml": "Docker",
	}

	for file, projectType := range indicators {
		if _, err := os.Stat(filepath.Join(path, file)); err == nil {
			return projectType
		}
	}

	return "Unknown"
}

// getRecommendedAgents returns recommended agents based on project type
func getRecommendedAgents(projectType string) []AgentRecommendation {
	// common recommendations for all projects
	common := []AgentRecommendation{
		{Name: "security-auditor", Critical: true, Reason: "Security is essential for any project"},
		{Name: "debugger", Critical: false, Reason: "Helpful for troubleshooting issues"},
	}

	switch projectType {
	case "Next.js", "Vite", "Nuxt.js", "Angular":
		return append(common,
			AgentRecommendation{Name: "frontend-specialist", Critical: true, Reason: "Essential for web frontends"},
			AgentRecommendation{Name: "test-engineer", Critical: true, Reason: "Testing is crucial for UIs"},
			AgentRecommendation{Name: "performance-optimizer", Critical: false, Reason: "Web performance matters"},
			AgentRecommendation{Name: "seo-specialist", Critical: false, Reason: "If SEO matters for your site"},
		)

	case "Node.js", "Python", "Go", "Rust", "Java (Maven)", "Java (Gradle)", "Ruby", "PHP":
		return append(common,
			AgentRecommendation{Name: "backend-specialist", Critical: true, Reason: "Essential for backend development"},
			AgentRecommendation{Name: "database-architect", Critical: true, Reason: "Most backends need database expertise"},
			AgentRecommendation{Name: "test-engineer", Critical: true, Reason: "Backend testing is crucial"},
		)

	case "React Native", "Expo", "Flutter":
		return append(common,
			AgentRecommendation{Name: "mobile-developer", Critical: true, Reason: "Essential for mobile apps"},
			AgentRecommendation{Name: "test-engineer", Critical: false, Reason: "Mobile testing helps"},
		)

	default:
		return common
	}
}

func hasAgent(info *ide.InstalledInfo, name string) bool {
	for _, a := range info.Agents {
		if a == name {
			return true
		}
	}
	return false
}

func latestOutdated(installed *ide.InstalledInfo) bool {
	if installed == nil {
		return true
	}
	return installed.Version != templates.GetLatestVersion()
}

func calculateHealthScore(installed *ide.InstalledInfo, recommendations []AgentRecommendation) int {
	if installed == nil {
		return 0
	}

	score := 50 // Base score for having templates installed

	// +20 for being up to date
	if !latestOutdated(installed) {
		score += 20
	}

	// +30 for having critical agents
	criticalCount := 0
	criticalInstalled := 0
	for _, rec := range recommendations {
		if rec.Critical {
			criticalCount++
			if hasAgent(installed, rec.Name) {
				criticalInstalled++
			}
		}
	}

	if criticalCount > 0 {
		score += (30 * criticalInstalled) / criticalCount
	}

	if score > 100 {
		score = 100
	}

	return score
}
