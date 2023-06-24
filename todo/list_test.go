package todo_test

import (
	"os"
	"testing"

	"github.com/pcla/todo"
)

func TestAdd(t *testing.T) {
	l := todo.List{}

	task := "Clean house"
	l.Add(task)

	expected := l[0].Task

	if expected != task {
		t.Errorf("expected %q, got %q instead", task, expected)
	}
}

func TestComplete(t *testing.T) {
	l := todo.List{}
	task := "Clean house"
	l.Add(task)

	if l[0].Task != task {
		t.Errorf("expected %q, got %q instead", task, l[0].Task)
	}

	if l[0].Done {
		t.Errorf("Task should not be completed unless you complete it")
	}

	l.Complete(1)

	if !l[0].Done {
		t.Errorf("Task should be completed when you complete it")
	}
}

func TestDelete(t *testing.T) {
	tasks := []string{
		"Task 1",
		"Task 2",
		"Task 3",
	}

	l := todo.List{}
	for _, i := range tasks {
		l.Add(i)
	}

	if l[0].Task != tasks[0] {
		t.Errorf("expected %q, got %q instead", tasks[0], l[0].Task)
	}

	l.Delete(1)

	if len(l) != 2 {
		t.Errorf("expected length of List will reduce after delete")
	}

	if l[0].Task == tasks[0] {
		t.Errorf("expected the task is remove from the list but not")
	}

	if l[0].Task != tasks[1] {
		t.Errorf("expected the 2nd task will become the first task")
	}
}

func TestSaveGet(t *testing.T) {
	l1 := todo.List{}
	l2 := todo.List{}

	task := "Task 1"
	l1.Add(task)

	if l1[0].Task != task {
		t.Errorf("expected task %q, but got %q instead", task, l1[0].Task)
	}

	tf, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("error creating temp file: %s", err)
	}

	defer os.Remove(tf.Name())

	err = l1.Save(tf.Name())
	if err != nil {
		t.Fatalf("error saving todo list: %s", err)
	}

	err = l2.Get(tf.Name())
	if err != nil {
		t.Fatalf("error in getting todo list from file %s", tf.Name())
	}

	if len(l2) == 0 {
		t.Errorf("no todo list retrieve from file")
	}

	if l2[0].Task != task {
		t.Errorf("expected first task to be %q but got %q instead", task, l2[0].Task)
	}
}
