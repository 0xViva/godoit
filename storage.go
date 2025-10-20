package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func getTaskFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home directory can't be determined
		return ".godoit.json"
	}
	return filepath.Join(home, ".godoit.json")
}

// SaveTasks saves the task list to disk
func SaveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(getTaskFile(), data, 0644)
}

// LoadTasks loads the task list from disk
func LoadTasks() ([]Task, error) {
	data, err := os.ReadFile(getTaskFile())
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}

// RemoveDoneTasks filters out completed tasks
func RemoveDoneTasks(tasks []Task) []Task {
	var remaining []Task
	for _, t := range tasks {
		if !t.Done {
			remaining = append(remaining, t)
		}
	}
	return remaining
}
