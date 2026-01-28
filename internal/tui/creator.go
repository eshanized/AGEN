// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Eshan Roy <eshanized@proton.me>
//
// AGEN - AI Agent Template Manager
// Agent Creator TUI

package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eshanized/agen/internal/templates"
)

// CreatorResult contains the data collected by the creator wizard
type CreatorResult struct {
	Name        string
	Description string
	Skills      []string
	Cancelled   bool
}

// creatorState tracks the current step
type creatorState int

const (
	stateNameInput creatorState = iota
	stateDescInput
	stateSkills
	stateCreatorConfirm
)

// CreatorModel is the bubbletea model for the agent creator
type CreatorModel struct {
	state creatorState

	nameInput textinput.Model
	descInput textinput.Model
	skillList list.Model

	selectedSkills map[string]bool
	templates      *templates.Templates

	width     int
	height    int
	quitting  bool
	confirmed bool
}

// NewCreator creates a new interactive creator model
func NewCreator(tmpl *templates.Templates) CreatorModel {
	// 1. Name Input
	ni := textinput.New()
	ni.Placeholder = "e.g. security-auditor"
	ni.Focus()
	ni.CharLimit = 50
	ni.Width = 40

	// 2. Description Input
	di := textinput.New()
	di.Placeholder = "What does this agent do?"
	di.CharLimit = 100
	di.Width = 60

	// 3. Skill List
	var skillItems []list.Item
	for name, skill := range tmpl.Skills {
		skillItems = append(skillItems, item{
			name:        name,
			description: skill.Description,
		})
	}
	skillDelegate := list.NewDefaultDelegate()
	skillList := list.New(skillItems, skillDelegate, 0, 0)
	skillList.Title = "Select Skills (space to toggle)"
	skillList.SetShowStatusBar(false)

	return CreatorModel{
		state:          stateNameInput,
		nameInput:      ni,
		descInput:      di,
		skillList:      skillList,
		selectedSkills: make(map[string]bool),
		templates:      tmpl,
	}
}

// Init implements tea.Model
func (m CreatorModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model
func (m CreatorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			switch m.state {
			case stateNameInput:
				if m.nameInput.Value() != "" {
					m.state = stateDescInput
					m.descInput.Focus()
					return m, textinput.Blink
				}
			case stateDescInput:
				if m.descInput.Value() != "" {
					m.state = stateSkills
					// blur inputs
					m.nameInput.Blur()
					m.descInput.Blur()
				}
			case stateSkills:
				m.state = stateCreatorConfirm
			case stateCreatorConfirm:
				m.confirmed = true
				return m, tea.Quit
			}

		case "space":
			if m.state == stateSkills {
				if i, ok := m.skillList.SelectedItem().(item); ok {
					m.selectedSkills[i.name] = !m.selectedSkills[i.name]
					// Update view status
					idx := m.skillList.Index()
					items := m.skillList.Items()
					items[idx] = item{
						name:        i.name,
						description: i.description,
						selected:    m.selectedSkills[i.name],
					}
					m.skillList.SetItems(items)
				}
				return m, nil
			}

		case "esc":
			if m.state > stateNameInput {
				m.state--
				// manage focus
				if m.state == stateNameInput {
					m.nameInput.Focus()
				} else if m.state == stateDescInput {
					m.descInput.Focus()
				}
				return m, nil
			} else {
				m.quitting = true
				return m, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.skillList.SetSize(msg.Width-4, msg.Height-10)
	}

	// Route events to components
	switch m.state {
	case stateNameInput:
		m.nameInput, cmd = m.nameInput.Update(msg)
	case stateDescInput:
		m.descInput, cmd = m.descInput.Update(msg)
	case stateSkills:
		m.skillList, cmd = m.skillList.Update(msg)
	}

	return m, cmd
}

// View implements tea.Model
func (m CreatorModel) View() string {
	if m.quitting {
		return "Cancelled.\n"
	}

	var b strings.Builder

	b.WriteString(titleStyle.Render("🧙 Agent Creator Wizard"))
	b.WriteString("\n\n")

	switch m.state {
	case stateNameInput:
		b.WriteString(subtitleStyle.Render("Step 1/4: Name your agent"))
		b.WriteString("\n")
		b.WriteString(m.nameInput.View())
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("enter: next • esc: cancel"))

	case stateDescInput:
		b.WriteString(subtitleStyle.Render("Step 2/4: Describe role & responsibilities"))
		b.WriteString("\n")
		b.WriteString(m.descInput.View())
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("enter: next • esc: back"))

	case stateSkills:
		b.WriteString(subtitleStyle.Render("Step 3/4: Add skills"))
		b.WriteString("\n")
		b.WriteString(m.skillList.View())
		b.WriteString(helpStyle.Render("space: toggle • enter: next • esc: back"))

	case stateCreatorConfirm:
		b.WriteString(subtitleStyle.Render("Step 4/4: Ready to create?"))
		b.WriteString("\n\n")

		b.WriteString(selectedStyle.Render("Name: "))
		b.WriteString(m.nameInput.Value())
		b.WriteString("\n")

		b.WriteString(selectedStyle.Render("Description: "))
		b.WriteString(m.descInput.Value())
		b.WriteString("\n\n")

		b.WriteString(selectedStyle.Render("Skills: "))
		var skills []string
		for name, selected := range m.selectedSkills {
			if selected {
				skills = append(skills, name)
			}
		}
		if len(skills) == 0 {
			b.WriteString("(none)")
		} else {
			b.WriteString(strings.Join(skills, ", "))
		}
		b.WriteString("\n\n")

		b.WriteString(helpStyle.Render("enter: create file • esc: back • ctrl+c: cancel"))
	}

	return b.String()
}

// GetResult returns the collected data
func (m CreatorModel) GetResult() CreatorResult {
	if m.quitting || !m.confirmed {
		return CreatorResult{Cancelled: true}
	}

	var skills []string
	for name, selected := range m.selectedSkills {
		if selected {
			skills = append(skills, name)
		}
	}

	return CreatorResult{
		Name:        m.nameInput.Value(),
		Description: m.descInput.Value(),
		Skills:      skills,
		Cancelled:   false,
	}
}

// RunCreator launches the wizard
func RunCreator(tmpl *templates.Templates) (CreatorResult, error) {
	m := NewCreator(tmpl)
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return CreatorResult{Cancelled: true}, err
	}

	return finalModel.(CreatorModel).GetResult(), nil
}
