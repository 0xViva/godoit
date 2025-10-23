package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var todosFile string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("cannot determine home directory")
	}

	// Use XDG_DATA_HOME or default to ~/.local/share
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome == "" {
		xdgDataHome = filepath.Join(home, ".local", "share")
	}

	dataDir := filepath.Join(xdgDataHome, "godoit")
	// Create directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		panic("cannot create data directory: " + err.Error())
	}

	todosFile = filepath.Join(dataDir, "todos.json")
}

// loadTodos reads todos from the todos.json file
func loadTodos() []Todo {
	data, err := os.ReadFile(todosFile)
	if err != nil {
		return []Todo{}
	}
	var todos []Todo
	_ = json.Unmarshal(data, &todos)
	return todos
}

// saveTodos writes todos to the todos.json file
func saveTodos(todos []Todo) {
	data, _ := json.MarshalIndent(todos, "", "  ")
	_ = os.WriteFile(todosFile, data, 0644)
}
