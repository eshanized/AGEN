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

	"github.com/eshanized/agen/internal/verify"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// verifyCmd runs verification scripts on the project.
// these are Go-based reimplementations of the original Python scripts.
var verifyCmd = &cobra.Command{
	Use:   "verify [path]",
	Short: "Run verification scripts",
	Long: `Run verification scripts to check project health.

Available checks:
  --security   Security scan (secrets, vulnerabilities)
  --lint       Lint and type checking
  --ux         UX audit (accessibility, usability)
  --seo        SEO check (meta tags, structure)
  --all        Run all checks (default)

Examples:
  agen verify                # Run all checks
  agen verify --security     # Only security scan
  agen verify --lint --ux    # Run multiple specific checks`,
	Args: cobra.MaximumNArgs(1),
	RunE: runVerify,
}

func init() {
	verifyCmd.Flags().Bool("security", false, "run security scan")
	verifyCmd.Flags().Bool("lint", false, "run lint check")
	verifyCmd.Flags().Bool("ux", false, "run UX audit")
	verifyCmd.Flags().Bool("seo", false, "run SEO check")
	verifyCmd.Flags().Bool("all", false, "run all checks")
	verifyCmd.Flags().Bool("fix", false, "attempt to auto-fix issues where possible")
	verifyCmd.Flags().StringP("output", "o", "text", "output format (text, json, markdown)")
}

// runVerify is the main logic for the verify command.
//
// How it works:
// 1. Parse which checks to run (default is all)
// 2. Initialize the verification runner
// 3. Run each check in priority order
// 4. Collect results and display summary
// 5. Exit with non-zero code if critical issues found
//
// Why Go-based? The original antigravity-kit used Python scripts.
// We rewrote them in Go so the tool is fully self-contained with no
// external dependencies. This makes installation way simpler.
func runVerify(cmd *cobra.Command, args []string) error {
	verbose := checkVerbose(cmd)

	// determine target directory
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// check directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", absPath)
	}

	// figure out which checks to run
	runSecurity, _ := cmd.Flags().GetBool("security")
	runLint, _ := cmd.Flags().GetBool("lint")
	runUX, _ := cmd.Flags().GetBool("ux")
	runSEO, _ := cmd.Flags().GetBool("seo")
	runAll, _ := cmd.Flags().GetBool("all")

	// if no specific checks requested, run all
	if !runSecurity && !runLint && !runUX && !runSEO {
		runAll = true
	}

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n🔍 AGEN Verification")
	fmt.Printf("Directory: %s\n\n", absPath)

	// create the runner
	runner := verify.NewRunner(absPath, verify.RunnerOptions{
		Verbose: verbose,
	})

	var results []verify.Result

	// Run checks in priority order (P0 → P5)
	if runAll || runSecurity {
		printInfo("Running security scan (P0)...")
		result := runner.RunSecurity()
		results = append(results, result)
		printCheckResult(result)
	}

	if runAll || runLint {
		printInfo("Running lint check (P1)...")
		result := runner.RunLint()
		results = append(results, result)
		printCheckResult(result)
	}

	if runAll || runUX {
		printInfo("Running UX audit (P4)...")
		result := runner.RunUX()
		results = append(results, result)
		printCheckResult(result)
	}

	if runAll || runSEO {
		printInfo("Running SEO check (P5)...")
		result := runner.RunSEO()
		results = append(results, result)
		printCheckResult(result)
	}

	// Print summary
	printVerifySummary(results)

	// exit with error if any critical issues
	for _, r := range results {
		if r.HasCritical {
			return fmt.Errorf("verification failed with critical issues")
		}
	}

	return nil
}

// printCheckResult shows the result of a single check
func printCheckResult(result verify.Result) {
	if result.Passed {
		printSuccess("%s: PASSED (%d issues)", result.Name, len(result.Issues))
	} else if result.HasCritical {
		printError("%s: FAILED (%d critical, %d warnings)",
			result.Name, result.CriticalCount, result.WarningCount)
	} else {
		printWarning("%s: WARNINGS (%d issues)", result.Name, len(result.Issues))
	}
}

// printVerifySummary shows the overall verification summary
func printVerifySummary(results []verify.Result) {
	fmt.Println()
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("📊 Summary")

	passed := 0
	failed := 0
	warnings := 0

	for _, r := range results {
		if r.Passed {
			passed++
		} else if r.HasCritical {
			failed++
		} else {
			warnings++
		}
	}

	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)

	green.Printf("  ✓ Passed:   %d\n", passed)
	yellow.Printf("  ⚠ Warnings: %d\n", warnings)
	red.Printf("  ✗ Failed:   %d\n", failed)
	fmt.Println()

	if failed > 0 {
		red.Println("❌ Verification FAILED - fix critical issues before proceeding")
	} else if warnings > 0 {
		yellow.Println("⚠️ Verification completed with warnings")
	} else {
		green.Println("✨ All checks passed!")
	}
	fmt.Println()
}
