package main

import (
	"github.com/charmbracelet/bubbles/help"
)

type model struct {
	todos         []Todo
	cursor        int
	mode          Mode
	input         string
	help          help.Model
	cursorVisible bool
}

func initialModel() model {
	todos := loadTodos()
	return model{
		todos:  todos,
		cursor: 0,
		mode:   Normal,
		input:  "",
		help:   help.New(),
	}
}
