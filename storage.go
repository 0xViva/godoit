package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

func getTaskFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "todos.json"
	}
	return filepath.Join(home, "todos.json")
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
	if err != nil {
		return nil, err
	}

	// Migrate old tasks to new format
	for i := range tasks {
		if tasks[i].Status == "" {
			if tasks[i].Done {
				tasks[i].Status = "done"
			} else {
				tasks[i].Status = "active"
			}
		}
		// Set CreatedAt if missing (backward compatibility)
		if tasks[i].CreatedAt.IsZero() {
			tasks[i].CreatedAt = time.Now()
		}
	}

	return tasks, nil
}

// RemoveDoneTasks now keeps all tasks (for backward compatibility)
// Tasks are marked with status instead of being removed
func RemoveDoneTasks(tasks []Task) []Task {
	// Keep all tasks - the status field tracks active/done/deleted
	return tasks
}
