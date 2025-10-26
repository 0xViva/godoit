package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	if RunCmd() {
		os.Exit(0)
	}

	p := tea.NewProgram(initialModel())
	p.Run()
}
