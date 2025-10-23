package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type Todo struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func (m model) todosToString() string {
	lines := []string{}
	for i, t := range m.todos {
		check := "[ ]"
		if t.Done {
			check = "[✔]"
		}

		cursor := "  "
		if m.mode == Edit && i == m.cursor {
			cursor = " ➤"
		} else if m.mode == Normal && i == m.cursor {
			cursor = " →"
		}

		text := t.Text
		if m.mode == Edit && i == m.cursor {
			text = lipgloss.NewStyle().Italic(true).Render(m.input)
			if m.cursorVisible {
				text += "_"
			}
		}

		// strike through only the text
		if t.Done && !(m.mode == Edit && i == m.cursor) {
			text = lipgloss.NewStyle().Strikethrough(true).Render(text)
		}

		line := fmt.Sprintf("%s%s %s", cursor, check, text)

		if i == m.cursor && (m.mode == Normal || m.mode == Edit) {
			line = lipgloss.NewStyle().Bold(true).Render(line)
		}

		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
