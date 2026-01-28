// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// AI-powered features: suggest, explain, compose

package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/eshanized/agen/internal/templates"
)

// Suggester analyzes projects and suggests appropriate agents
type Suggester struct {
	templates *templates.Templates
}

// Suggestion represents a recommended agent/skill
type Suggestion struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"` // agent, skill
	Score       float64 `json:"score"`
	Reason      string  `json:"reason"`
	Description string  `json:"description"`
}

// ProjectAnalysis contains analyzed project information
type ProjectAnalysis struct {
	ProjectType  string          `json:"project_type"`
	Languages    []string        `json:"languages"`
	Frameworks   []string        `json:"frameworks"`
	HasTests     bool            `json:"has_tests"`
	HasDocker    bool            `json:"has_docker"`
	HasCI        bool            `json:"has_ci"`
	Dependencies map[string]bool `json:"dependencies"`
	FileCount    int             `json:"file_count"`
}

// NewSuggester creates a new suggester
func NewSuggester() (*Suggester, error) {
	tmpl, err := templates.LoadEmbedded()
	if err != nil {
		return nil, err
	}

	return &Suggester{templates: tmpl}, nil
}

// Suggest analyzes a project and suggests agents
func (s *Suggester) Suggest(projectDir string) ([]Suggestion, error) {
	analysis := s.analyzeProject(projectDir)
	suggestions := s.generateSuggestions(analysis)

	// Sort by score descending
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Score > suggestions[j].Score
	})

	return suggestions, nil
}

// analyzeProject scans the project to understand its nature
func (s *Suggester) analyzeProject(dir string) *ProjectAnalysis {
	analysis := &ProjectAnalysis{
		Languages:    []string{},
		Frameworks:   []string{},
		Dependencies: make(map[string]bool),
	}

	// Detect languages by file extensions
	extCount := make(map[string]int)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// Skip hidden directories
		if strings.Contains(path, "/node_modules/") ||
			strings.Contains(path, "/.git/") ||
			strings.Contains(path, "/vendor/") {
			return filepath.SkipDir
		}

		ext := filepath.Ext(info.Name())
		extCount[ext]++
		analysis.FileCount++

		// Check for specific files
		name := info.Name()
		switch name {
		case "package.json":
			analysis.Dependencies["nodejs"] = true
		case "go.mod":
			analysis.Dependencies["go"] = true
		case "requirements.txt", "pyproject.toml":
			analysis.Dependencies["python"] = true
		case "Cargo.toml":
			analysis.Dependencies["rust"] = true
		case "Dockerfile", "docker-compose.yml":
			analysis.HasDocker = true
		case ".github", ".gitlab-ci.yml", "Jenkinsfile":
			analysis.HasCI = true
		}

		return nil
	})

	// Determine primary language
	langMap := map[string]string{
		".js": "javascript", ".ts": "typescript",
		".py": "python", ".go": "go",
		".rs": "rust", ".java": "java",
		".rb": "ruby", ".php": "php",
		".swift": "swift", ".kt": "kotlin",
	}

	for ext, lang := range langMap {
		if count, ok := extCount[ext]; ok && count > 5 {
			analysis.Languages = append(analysis.Languages, lang)
		}
	}

	// Detect frameworks
	if _, err := os.Stat(filepath.Join(dir, "next.config.js")); err == nil {
		analysis.Frameworks = append(analysis.Frameworks, "nextjs")
		analysis.ProjectType = "web"
	}
	if _, err := os.Stat(filepath.Join(dir, "vite.config.ts")); err == nil {
		analysis.Frameworks = append(analysis.Frameworks, "vite")
		analysis.ProjectType = "web"
	}
	if _, err := os.Stat(filepath.Join(dir, "app.json")); err == nil {
		analysis.Frameworks = append(analysis.Frameworks, "react-native")
		analysis.ProjectType = "mobile"
	}

	// Check for tests
	testDirs := []string{"test", "tests", "__tests__", "spec"}
	for _, td := range testDirs {
		if _, err := os.Stat(filepath.Join(dir, td)); err == nil {
			analysis.HasTests = true
			break
		}
	}

	return analysis
}

// generateSuggestions creates recommendations based on analysis
func (s *Suggester) generateSuggestions(analysis *ProjectAnalysis) []Suggestion {
	var suggestions []Suggestion

	// Always suggest core agents
	suggestions = append(suggestions, Suggestion{
		Name:        "security-auditor",
		Type:        "agent",
		Score:       0.9,
		Reason:      "Security is always important",
		Description: "Elite cybersecurity expert",
	})

	// Language-specific suggestions
	for _, lang := range analysis.Languages {
		switch lang {
		case "javascript", "typescript":
			suggestions = append(suggestions, Suggestion{
				Name:        "frontend-specialist",
				Type:        "agent",
				Score:       0.85,
				Reason:      fmt.Sprintf("Detected %s files", lang),
				Description: "Senior Frontend Architect",
			})
			suggestions = append(suggestions, Suggestion{
				Name:        "nextjs-react-expert",
				Type:        "skill",
				Score:       0.8,
				Reason:      "For React/Next.js optimization",
				Description: "Performance optimization expertise",
			})

		case "python":
			suggestions = append(suggestions, Suggestion{
				Name:        "backend-specialist",
				Type:        "agent",
				Score:       0.85,
				Reason:      "Detected Python project",
				Description: "Backend Development Architect",
			})
			suggestions = append(suggestions, Suggestion{
				Name:        "python-patterns",
				Type:        "skill",
				Score:       0.8,
				Reason:      "Python best practices",
				Description: "Python development patterns",
			})

		case "go":
			suggestions = append(suggestions, Suggestion{
				Name:        "backend-specialist",
				Type:        "agent",
				Score:       0.85,
				Reason:      "Detected Go project",
				Description: "Backend Development Architect",
			})
		}
	}

	// Framework-specific
	for _, fw := range analysis.Frameworks {
		switch fw {
		case "nextjs":
			// Already added via language detection
		case "react-native":
			suggestions = append(suggestions, Suggestion{
				Name:        "mobile-developer",
				Type:        "agent",
				Score:       0.9,
				Reason:      "React Native project detected",
				Description: "Mobile development expert",
			})
			suggestions = append(suggestions, Suggestion{
				Name:        "mobile-design",
				Type:        "skill",
				Score:       0.85,
				Reason:      "Mobile UI/UX patterns",
				Description: "Mobile-first design thinking",
			})
		}
	}

	// Testing suggestions
	if !analysis.HasTests {
		suggestions = append(suggestions, Suggestion{
			Name:        "test-engineer",
			Type:        "agent",
			Score:       0.75,
			Reason:      "No test directory found",
			Description: "Testing and TDD expert",
		})
	}

	// DevOps suggestions
	if analysis.HasDocker {
		suggestions = append(suggestions, Suggestion{
			Name:        "devops-engineer",
			Type:        "agent",
			Score:       0.8,
			Reason:      "Docker detected",
			Description: "DevOps and infrastructure",
		})
	}

	// CI/CD
	if analysis.HasCI {
		suggestions = append(suggestions, Suggestion{
			Name:        "deployment-procedures",
			Type:        "skill",
			Score:       0.7,
			Reason:      "CI pipeline detected",
			Description: "Deployment best practices",
		})
	}

	// Always useful skills
	suggestions = append(suggestions, Suggestion{
		Name:        "clean-code",
		Type:        "skill",
		Score:       0.95,
		Reason:      "Universal best practice",
		Description: "Pragmatic coding standards",
	})

	suggestions = append(suggestions, Suggestion{
		Name:        "brainstorming",
		Type:        "skill",
		Score:       0.7,
		Reason:      "Helps with complex decisions",
		Description: "Socratic questioning protocol",
	})

	return suggestions
}

// ExplainAgent provides detailed explanation of an agent
func (s *Suggester) ExplainAgent(name string) (string, error) {
	agent, ok := s.templates.Agents[name]
	if !ok {
		return "", fmt.Errorf("agent not found: %s", name)
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", name))
	sb.WriteString(fmt.Sprintf("**Description:** %s\n\n", agent.Description))

	if len(agent.Skills) > 0 {
		sb.WriteString("**Skills:**\n")
		for _, skill := range agent.Skills {
			sb.WriteString(fmt.Sprintf("- %s\n", skill))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Full Content\n\n")
	sb.WriteString(agent.Content)

	return sb.String(), nil
}

// ExplainSkill provides detailed explanation of a skill
func (s *Suggester) ExplainSkill(name string) (string, error) {
	skill, ok := s.templates.Skills[name]
	if !ok {
		return "", fmt.Errorf("skill not found: %s", name)
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", name))
	sb.WriteString(fmt.Sprintf("**Description:** %s\n\n", skill.Description))
	sb.WriteString("## Full Content\n\n")
	sb.WriteString(skill.Content)

	return sb.String(), nil
}

// ComposeAgent creates a custom agent based on description
type ComposedAgent struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Skills      []string `json:"skills"`
	Content     string   `json:"content"`
}

// Compose generates a custom agent based on user description
func (s *Suggester) Compose(name, description string, baseAgents []string) (*ComposedAgent, error) {
	composed := &ComposedAgent{
		Name:        name,
		Description: description,
		Skills:      []string{},
	}

	var sb strings.Builder

	// Generate frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("name: %s\n", name))
	sb.WriteString(fmt.Sprintf("description: %s\n", description))

	// Collect skills from base agents
	skillSet := make(map[string]bool)
	for _, baseAgent := range baseAgents {
		if agent, ok := s.templates.Agents[baseAgent]; ok {
			for _, skill := range agent.Skills {
				skillSet[skill] = true
			}
		}
	}

	if len(skillSet) > 0 {
		sb.WriteString("skills: ")
		var skills []string
		for skill := range skillSet {
			skills = append(skills, skill)
		}
		sb.WriteString(strings.Join(skills, ", "))
		sb.WriteString("\n")
		composed.Skills = skills
	}

	sb.WriteString("---\n\n")

	// Generate content header
	sb.WriteString(fmt.Sprintf("# %s\n\n", name))
	sb.WriteString(fmt.Sprintf("> %s\n\n", description))

	// Add composed rules from base agents
	sb.WriteString("## Core Principles\n\n")

	for _, baseAgent := range baseAgents {
		if agent, ok := s.templates.Agents[baseAgent]; ok {
			sb.WriteString(fmt.Sprintf("### From %s\n\n", baseAgent))
			// Extract key points (simplified - just take first 500 chars)
			content := agent.Content
			if len(content) > 500 {
				content = content[:500] + "..."
			}
			sb.WriteString(content)
			sb.WriteString("\n\n")
		}
	}

	composed.Content = sb.String()
	return composed, nil
}
