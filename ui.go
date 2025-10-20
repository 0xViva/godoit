package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// InitialModel creates the initial application model
func InitialModel() Model {
	tasks, err := LoadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		tasks = []Task{}
	}
	m := Model{
		Tasks:  tasks,
		Filter: "",
		Cursor: 0,
	}
	// Initialize cursor to first visible task
	visibleIndices := m.getVisibleTaskIndices()
	if len(visibleIndices) > 0 {
		m.Cursor = visibleIndices[0]
	}
	return m
}

// getVisibleTaskIndices returns the indices of all tasks that should be visible
func (m Model) getVisibleTaskIndices() []int {
	var indices []int
	for i, t := range m.Tasks {
		// Apply filter if set
		if m.Filter != "" && t.Priority != m.Filter {
			continue
		}
		// All tasks are visible (active, done, deleted)
		indices = append(indices, i)
	}
	return indices
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
			DisplayTasks(m.Tasks, "", "all")
			return m, tea.Quit
		}

		if !m.ActiveCmd {
			// Navigation / single-key actions
			switch msg.String() {
			case "j":
				// Get visible tasks to determine max cursor position
				visibleIndices := m.getVisibleTaskIndices()
				if len(visibleIndices) > 0 {
					// Find current position in visible list
					currentPos := -1
					for i, idx := range visibleIndices {
						if idx == m.Cursor {
							currentPos = i
							break
						}
					}
					// Move to next visible task (wrap around to top)
					if currentPos >= 0 && currentPos < len(visibleIndices)-1 {
						m.Cursor = visibleIndices[currentPos+1]
					} else if currentPos == len(visibleIndices)-1 {
						// Wrap to first visible task
						m.Cursor = visibleIndices[0]
					}
				}
				return m, nil
			case "k":
				// Get visible tasks to determine navigation
				visibleIndices := m.getVisibleTaskIndices()
				if len(visibleIndices) > 0 {
					// Find current position in visible list
					currentPos := -1
					for i, idx := range visibleIndices {
						if idx == m.Cursor {
							currentPos = i
							break
						}
					}
					// Move to previous visible task (wrap around to bottom)
					if currentPos > 0 {
						m.Cursor = visibleIndices[currentPos-1]
					} else if currentPos == 0 {
						// Wrap to last visible task
						m.Cursor = visibleIndices[len(visibleIndices)-1]
					}
				}
				return m, nil
			case "x":
				// Mark task as done (only works on active tasks)
				if len(m.Tasks) > 0 && m.Cursor < len(m.Tasks) {
					task := &m.Tasks[m.Cursor]
					if task.Status == "active" || task.Status == "" {
						task.Done = true
						task.Status = "done"
						now := time.Now()
						task.CompletedAt = &now

						// Find next active task
						nextActive := -1
						for i := m.Cursor + 1; i < len(m.Tasks); i++ {
							if m.Tasks[i].Status == "active" || m.Tasks[i].Status == "" {
								nextActive = i
								break
							}
						}
						if nextActive != -1 {
							m.Cursor = nextActive
						} else {
							// No active tasks after, look before
							for i := m.Cursor - 1; i >= 0; i-- {
								if m.Tasks[i].Status == "active" || m.Tasks[i].Status == "" {
									m.Cursor = i
									break
								}
							}
						}
					}
				}
				return m, nil
			case "d":
				// Delete task (only works on active or done tasks)
				if len(m.Tasks) > 0 && m.Cursor < len(m.Tasks) {
					task := &m.Tasks[m.Cursor]
					if task.Status == "active" || task.Status == "" || task.Status == "done" {
						wasActive := (task.Status == "active" || task.Status == "")
						task.Status = "deleted"
						task.Done = false
						now := time.Now()
						task.DeletedAt = &now

						// If we deleted an active task, find next active task
						if wasActive {
							// Look for next active task after current position
							nextActive := -1
							for i := m.Cursor + 1; i < len(m.Tasks); i++ {
								if m.Tasks[i].Status == "active" || m.Tasks[i].Status == "" {
									nextActive = i
									break
								}
							}
							if nextActive != -1 {
								m.Cursor = nextActive
							} else {
								// No active tasks after, look before
								for i := m.Cursor - 1; i >= 0; i-- {
									if m.Tasks[i].Status == "active" || m.Tasks[i].Status == "" {
										m.Cursor = i
										break
									}
								}
							}
						}
					}
				}
				return m, nil
			case "u":
				// Undo - restore done or deleted tasks back to active
				if len(m.Tasks) > 0 && m.Cursor < len(m.Tasks) {
					task := &m.Tasks[m.Cursor]
					if task.Status == "done" || task.Status == "deleted" {
						task.Status = "active"
						task.Done = false
						task.CompletedAt = nil
						task.DeletedAt = nil
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
			case "e":
				if len(m.Tasks) > 0 && m.Cursor < len(m.Tasks) {
					m.Command = fmt.Sprintf("edit %d %s", m.Cursor+1, m.Tasks[m.Cursor].Name)
					m.ActiveCmd = true
				}
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

	b.WriteString("📝  TODOs:\n")
	b.WriteString("Controls: ↑/↓ move | x done | d delete | u undo | a add | e edit | p priority | f filter | q quit\n\n")

	// Separate tasks by status for display
	var activeTasks, doneTasks, deletedTasks []Task
	var activeIndices, doneIndices, deletedIndices []int

	for i, t := range m.Tasks {
		if m.Filter != "" && t.Priority != m.Filter {
			continue
		}
		switch t.Status {
		case "deleted":
			deletedTasks = append(deletedTasks, t)
			deletedIndices = append(deletedIndices, i)
		case "done":
			doneTasks = append(doneTasks, t)
			doneIndices = append(doneIndices, i)
		default: // "active" or empty (for backward compatibility)
			activeTasks = append(activeTasks, t)
			activeIndices = append(activeIndices, i)
		}
	}

	// Display active tasks
	if len(activeTasks) > 0 {
		b.WriteString("=== ACTIVE ===\n")
		for i, t := range activeTasks {
			cursor := " "
			if activeIndices[i] == m.Cursor {
				cursor = ">"
			}
			age := ""
			if !t.CreatedAt.IsZero() {
				age = fmt.Sprintf(" [%s]", FormatTaskAge(t.CreatedAt))
			}
			// Show task ID (index+1) instead of display order
			b.WriteString(fmt.Sprintf("%s [ ] %d. %s (%s)%s\n", cursor, activeIndices[i]+1, t.Name, t.Priority, age))
		}
		b.WriteString("\n")
	}

	// Display done tasks
	if len(doneTasks) > 0 {
		b.WriteString("=== DONE ===\n")
		for i, t := range doneTasks {
			cursor := " "
			if doneIndices[i] == m.Cursor {
				cursor = ">"
			}
			age := ""
			if !t.CreatedAt.IsZero() {
				age = fmt.Sprintf(" [%s]", FormatTaskAge(t.CreatedAt))
			}
			// Show task ID (index+1) instead of display order
			b.WriteString(fmt.Sprintf("%s [✓] %d. %s (%s)%s\n", cursor, doneIndices[i]+1, t.Name, t.Priority, age))
		}
		b.WriteString("\n")
	}

	// Display deleted tasks
	if len(deletedTasks) > 0 {
		b.WriteString("=== DELETED ===\n")
		for i, t := range deletedTasks {
			cursor := " "
			if deletedIndices[i] == m.Cursor {
				cursor = ">"
			}
			age := ""
			if !t.CreatedAt.IsZero() {
				age = fmt.Sprintf(" [%s]", FormatTaskAge(t.CreatedAt))
			}
			// Show task ID (index+1) instead of display order
			b.WriteString(fmt.Sprintf("%s [DELETED] %d. %s (%s)%s\n", cursor, deletedIndices[i]+1, t.Name, t.Priority, age))
		}
		b.WriteString("\n")
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
