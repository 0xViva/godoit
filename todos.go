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

func (m model) todosToStringPlain() string {
	var b strings.Builder

	maxID := 0
	maxTextLen := 0
	for _, t := range m.todos {
		if t.ID > maxID {
			maxID = t.ID
		}
		if l := len([]rune(t.Text)); l > maxTextLen {
			maxTextLen = l
		}
	}
	width := maxTextLen + paddingAfterText

	for _, todo := range m.todos {
		// ID
		idStr := fmt.Sprintf("%*d", len(fmt.Sprintf("%d", maxID)), todo.ID)

		checkBox := checkBox
		if todo.Done {
			checkBox = fmt.Sprintf("[%s]", checkMark)
		}

		text := todo.Text
		if todo.Done {
			text = strikethroughOn + text + strikethroughOff
		}

		textPad := max(0, width-len([]rune(todo.Text)))

		age := FormatTaskAge(todo.CreatedAt)

		line := fmt.Sprintf("%s %s %s%s", idStr, checkBox, text, strings.Repeat(" ", textPad)+age)
		b.WriteString(line + "\n")
	}

	return b.String()
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
	width := maxLen + paddingAfterText

	for i, todo := range m.todos {
		cursor := " "
		if m.cursorLine == i && (m.mode == Normal || m.mode == Edit) {
			cursor = taskCursor
		}

		checkbox := checkBox
		if todo.Done {

			checkbox = fmt.Sprintf("[%s]", checkMark)
		}

		text := todo.Text
		textStyle := lipgloss.NewStyle()

		if todo.Done {
			textStyle = textStyle.Strikethrough(true)
		}

		if m.cursorLine == i && (m.mode == Normal || m.mode == Edit) {
			textStyle = textStyle.Bold(true)
		}

		if m.mode == Edit && i == m.cursorLine {
			cursor = " "
			if m.cursorVisible {
				cursor = inputCursor
			}
			text = m.input + inputCursor
			textStyle = textStyle.Bold(true)
		}

		text = textStyle.Render(text)

		age := FormatTaskAge(todo.CreatedAt)
		fadedAge := ageStyle.Render(age)

		idStr := fmt.Sprintf("%*d", len(fmt.Sprintf("%d", maxID)), todo.ID)
		if m.cursorLine == i && (m.mode == Normal || m.mode == Edit) {
			idStr = cursorIDStyle.Render(idStr)
		} else {
			idStr = idStyle.Render(idStr)
		}

		textVisibleWidth := lipgloss.Width(text)

		textPad := max(0, width-textVisibleWidth)

		var line string
		if m.cursorLine == i && (m.mode == Normal || m.mode == Edit) {
			line = fmt.Sprintf(" %s%s%s %s%s", idStr, cursor, checkbox, text, strings.Repeat(" ", textPad)+fadedAge)
		} else {
			line = fmt.Sprintf("%s %s%s %s%s", idStr, cursor, checkbox, text, strings.Repeat(" ", textPad)+fadedAge)
		}

		if m.cursorLine == i {
			line = lineHighlight.Render(line)
		}

		b.WriteString(line + "\n")
	}

	return b.String()
}
