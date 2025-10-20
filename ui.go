package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// InitialModel creates the initial application model
func InitialModel() Model {
	tasks, err := LoadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		tasks = []Task{}
	}
	return Model{
		Tasks:  tasks,
		Filter: "",
	}
}

// Init initializes the Bubble Tea model
func (m Model) Init() tea.Cmd {
	return BlinkCursor()
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case CursorBlinkMsg:
		if m.ActiveCmd {
			m.ShowCursor = !m.ShowCursor
		} else {
			m.ShowCursor = false
		}
		return m, BlinkCursor()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			// Remove done tasks before quitting
			m.Tasks = RemoveDoneTasks(m.Tasks)
			if err := SaveTasks(m.Tasks); err != nil {
				m.CommandMsg = fmt.Sprintf("Error saving tasks: %v", err)
			}
			fmt.Print("\033[H\033[2J")
			DisplayTasks(m.Tasks, "")
			return m, tea.Quit
		}

		if !m.ActiveCmd {
			// Navigation / single-key actions
			switch msg.String() {
			case "j":
				if m.Cursor < len(m.Tasks)-1 {
					m.Cursor++
				}
				return m, nil
			case "k":
				if m.Cursor > 0 {
					m.Cursor--
				}
				return m, nil
			case "x":
				if len(m.Tasks) > 0 && m.Cursor < len(m.Tasks) {
					m.Tasks[m.Cursor].Done = !m.Tasks[m.Cursor].Done
				}
				return m, nil
			case "d":
				if len(m.Tasks) > 0 {
					m.Tasks = append(m.Tasks[:m.Cursor], m.Tasks[m.Cursor+1:]...)
					if m.Cursor >= len(m.Tasks) && m.Cursor > 0 {
						m.Cursor--
					}
				}
				return m, nil
			case "a":
				m.Command = "add "
				m.ActiveCmd = true
				return m, nil
			case "p":
				m.Command = fmt.Sprintf("priority %d ", m.Cursor+1)
				m.ActiveCmd = true
				return m, nil
			case "f":
				m.Command = "filter "
				m.ActiveCmd = true
				return m, nil
			}
		} else {
			// Command typing
			if msg.Type == tea.KeyRunes || msg.String() == " " {
				m.Command += msg.String()

			} else if msg.Type == tea.KeyBackspace {

				if len(m.Command) > 0 {
					m.Command = m.Command[:len(m.Command)-1]
				} else {
					m.ActiveCmd = false
				}
			} else if msg.Type == tea.KeyEnter {
				if m.Command != "" {
					m.Tasks, m.Filter, m.CommandMsg = ExecuteCommand(m.Tasks, m.Command, m.Filter)
					m.Command = ""
				}
				m.ActiveCmd = false
			} else if msg.String() == "esc" {
				m.Command = ""
				m.ActiveCmd = false
			}
		}
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	var b strings.Builder

	b.WriteString("📝  TODO List (Interactive Mode)\n")
	b.WriteString("Controls: ↑/↓ move | x toggle done | d delete | a add | p priority | f filter | q quit\n\n")

	tasksToShow := m.Tasks
	if m.Filter != "" {
		var filtered []Task
		for _, t := range m.Tasks {
			if t.Priority == m.Filter {
				filtered = append(filtered, t)
			}
		}
		tasksToShow = filtered
	}

	for i, t := range tasksToShow {
		cursor := " "
		if i == m.Cursor {
			cursor = ">"
		}
		status := "[ ]"
		if t.Done {
			status = "[x]"
		}
		b.WriteString(fmt.Sprintf("%s %s %s (%s)\n", cursor, status, t.Name, t.Priority))
	}

	// Show command line only when actively typing
	if m.ActiveCmd {
		cursor := " "
		if m.ShowCursor {
			cursor = "|" // blinking cursor
		}
		b.WriteString("\n> " + m.Command + cursor)
	}

	if m.CommandMsg != "" {
		b.WriteString("\n" + m.CommandMsg)
	}

	return b.String()
}
