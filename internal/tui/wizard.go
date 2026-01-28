// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// A cross-platform CLI tool for managing AI agent templates

package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eshanized/agen/internal/templates"
)

// WizardResult contains the user's selections from the wizard
type WizardResult struct {
	IDE       string
	Agents    []string
	Skills    []string
	Cancelled bool
}

// wizardState tracks which step we're on
type wizardState int

const (
	stateIDESelect wizardState = iota
	stateAgentSelect
	stateSkillSelect
	stateConfirm
)

// Model is the bubbletea model for the interactive wizard
type Model struct {
	state     wizardState
	ideList   list.Model
	agentList list.Model
	skillList list.Model

	selectedIDE    string
	selectedAgents map[string]bool
	selectedSkills map[string]bool

	templates *templates.Templates
	width     int
	height    int
	quitting  bool
	confirmed bool
}

// item implements list.Item for our selection lists
type item struct {
	name        string
	description string
	selected    bool
}

func (i item) Title() string {
	if i.selected {
		return "✓ " + i.name
	}
	return "  " + i.name
}
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.name }

// Styles for the wizard
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginBottom(2)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(2)
)

// NewWizard creates a new interactive wizard model
func NewWizard(tmpl *templates.Templates) Model {
	// IDE options
	ideItems := []list.Item{
		item{name: "antigravity", description: "Claude Code / Antigravity (full .agent/ folder)"},
		item{name: "cursor", description: "Cursor IDE (.cursorrules file)"},
		item{name: "windsurf", description: "Windsurf IDE (.windsurfrules file)"},
		item{name: "zed", description: "Zed Editor (.zed/ folder with prompts)"},
	}

	ideDelegate := list.NewDefaultDelegate()
	ideList := list.New(ideItems, ideDelegate, 60, 10)
	ideList.Title = "Select IDE"
	ideList.SetShowStatusBar(false)
	ideList.SetFilteringEnabled(false)

	// Agent options
	var agentItems []list.Item
	for name, agent := range tmpl.Agents {
		agentItems = append(agentItems, item{
			name:        name,
			description: agent.Description,
		})
	}
	agentDelegate := list.NewDefaultDelegate()
	agentList := list.New(agentItems, agentDelegate, 60, 15)
	agentList.Title = "Select Agents (space to toggle, enter to continue)"
	agentList.SetShowStatusBar(false)

	// Skill options
	var skillItems []list.Item
	for name, skill := range tmpl.Skills {
		skillItems = append(skillItems, item{
			name:        name,
			description: skill.Description,
		})
	}
	skillDelegate := list.NewDefaultDelegate()
	skillList := list.New(skillItems, skillDelegate, 60, 15)
	skillList.Title = "Select Skills (space to toggle, enter to continue)"
	skillList.SetShowStatusBar(false)

	return Model{
		state:          stateIDESelect,
		ideList:        ideList,
		agentList:      agentList,
		skillList:      skillList,
		selectedAgents: make(map[string]bool),
		selectedSkills: make(map[string]bool),
		templates:      tmpl,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			switch m.state {
			case stateIDESelect:
				if i, ok := m.ideList.SelectedItem().(item); ok {
					m.selectedIDE = i.name
					m.state = stateAgentSelect
				}

			case stateAgentSelect:
				m.state = stateSkillSelect

			case stateSkillSelect:
				m.state = stateConfirm

			case stateConfirm:
				m.confirmed = true
				return m, tea.Quit
			}
			return m, nil

		case " ":
			// Toggle selection
			switch m.state {
			case stateAgentSelect:
				if i, ok := m.agentList.SelectedItem().(item); ok {
					m.selectedAgents[i.name] = !m.selectedAgents[i.name]
					// Update list item
					idx := m.agentList.Index()
					items := m.agentList.Items()
					items[idx] = item{
						name:        i.name,
						description: i.description,
						selected:    m.selectedAgents[i.name],
					}
					m.agentList.SetItems(items)
				}

			case stateSkillSelect:
				if i, ok := m.skillList.SelectedItem().(item); ok {
					m.selectedSkills[i.name] = !m.selectedSkills[i.name]
					idx := m.skillList.Index()
					items := m.skillList.Items()
					items[idx] = item{
						name:        i.name,
						description: i.description,
						selected:    m.selectedSkills[i.name],
					}
					m.skillList.SetItems(items)
				}
			}
			return m, nil

		case "a":
			// Select all
			switch m.state {
			case stateAgentSelect:
				items := m.agentList.Items()
				for idx, itm := range items {
					if i, ok := itm.(item); ok {
						m.selectedAgents[i.name] = true
						items[idx] = item{name: i.name, description: i.description, selected: true}
					}
				}
				m.agentList.SetItems(items)

			case stateSkillSelect:
				items := m.skillList.Items()
				for idx, itm := range items {
					if i, ok := itm.(item); ok {
						m.selectedSkills[i.name] = true
						items[idx] = item{name: i.name, description: i.description, selected: true}
					}
				}
				m.skillList.SetItems(items)
			}
			return m, nil

		case "n":
			// Select none
			switch m.state {
			case stateAgentSelect:
				items := m.agentList.Items()
				for idx, itm := range items {
					if i, ok := itm.(item); ok {
						m.selectedAgents[i.name] = false
						items[idx] = item{name: i.name, description: i.description, selected: false}
					}
				}
				m.agentList.SetItems(items)

			case stateSkillSelect:
				items := m.skillList.Items()
				for idx, itm := range items {
					if i, ok := itm.(item); ok {
						m.selectedSkills[i.name] = false
						items[idx] = item{name: i.name, description: i.description, selected: false}
					}
				}
				m.skillList.SetItems(items)
			}
			return m, nil

		case "esc":
			// go back
			if m.state > stateIDESelect {
				m.state--
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ideList.SetSize(msg.Width-4, msg.Height-10)
		m.agentList.SetSize(msg.Width-4, msg.Height-10)
		m.skillList.SetSize(msg.Width-4, msg.Height-10)
	}

	// Update the current list
	var cmd tea.Cmd
	switch m.state {
	case stateIDESelect:
		m.ideList, cmd = m.ideList.Update(msg)
	case stateAgentSelect:
		m.agentList, cmd = m.agentList.Update(msg)
	case stateSkillSelect:
		m.skillList, cmd = m.skillList.Update(msg)
	}

	return m, cmd
}

// View implements tea.Model
func (m Model) View() string {
	if m.quitting {
		return "Cancelled.\n"
	}

	var b strings.Builder

	b.WriteString(titleStyle.Render("🤖 AGEN Interactive Setup Wizard"))
	b.WriteString("\n")

	switch m.state {
	case stateIDESelect:
		b.WriteString(subtitleStyle.Render("Step 1/4: Choose your IDE"))
		b.WriteString("\n")
		b.WriteString(m.ideList.View())

	case stateAgentSelect:
		b.WriteString(subtitleStyle.Render(fmt.Sprintf("Step 2/4: Select Agents (IDE: %s)", m.selectedIDE)))
		b.WriteString("\n")
		b.WriteString(m.agentList.View())
		b.WriteString(helpStyle.Render("space: toggle • a: all • n: none • enter: continue • esc: back"))

	case stateSkillSelect:
		b.WriteString(subtitleStyle.Render("Step 3/4: Select Skills"))
		b.WriteString("\n")
		b.WriteString(m.skillList.View())
		b.WriteString(helpStyle.Render("space: toggle • a: all • n: none • enter: continue • esc: back"))

	case stateConfirm:
		b.WriteString(subtitleStyle.Render("Step 4/4: Confirm Selection"))
		b.WriteString("\n\n")

		b.WriteString(selectedStyle.Render("IDE: "))
		b.WriteString(m.selectedIDE)
		b.WriteString("\n\n")

		b.WriteString(selectedStyle.Render("Agents: "))
		var agents []string
		for name, selected := range m.selectedAgents {
			if selected {
				agents = append(agents, name)
			}
		}
		if len(agents) == 0 {
			b.WriteString("(all)")
		} else {
			b.WriteString(strings.Join(agents, ", "))
		}
		b.WriteString("\n\n")

		b.WriteString(selectedStyle.Render("Skills: "))
		var skills []string
		for name, selected := range m.selectedSkills {
			if selected {
				skills = append(skills, name)
			}
		}
		if len(skills) == 0 {
			b.WriteString("(all)")
		} else {
			b.WriteString(strings.Join(skills, ", "))
		}
		b.WriteString("\n\n")

		b.WriteString(helpStyle.Render("enter: install • esc: back • q: cancel"))
	}

	return b.String()
}

// GetResult returns the wizard result after it completes
func (m Model) GetResult() WizardResult {
	if m.quitting || !m.confirmed {
		return WizardResult{Cancelled: true}
	}

	var agents, skills []string
	for name, selected := range m.selectedAgents {
		if selected {
			agents = append(agents, name)
		}
	}
	for name, selected := range m.selectedSkills {
		if selected {
			skills = append(skills, name)
		}
	}

	return WizardResult{
		IDE:       m.selectedIDE,
		Agents:    agents,
		Skills:    skills,
		Cancelled: false,
	}
}

// RunWizard runs the interactive wizard and returns the result
func RunWizard(tmpl *templates.Templates) (WizardResult, error) {
	m := NewWizard(tmpl)
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return WizardResult{Cancelled: true}, err
	}

	return finalModel.(Model).GetResult(), nil
}

// Keymap bindings that are used by the wizard
type keyMap struct {
	Toggle key.Binding
	All    key.Binding
	None   key.Binding
	Back   key.Binding
	Quit   key.Binding
}

// Unused but kept for reference
var _ = textinput.Model{}
