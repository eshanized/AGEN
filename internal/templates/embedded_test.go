// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Unit tests for template loading and parsing

package templates

import (
	"testing"
)

func TestLoadEmbedded(t *testing.T) {
	tmpl, err := LoadEmbedded()
	if err != nil {
		t.Fatalf("LoadEmbedded() failed: %v", err)
	}

	if tmpl == nil {
		t.Fatal("LoadEmbedded() returned nil")
	}

	// Check we have templates loaded
	if len(tmpl.Agents) == 0 {
		t.Error("No agents loaded")
	}

	if len(tmpl.Skills) == 0 {
		t.Error("No skills loaded")
	}

	if len(tmpl.Workflows) == 0 {
		t.Error("No workflows loaded")
	}

	t.Logf("Loaded %d agents, %d skills, %d workflows",
		len(tmpl.Agents), len(tmpl.Skills), len(tmpl.Workflows))
}

func TestParseFrontmatter(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantDesc string
		wantBody string
	}{
		{
			name: "with frontmatter",
			content: `---
description: Test description
name: test
---

# Content here`,
			wantDesc: "Test description",
			wantBody: "# Content here",
		},
		{
			name:     "without frontmatter",
			content:  "# Just content",
			wantDesc: "",
			wantBody: "# Just content",
		},
		{
			name: "malformed frontmatter",
			content: `---
this is not valid yaml
---
Content`,
			wantDesc: "",
			wantBody: `---
this is not valid yaml
---
Content`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, body := parseFrontmatter(tt.content)

			if fm != nil {
				if desc, ok := fm["description"].(string); ok {
					if desc != tt.wantDesc {
						t.Errorf("description = %q, want %q", desc, tt.wantDesc)
					}
				}
			}

			if body != tt.wantBody {
				t.Errorf("body = %q, want %q", body, tt.wantBody)
			}
		})
	}
}

func TestTemplatesFilter(t *testing.T) {
	// Create a mock templates struct
	tmpl := &Templates{
		Version: "1.0.0",
		Agents: map[string]Agent{
			"frontend": {Name: "frontend"},
			"backend":  {Name: "backend"},
			"security": {Name: "security"},
		},
		Skills: map[string]Skill{
			"clean-code": {Name: "clean-code"},
			"testing":    {Name: "testing"},
		},
		Workflows: map[string]Workflow{
			"create": {Name: "create"},
		},
	}

	// Test filtering agents
	filtered := tmpl.Filter([]string{"frontend", "backend"}, nil)

	if len(filtered.Agents) != 2 {
		t.Errorf("Expected 2 agents, got %d", len(filtered.Agents))
	}

	if _, ok := filtered.Agents["frontend"]; !ok {
		t.Error("frontend agent should be in filtered results")
	}

	if _, ok := filtered.Agents["security"]; ok {
		t.Error("security agent should NOT be in filtered results")
	}

	// Workflows should still be included
	if len(filtered.Workflows) != len(tmpl.Workflows) {
		t.Error("Workflows should not be filtered")
	}
}

func TestGetLatestVersion(t *testing.T) {
	version := GetLatestVersion()
	if version == "" {
		t.Error("GetLatestVersion() returned empty string")
	}

	// Should match CurrentVersion
	if version != CurrentVersion {
		t.Errorf("GetLatestVersion() = %q, want %q", version, CurrentVersion)
	}
}
