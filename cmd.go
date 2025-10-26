package main

import (
	"flag"
	"fmt"
)

func RunCmd() bool {
	listFlag := flag.Bool("ls", false, "List all todos")
	flag.Parse()

	switch {
	case *listFlag:
		m := model{
			todos: loadTodos(),
			mode:  Normal,
		}

		if len(m.todos) == 0 {
			fmt.Println("No todos found.")
			return true
		}

		fmt.Println(m.todosToString())
		return true
	default:
		return false
	}
}
