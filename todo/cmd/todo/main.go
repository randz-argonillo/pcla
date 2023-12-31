package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pcla/todo"
)

var todoFilename = ".todo.json"

func main() {
	if os.Getenv("TODO_FNAME") != "" {
		todoFilename = os.Getenv("TODO_FNAME")
	}

	l := &todo.List{}

	addFlag := flag.Bool("add", false, "Add task flag")
	delete := flag.Int("del", 0, "Delete a task item")

	listFlag := flag.Bool("list", false, "List all incomplete tasks")
	listVerboseFlag := flag.Bool("verbose", false, "List tasks verbosely")
	pendingFlag := flag.Bool("pending", false, "List pending tasks")

	complete := flag.Int("complete", 0, "The task item to be completed")

	flag.Parse()

	if err := l.Get(todoFilename); err != nil {
		printErrorAndExit(err)
	}

	switch {
	case *listVerboseFlag && *listFlag && *pendingFlag:
		fmt.Println(l.AllDetails(true))
	case *listVerboseFlag && *listFlag:
		fmt.Println(l.AllDetails(false))
	case *listFlag && *pendingFlag:
		fmt.Println(l.AllDetails(true))
	case *listFlag:
		fmt.Println(l)
	case *addFlag:
		t, err := getTask(os.Stdin, flag.Args()...)

		if err != nil {
			printErrorAndExit(err)
		}

		l.Add(t)

		if err = l.Save(todoFilename); err != nil {
			printErrorAndExit(err)
		}
	case *delete > 0:
		if err := l.Delete(*delete); err != nil {
			printErrorAndExit(err)
		}

		if err := l.Save(todoFilename); err != nil {
			printErrorAndExit(err)
		}
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			printErrorAndExit(err)
		}

		if err := l.Save(todoFilename); err != nil {
			printErrorAndExit(err)
		}
	default:
		fmt.Fprintln(os.Stderr, "Invalid command line option")
		os.Exit(1)
	}
}

func printErrorAndExit(err error) {
	fmt.Fprint(os.Stderr, err)
	os.Exit(1)
}

func getTask(reader io.Reader, args ...string) (string, error) {
	task := strings.Join(args, " ")
	if task != "" {
		return task, nil
	}

	task, err := getTaskFromIO(reader)
	if err != nil {
		return "", err
	}

	return task, nil
}

func getTaskFromIO(reader io.Reader) (string, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Scan()

	if err := scanner.Err(); err != nil {
		return "", err
	}

	t := scanner.Text()
	if len(t) == 0 {
		return "", fmt.Errorf("you added an empty todo")
	}

	return t, nil
}
