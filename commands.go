package main

import (
	"fmt"
	"strings"
	"time"
)

// ExecuteCommand processes a command string and returns updated tasks, filter, and message
func ExecuteCommand(tasks []Task, cmd string, filter string) ([]Task, string, string) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return tasks, filter, ""
	}

	switch parts[0] {
	case "add":
		if len(parts) > 1 {
			newTask := Task{
				Name:      strings.Join(parts[1:], " "),
				Priority:  "low",
				Status:    "active",
				CreatedAt: time.Now(),
			}
			tasks = append(tasks, newTask)
			return tasks, filter, "Added task."
		}
	case "remove":
		if len(parts) > 1 {
			idx := parseIndex(parts[1])
			if idx >= 0 && idx < len(tasks) {
				tasks = append(tasks[:idx], tasks[idx+1:]...)
				return tasks, filter, "Removed task."
			}
		}
	case "priority":
		if len(parts) > 2 {
			idx := parseIndex(parts[1])
			if idx >= 0 && idx < len(tasks) {
				tasks[idx].Priority = parts[2]
				return tasks, filter, "Changed priority."
			}
		}
	case "filter":
		if len(parts) > 1 {
			return tasks, parts[1], "Filter applied."
		}
		return tasks, "", "Filter cleared."
	case "edit":
		if len(parts) > 2 {
			idx := parseIndex(parts[1])
			if idx >= 0 && idx < len(tasks) {
				tasks[idx].Name = strings.Join(parts[2:], " ")
				return tasks, filter, "Task edited."
			}
		}
	}
	return tasks, filter, "Unknown command."
}

// parseIndex converts a 1-indexed string to 0-indexed int
func parseIndex(s string) int {
	var idx int
	_, err := fmt.Sscanf(s, "%d", &idx)
	if err != nil {
		return -1
	}
	return idx - 1
}
