package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Task represents a single TODO item
type Task struct {
	Name     string `json:"name"`
	Priority string `json:"priority"`
	Done     bool   `json:"done"`
}

// Model represents the application state for Bubble Tea
type Model struct {
	Tasks         []Task
	Cursor        int
	Filter        string
	Command       string
	CommandMsg    string
	ActiveCmd     bool
	ShowCursor    bool
	CommandCursor int
}

// CursorBlinkMsg is sent periodically to blink the cursor
type CursorBlinkMsg struct{}

// BlinkCursor returns a command that sends cursor blink messages
func BlinkCursor() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(time.Time) tea.Msg {
		return CursorBlinkMsg{}
	})
}
