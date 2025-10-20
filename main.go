package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Task struct {
	Name     string `json:"name"`
	Priority string `json:"priority"`
	Done     bool   `json:"done"`
}

// ---------- MODEL ----------
type model struct {
	tasks         []Task
	cursor        int
	filter        string
	command       string
	commandMsg    string
	activeCmd     bool
	showCursor    bool
	commandCursor int
}

const taskFile = "$HOME/.todo_tasks.json" // change to a preferred location

func initialModel() model {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		tasks = []Task{}
	}
	return model{
		tasks:  tasks,
		filter: "",
	}
}

// ---------- MAIN ----------
func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-i":
			runInteractive()
			return
		case "-l":
			displayTasks(initialModel().tasks, "")
			return
		default:
			fmt.Println("Usage: terminal-todo [-i for interactive|-l for list]")
			return
		}
	}

	// Default: interactive mode
	runInteractive()
}

// ---------- DISPLAY-ONLY MODE ----------
func displayTasks(tasks []Task, filter string) {
	fmt.Println("TODO List:")

	tasksToShow := tasks
	if filter != "" {
		var filtered []Task
		for _, t := range tasks {
			if t.Priority == filter {
				filtered = append(filtered, t)
			}
		}
		tasksToShow = filtered
	}

	for i, t := range tasksToShow {
		status := "[ ]"
		if t.Done {
			status = "[x]"
		}
		// Always start printing at column 0
		fmt.Printf("%d. %s %s (%s)\n", i+1, status, t.Name, t.Priority)
	}
}

// ---------- INTERACTIVE MODE ----------
func runInteractive() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

// ---------- TEA MODEL METHODS ----------
func (m model) Init() tea.Cmd {
	return blinkCursor()
}

func blinkCursor() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(time.Time) tea.Msg {
		return cursorBlinkMsg{}
	})
}

type cursorBlinkMsg struct{}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case cursorBlinkMsg:
		if m.activeCmd {
			m.showCursor = !m.showCursor
		} else {
			m.showCursor = false
		}
		return m, blinkCursor()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			// Remove done tasks before quitting
			m.tasks = removeDoneTasks(m.tasks)
			if err := saveTasks(m.tasks); err != nil {
				m.commandMsg = fmt.Sprintf("Error saving tasks: %v", err)
			}
			fmt.Print("\033[H\033[2J")
			displayTasks(m.tasks, "")
			return m, tea.Quit

		}

		if !m.activeCmd {
			// Navigation / single-key actions
			switch msg.String() {
			case "j":
				if m.cursor < len(m.tasks)-1 {
					m.cursor++
				}
				return m, nil
			case "k":
				if m.cursor > 0 {
					m.cursor--
				}
				return m, nil
			case "x":
				m.tasks[m.cursor].Done = !m.tasks[m.cursor].Done
				return m, nil
			case "d":
				if len(m.tasks) > 0 {
					m.tasks = append(m.tasks[:m.cursor], m.tasks[m.cursor+1:]...)
					if m.cursor >= len(m.tasks) && m.cursor > 0 {
						m.cursor--
					}
				}
				return m, nil
			case "a":
				m.command = "add "
				m.activeCmd = true
				return m, nil
			case "p":
				m.command = fmt.Sprintf("priority %d ", m.cursor+1)
				m.activeCmd = true
				return m, nil
			case "f":
				m.command = "filter "
				m.activeCmd = true
				return m, nil
			}
		} else {
			// Command typing
			if msg.Type == tea.KeyRunes || msg.String() == " " {
				m.command += msg.String()

			} else if msg.Type == tea.KeyBackspace {

				if len(m.command) > 0 {
					m.command = m.command[:len(m.command)-1]
				} else {
					m.activeCmd = false
				}
			} else if msg.Type == tea.KeyEnter {
				if m.command != "" {
					m.tasks, m.filter, m.commandMsg = executeCommand(m.tasks, m.command, m.filter)
					m.command = ""
				}
				m.activeCmd = false
			} else if msg.String() == "esc" {
				m.command = ""
				m.activeCmd = false
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString("📝  TODO List (Interactive Mode)\n")
	b.WriteString("Controls: ↑/↓ move | x toggle done | d delete | a add | p priority | f filter | q quit\n\n")
	tasksToShow := m.tasks
	if m.filter != "" {
		var filtered []Task
		for _, t := range m.tasks {
			if t.Priority == m.filter {
				filtered = append(filtered, t)
			}
		}
		tasksToShow = filtered
	}

	for i, t := range tasksToShow {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		status := "[ ]"
		if t.Done {
			status = "[x]"
		}
		b.WriteString(fmt.Sprintf("%s %s %s (%s)\n", cursor, status, t.Name, t.Priority))
	}

	// Show command line only when actively typing
	if m.activeCmd {
		cursor := " "
		if m.showCursor {
			cursor = "|" // blinking cursor
		}
		b.WriteString("\n> " + m.command + cursor)
	}

	if m.commandMsg != "" {
		b.WriteString("\n" + m.commandMsg)
	}

	return b.String()
}

// ---------- COMMAND HANDLER ----------
func executeCommand(tasks []Task, cmd string, filter string) ([]Task, string, string) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return tasks, filter, ""
	}

	switch parts[0] {
	case "add":
		if len(parts) > 1 {
			newTask := Task{Name: strings.Join(parts[1:], " "), Priority: "low"}
			tasks = append(tasks, newTask)
			return tasks, filter, "Added task."
		}
	case "remove":
		if len(parts) > 1 {
			idx := parseIndex(parts[1])
			if idx >= 0 && idx < len(tasks) {
				tasks = append(tasks[:idx], tasks[idx+1:]...)
				return tasks, filter, "Removed task."
			}
		}
	case "priority":
		if len(parts) > 2 {
			idx := parseIndex(parts[1])
			if idx >= 0 && idx < len(tasks) {
				tasks[idx].Priority = parts[2]
				return tasks, filter, "Changed priority."
			}
		}
	case "filter":
		if len(parts) > 1 {
			return tasks, parts[1], "Filter applied."
		}
		return tasks, "", "Filter cleared."
	}
	return tasks, filter, "Unknown command."
}

// ---------- HELPER FUNCTIONS ----------
func parseIndex(s string) int {
	var idx int
	_, err := fmt.Sscanf(s, "%d", &idx)
	if err != nil {
		return -1
	}
	return idx - 1
}

func removeDoneTasks(tasks []Task) []Task {
	var remaining []Task
	for _, t := range tasks {
		if !t.Done {
			remaining = append(remaining, t)
		}
	}
	return remaining
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(taskFile, data, 0644)
}

func loadTasks() ([]Task, error) {
	data, err := os.ReadFile(taskFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}
