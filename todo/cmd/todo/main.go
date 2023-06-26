package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pcla/todo"
)

const todoFilename = ".todo.json"

func main() {
	l := &todo.List{}

	if err := l.Get(todoFilename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case len(os.Args) == 1:
		for _, t := range *l {
			fmt.Println(t.Task)
		}
	default:
		task := strings.Join(os.Args[1:], " ")
		l.Add(task)

		if err := l.Save(todoFilename); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}
}
