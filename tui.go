package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"time"
)

type Mode int

const (
	Normal Mode = iota
	Insert
	Edit
)

func blinkCursor() tea.Msg {
	return struct{}{}
}

func blinkCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(time.Time) tea.Msg {
		return blinkCursor()
	})
}

func (m model) Init() tea.Cmd {
	return blinkCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		}

		switch m.mode {
		case Normal:
			switch {
			case key.Matches(msg, keys.Down):
				if m.cursor < len(m.todos)-1 {
					m.cursor++
				}
			case key.Matches(msg, keys.Up):
				if m.cursor > 0 {
					m.cursor--
				}
			case key.Matches(msg, keys.Add):
				m.mode = Insert
				m.input = ""
				m.cursorVisible = true
				cmd = blinkCmd() // start blinking
			case key.Matches(msg, keys.Edit):
				if len(m.todos) > 0 {
					m.mode = Edit
					m.input = m.todos[m.cursor].Text
					m.cursorVisible = true
					cmd = blinkCmd() // start blinking
				}
			case key.Matches(msg, keys.Delete):
				if len(m.todos) > 0 {
					m.todos = append(m.todos[:m.cursor], m.todos[m.cursor+1:]...)
					if m.cursor >= len(m.todos) && m.cursor > 0 {
						m.cursor--
					}
					saveTodos(m.todos)
				}
			case key.Matches(msg, keys.Toggle):
				if len(m.todos) > 0 {
					m.todos[m.cursor].Done = !m.todos[m.cursor].Done
					saveTodos(m.todos)
				}
			}

		case Insert, Edit:
			switch {
			case key.Matches(msg, keys.Enter):
				if m.mode == Insert {
					m.todos = append(m.todos, Todo{Text: strings.TrimSpace(m.input)})
					m.cursor = len(m.todos) - 1
				} else if m.mode == Edit && len(m.todos) > 0 {
					m.todos[m.cursor].Text = strings.TrimSpace(m.input)
				}
				saveTodos(m.todos)
				m.mode = Normal
				m.input = ""
				m.cursorVisible = false
			case key.Matches(msg, keys.Esc):
				m.mode = Normal
				m.input = ""
				m.cursorVisible = false
			case msg.String() == "backspace":
				if len(m.input) > 0 {
					m.input = m.input[:len(m.input)-1]
				}
			default:
				if len(msg.String()) == 1 {
					m.input += msg.String()
				}
			}
		}

	case struct{}:
		if m.mode == Insert || m.mode == Edit {
			m.cursorVisible = !m.cursorVisible
			cmd = blinkCmd()
		}
	}

	if cmd != nil {
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	modeText := ""
	switch m.mode {
	case Normal:
		modeText = "Normal"
	case Insert:
		modeText = "Insert"
	case Edit:
		modeText = "Edit"
	}

	header := lipgloss.NewStyle().Bold(true).Render("📝 GoDoIt")
	body := m.todosToString()

	if m.mode == Insert {
		cursorChar := " "
		if m.cursorVisible {
			cursorChar = "_"
		}
		inputLine := fmt.Sprintf(" ➤[ ] %s%s", m.input, cursorChar)
		body += "\n" + lipgloss.NewStyle().Bold(true).Italic(true).Render(inputLine)
	}

	var help = m.help.View(keys)

	modeBar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Render(fmt.Sprintf("-- %s --", modeText))

	if m.mode == Insert || m.mode == Edit {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			"",
			body,
			"",
			modeBar,
			"",
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		body,
		"",
		"",
		"",
		help,
	)
}
