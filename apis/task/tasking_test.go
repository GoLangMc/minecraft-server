package task

import (
	"fmt"
	"testing"
)

func TestTasker_Load(t *testing.T) {

	tasker := NewTasking(1_000 / 20)
	tasker.Load()

	tasker.Every(20, printCurrentTask)

	done := <-tasker.kill
	fmt.Printf("Tasker done: %t\n", done)
}

var count = 0

func printCurrentTask(_ *Task) {
	fmt.Printf("Running print task %d\n", count)

	if count++; count >= 2 {
		panic("hi")
	}
}
