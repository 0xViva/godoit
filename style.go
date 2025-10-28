package main

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	checkBox         = "[ ]"
	checkMark        = "x"
	strikethroughOn  = "\033[9m"
	strikethroughOff = "\033[0m"
	paddingAfterText = 3
	taskCursor       = "➤"
	inputCursor      = "_"
)

var (
	ageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Faint(true)

	idStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Faint(true)

	lineHighlight = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15"))

	cursorIDStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Bold(true)

	modeBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	inputStyle = lipgloss.NewStyle().
			Bold(true).
			Italic(true)

	inputIDStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Bold(true)

	header = lipgloss.NewStyle().Bold(true).Render("📝 GoDoIt")
)
