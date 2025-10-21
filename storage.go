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

func SaveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(getTaskFile(), data, 0644)
}

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

func RemoveDoneTasks(tasks []Task) []Task {
	return tasks
}
