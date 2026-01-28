// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// A cross-platform CLI tool for managing AI agent templates

package main

import (
	"os"

	"github.com/eshanized/agen/internal/cli"
)

// main is the entry point for the agen CLI tool.
//
// How it works:
// 1. We hand off control to the cli package which sets up all commands
// 2. If there's an error during execution, we exit with code 1
// 3. otherwise we exit cleanly with code 0
//
// Why separate the CLI logic? Keeps main.go minimal and testable.
// The cli package handles all the Cobra command setup and routing.
func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
