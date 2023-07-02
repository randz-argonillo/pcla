package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type List []item

func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*l = append(*l, t)
}

func (l *List) Complete(i int) error {
	list := *l

	if i <= 0 || i > len(list) {
		return fmt.Errorf("item %d does not exist. ", i)
	}

	list[i-1].Done = true
	list[i-1].CompletedAt = time.Now()

	return nil
}

func (l *List) Delete(i int) error {
	list := *l

	if i <= 0 || i > len(list) {
		return fmt.Errorf("item %d does not exist", i)
	}

	before := list[:i-1]
	after := list[i:]

	*l = append(before, after...)

	return nil
}

func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, js, 0644)
}

func (l *List) Get(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(content) == 0 {
		return nil
	}

	return json.Unmarshal(content, l)
}

func (l *List) String() string {
	format := ""

	for i, t := range *l {
		doneMark := "  "
		if t.Done {
			doneMark = "X "
		}

		format += fmt.Sprintf("%s%d: %s\n", doneMark, i+1, t.Task)
	}

	return format
}

func (l *List) AllDetails(pendingOnly bool) string {
	format := ""

	for i, t := range *l {
		if pendingOnly && t.Done {
			continue // exclude done task
		}

		doneMark := "  "
		doneDate := ""
		if t.Done {
			doneMark = "X "
			doneDate = fmt.Sprintf("done at %s", t.CompletedAt.Format("Jan 2"))
		}

		format += fmt.Sprintf("%s%d: %s (created at %s, done at %s)\n",
			doneMark,
			i+1,
			t.Task,
			t.CreatedAt.Format("Jan 2"),
			doneDate)
	}

	return format
}
