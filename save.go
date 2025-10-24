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
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome == "" {
		xdgDataHome = filepath.Join(home, ".local", "share")
	}

	dataDir := filepath.Join(xdgDataHome, "godoit")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		panic("cannot create data directory: " + err.Error())
	}

	todosFile = filepath.Join(dataDir, "todos.json")
}

func loadTodos() []Todo {
	data, err := os.ReadFile(todosFile)
	if err != nil {
		return []Todo{}
	}
	var todos []Todo
	_ = json.Unmarshal(data, &todos)
	return todos
}

func saveTodos(todos []Todo) {
	data, _ := json.MarshalIndent(todos, "", "  ")
	_ = os.WriteFile(todosFile, data, 0644)
}
