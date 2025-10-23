package main

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"
)

const todosFile = "todos.json"

type Todo struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type Mode int

const (
	Normal Mode = iota
	Insert
	Edit
)

type model struct {
	textarea textarea.Model
	todos    []Todo
	Mode
}

func loadTodos() []Todo {
	data, err := os.ReadFile(todosFile)
	if err != nil {
		return []Todo{}
	}
	var todos []Todo
	_ = json.Unmarshal(data, &todos)
	return todos
}

func saveTodos(todos []Todo) {
	data, _ := json.MarshalIndent(todos, "", "  ")
	_ = os.WriteFile(todosFile, data, 0644)
}

func todosToText(todos []Todo) string {
	lines := []string{}
	for _, t := range todos {
		check := "[ ]"
		if t.Done {
			check = "[x]"
		}
		lines = append(lines, fmt.Sprintf("%s %s", check, t.Text))
	}
	return strings.Join(lines, "\n")
}

func textToTodos(text string) []Todo {
	lines := strings.Split(text, "\n")
	todos := []Todo{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		done := false
		if strings.HasPrefix(line, "[x]") || strings.HasPrefix(line, "[X]") {
			done = true
			line = strings.TrimSpace(line[3:])
		} else if strings.HasPrefix(line, "[ ]") {
			line = strings.TrimSpace(line[3:])
		}
		todos = append(todos, Todo{Text: line, Done: done})
	}
	return todos
}

func initialModel() model {
	todos := loadTodos()
	ta := textarea.New()
	ta.SetValue(todosToText(todos))
	ta.Focus()
	return model{textarea: ta, todos: todos, Mode: Normal}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "a":
			if m.Mode == Normal {
				m.Mode = Insert
				m.todos = append(m.todos, Todo{Text: "", Done: false})
				m.textarea.SetValue(todosToText(m.todos))
				m.textarea.CursorEnd()
				m.textarea.Focus()
				return m, nil
			}

		case "e":
			if m.Mode == Normal && len(m.todos) > 0 {
				m.Mode = Edit
				m.textarea.Focus()
				return m, nil
			}

		case "enter":
			if m.Mode == Insert || m.Mode == Edit {
				todos := textToTodos(m.textarea.Value())

				// Remove empty tasks
				nonEmptyTodos := []Todo{}
				for _, t := range todos {
					if strings.TrimSpace(t.Text) != "" {
						nonEmptyTodos = append(nonEmptyTodos, t)
					}
				}

				m.todos = nonEmptyTodos
				saveTodos(m.todos)
				m.textarea.SetValue(todosToText(m.todos))
				m.Mode = Normal
			}

		// Navigation in Normal mode
		case "k", "up":
			if m.Mode == Normal {
				m.textarea.CursorUp()
			}
		case "j", "down":
			if m.Mode == Normal {
				m.textarea.CursorDown()
			}
		}

	case tea.WindowSizeMsg:
		m.textarea.SetWidth(msg.Width)
		m.textarea.SetHeight(15)
	}

	// Only let textarea process key events in Insert/Edit mode
	if m.Mode == Insert || m.Mode == Edit {
		m.textarea, cmd = m.textarea.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	modeText := ""
	switch m.Mode {
	case Normal:
		modeText = "Normal"
	case Insert:
		modeText = "Insert"
	case Edit:
		modeText = "Edit"
	}

	header := lipgloss.NewStyle().Bold(true).Render("📝 TODOs")
	modeBar := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(fmt.Sprintf("Mode: %s", modeText))
	footer := lipgloss.NewStyle().Faint(true).Render("a: add | e: edit | enter: save | q: quit")

	layout := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		modeBar,
		m.textarea.View(),
		footer,
	)

	return layout
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	p.Run()
}
