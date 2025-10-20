package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-i":
			runInteractive()
			return
		case "-l":
			DisplayTasks(InitialModel().Tasks, "")
			return
		default:
			fmt.Println("Usage: terminal-todo [-i for interactive|-l for list]")
			return
		}
	}

	// Default: interactive mode
	runInteractive()
}

func runInteractive() {
	p := tea.NewProgram(InitialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
