package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"time"
)

type Todo struct {
	ID          int        `json:"id"`
	Text        string     `json:"text"`
	Done        bool       `json:"done"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

func (m model) todosToString() string {
	var b strings.Builder

	maxLen := 0
	maxID := 0
	for _, t := range m.todos {
		l := len([]rune(t.Text))
		if l > maxLen {
			maxLen = l
		}
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	padding := 3
	width := maxLen + padding

	ageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Faint(true)

	idStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Faint(true)

	lineHighlight := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15"))

	cursorIDStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Bold(true)

	for i, todo := range m.todos {
		cursor := " "
		if m.cursor == i && (m.mode == Normal || m.mode == Edit) {
			cursor = "➤"
		}

		check := "[ ]"
		if todo.Done {
			check = "[x]"
		}

		text := todo.Text
		textStyle := lipgloss.NewStyle()

		if todo.Done {
			textStyle = textStyle.Strikethrough(true)
		}

		if m.cursor == i && (m.mode == Normal || m.mode == Edit) {
			textStyle = textStyle.Bold(true)
		}

		if m.mode == Edit && i == m.cursor {
			cursorChar := " "
			if m.cursorVisible {
				cursorChar = "_"
			}
			text = m.input + cursorChar
			textStyle = textStyle.Bold(true)
		}

		text = textStyle.Render(text)

		age := FormatTaskAge(todo.CreatedAt)
		fadedAge := ageStyle.Render(age)

		idStr := fmt.Sprintf("%*d", len(fmt.Sprintf("%d", maxID)), todo.ID)
		if m.cursor == i && (m.mode == Normal || m.mode == Edit) {
			idStr = cursorIDStyle.Render(idStr)
		} else {
			idStr = idStyle.Render(idStr)
		}

		textVisibleWidth := lipgloss.Width(text)
		textPad := width - textVisibleWidth
		if textPad < 0 {
			textPad = 0
		}

		var line string
		if m.cursor == i && (m.mode == Normal || m.mode == Edit) {
			line = fmt.Sprintf(" %s%s%s %s%s", idStr, cursor, check, text, strings.Repeat(" ", textPad)+fadedAge)
		} else {
			line = fmt.Sprintf("%s %s%s %s%s", idStr, cursor, check, text, strings.Repeat(" ", textPad)+fadedAge)
		}

		if m.cursor == i {
			line = lineHighlight.Render(line)
		}

		b.WriteString(line + "\n")
	}

	return b.String()
}
