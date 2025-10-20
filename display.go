package main

import "fmt"

// DisplayTasks prints the task list to stdout in non-interactive mode
func DisplayTasks(tasks []Task, filter string) {
	fmt.Println("TODO List:")

	tasksToShow := tasks
	if filter != "" {
		var filtered []Task
		for _, t := range tasks {
			if t.Priority == filter {
				filtered = append(filtered, t)
			}
		}
		tasksToShow = filtered
	}

	for i, t := range tasksToShow {
		status := "[ ]"
		if t.Done {
			status = "[x]"
		}
		// Always start printing at column 0
		fmt.Printf("%d. %s %s (%s)\n", i+1, status, t.Name, t.Priority)
	}
}
