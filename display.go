package main

import (
	"fmt"
	"time"
)

// FormatTaskAge returns a human-readable age string (e.g., "5m", "2h", "3d")
func FormatTaskAge(createdAt time.Time) string {
	duration := time.Since(createdAt)
	minutes := int(duration.Minutes())
	hours := int(duration.Hours())
	days := int(duration.Hours() / 24)

	if days > 0 {
		return fmt.Sprintf("%dd", days)
	} else if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	} else {
		return "just now"
	}
}

// DisplayTasks prints the task list to stdout in non-interactive mode
// mode can be: "active" (show only active), "deleted" (show only deleted), "all" (show everything)
func DisplayTasks(tasks []Task, filter string, mode string) {
	// Separate tasks by status
	var activeTasks, doneTasks, deletedTasks []Task
	for _, t := range tasks {
		if filter != "" && t.Priority != filter {
			continue
		}
		switch t.Status {
		case "deleted":
			deletedTasks = append(deletedTasks, t)
		case "done":
			doneTasks = append(doneTasks, t)
		default: // "active" or empty (for backward compatibility)
			activeTasks = append(activeTasks, t)
		}
	}

	// Display based on mode
	switch mode {
	case "active":
		// Show only active tasks
		if len(activeTasks) > 0 {
			fmt.Println("\nActive Tasks:")
			// Find original task indices
			taskIdx := 0
			for i, t := range tasks {
				if filter != "" && t.Priority != filter {
					continue
				}
				if t.Status != "active" && t.Status != "" {
					continue
				}
				age := ""
				if !t.CreatedAt.IsZero() {
					age = fmt.Sprintf(" [%s]", FormatTaskAge(t.CreatedAt))
				}
				fmt.Printf("[ ] %d. %s (%s)%s\n", i+1, t.Name, t.Priority, age)
				taskIdx++
			}
			fmt.Println()
		} else {
			fmt.Println("\nNo active tasks.")
		}

	case "deleted":
		// Show only deleted tasks
		if len(deletedTasks) > 0 {
			fmt.Println("\nDeleted Tasks:")
			// Find original task indices
			for i, t := range tasks {
				if filter != "" && t.Priority != filter {
					continue
				}
				if t.Status != "deleted" {
					continue
				}
				age := ""
				if !t.CreatedAt.IsZero() {
					age = fmt.Sprintf(" [%s]", FormatTaskAge(t.CreatedAt))
				}
				fmt.Printf("[DELETED] %d. %s (%s)%s\n", i+1, t.Name, t.Priority, age)
			}
			fmt.Println()
		} else {
			fmt.Println("\nNo deleted tasks.")
		}

	default:
		// Show all tasks in sections (for interactive mode on quit)
		// Display active tasks
		hasActive := false
		for i, t := range tasks {
			if filter != "" && t.Priority != filter {
				continue
			}
			if t.Status == "active" || t.Status == "" {
				if !hasActive {
					fmt.Println("\n=== ACTIVE TASKS ===")
					hasActive = true
				}
				age := ""
				if !t.CreatedAt.IsZero() {
					age = fmt.Sprintf(" [%s]", FormatTaskAge(t.CreatedAt))
				}
				fmt.Printf("[ ] %d. %s (%s)%s\n", i+1, t.Name, t.Priority, age)
			}
		}

		// Display done tasks
		hasDone := false
		for i, t := range tasks {
			if filter != "" && t.Priority != filter {
				continue
			}
			if t.Status == "done" {
				if !hasDone {
					fmt.Println("\n=== DONE TASKS ===")
					hasDone = true
				}
				age := ""
				if !t.CreatedAt.IsZero() {
					age = fmt.Sprintf(" [%s]", FormatTaskAge(t.CreatedAt))
				}
				fmt.Printf("[✓] %d. %s (%s)%s\n", i+1, t.Name, t.Priority, age)
			}
		}

		// Display deleted tasks
		hasDeleted := false
		for i, t := range tasks {
			if filter != "" && t.Priority != filter {
				continue
			}
			if t.Status == "deleted" {
				if !hasDeleted {
					fmt.Println("\n=== DELETED TASKS ===")
					hasDeleted = true
				}
				age := ""
				if !t.CreatedAt.IsZero() {
					age = fmt.Sprintf(" [%s]", FormatTaskAge(t.CreatedAt))
				}
				fmt.Printf("[DELETED] %d. %s (%s)%s\n", i+1, t.Name, t.Priority, age)
			}
		}

		if len(activeTasks) == 0 && len(doneTasks) == 0 && len(deletedTasks) == 0 {
			fmt.Println("\nNo tasks found.")
		}
		fmt.Println()
	}
}
