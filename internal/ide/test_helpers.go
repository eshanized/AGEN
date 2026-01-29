// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Common test helpers for IDE adapter tests

package ide

import (
	"github.com/eshanized/agen/internal/templates"
)

// createMockTemplates creates a mock Templates struct for testing
func createMockTemplates() *templates.Templates {
	return &templates.Templates{
		Agents: map[string]templates.Agent{
			"test-agent": {
				Name:        "test-agent",
				Description: "A test agent for unit testing",
				Content:     "# Test Agent\n\nThis is a test agent.",
			},
			"another-agent": {
				Name:        "another-agent",
				Description: "Another test agent",
				Content:     "# Another Agent\n\nAnother test agent.",
			},
		},
		Skills: map[string]templates.Skill{
			"test-skill": {
				Name:        "test-skill",
				Description: "A test skill for unit testing",
				Content:     "# Test Skill\n\nThis is a test skill.",
			},
			"api-patterns": {
				Name:        "api-patterns",
				Description: "API design patterns",
				Content:     "# API Patterns\n\nREST, GraphQL, etc.",
			},
		},
		Workflows: map[string]templates.Workflow{
			"test-workflow": {
				Name:        "test-workflow",
				Description: "A test workflow for unit testing",
				Content:     "# Test Workflow\n\nThis is a test workflow.",
			},
			"deploy": {
				Name:        "deploy",
				Description: "Deploy to production",
				Content:     "# Deploy\n\nDeployment workflow.",
			},
		},
	}
}
