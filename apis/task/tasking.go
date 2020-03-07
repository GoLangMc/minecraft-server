package task

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/golangmc/minecraft-server/apis/base"
)

type Task struct {
	uuid int64
	exec *func(task *Task)

	cancel bool
	period int64
	paused int64
	tasker *Tasking
}

type Tasking struct {
	// milliseconds per tick
	mpt int64

	// uuid -> task
	tasks map[int64]*Task

	// task -> last ran
	ticks map[*Task]int64
	// time -> tasks
	queue map[int64][]*Task

	next uint64
	done bool
	kill chan bool
}

func NewTasking(mpt int64) *Tasking {
	return &Tasking{
		mpt: mpt,

		tasks: make(map[int64]*Task),
		ticks: make(map[*Task]int64),
		queue: make(map[int64][]*Task),
	}
}

func (t *Tasking) Load() {
	t.done = false
	t.kill = make(chan bool, 1)

	go t.tick()
}

func (t *Tasking) Kill() {
	if t.done {
		return
	}

	t.done = true
	t.kill <- true

	for k := range t.ticks {
		delete(t.ticks, k)
	}

	for k := range t.queue {
		delete(t.queue, k)
	}

	for k, v := range t.tasks {
		delete(t.tasks, k)
		v.Cancel()
	}

	close(t.kill)
}

func (t *Tasking) tick() {
	tick := time.NewTicker(time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-t.kill:
			return
		case curr := <-tick.C:
			t.tickQueue(curr)
			t.tickTasks(curr)
		}
	}
}

func (t *Tasking) tickTasks(curr time.Time) {
	unix := curr.UnixNano() / 1e6

	for task, last := range t.ticks {
		if unix-last < task.period {
			continue // not ready to be executed
		}

		if err := task.attemptExec(); err != nil {
			task.cancel = true
			fmt.Printf("%v", err)
		}

		if task.cancel || task.period <= 0 {
			delete(t.ticks, task)
		} else {
			t.ticks[task] = unix
		}
	}
}

func (t *Tasking) tickQueue(curr time.Time) {
	// this is for handling the delayed tasks
	// delay gets counted down, and then the tasks are added to the queue used for tickTasks

	unix := curr.UnixNano() / 1e6

	for when, tasks := range t.queue {
		if unix < when {
			continue // not ready to be executed
		}

		delete(t.queue, when)

		for _, task := range tasks {
			t.ticks[task] = 0
		}
	}
}

func (t *Tasking) nextTaskU() int64 {
	return int64(atomic.AddUint64(&t.next, 1))
}

func (t *Tasking) repeats(period int64, function func(task *Task)) {
	task := t.newTask(period, 0, &function)

	t.ticks[task] = 0
	t.tasks[task.uuid] = task
}

func (t *Tasking) delayed(paused int64, function func(task *Task)) {
	task := t.newTask(0, paused, &function)

	unix := time.Now().UnixNano() / 1e6
	when := unix + paused

	queue, exists := t.queue[when]

	if !exists {
		queue = make([]*Task, 0)
	}

	queue = append(queue, task)

	t.queue[when] = queue
	t.tasks[task.uuid] = task
}

// repeats the function every period, in ticks
func (t *Tasking) Every(period int64, function func(task *Task)) {
	t.repeats(period*t.mpt, function)
}

// executes the function after paused, in ticks
func (t *Tasking) After(paused int64, function func(task *Task)) {
	t.delayed(paused*t.mpt, function)
}

func (t *Tasking) EveryTime(period int64, duration time.Duration, function func(task *Task)) {
	t.repeats(period*duration.Milliseconds(), function)
}

func (t *Tasking) AfterTime(paused int64, duration time.Duration, function func(task *Task)) {
	t.delayed(paused*duration.Milliseconds(), function)
}

func (t *Tasking) newTask(period int64, paused int64, function *func(task *Task)) *Task {
	return &Task{
		tasker: t,
		uuid:   t.nextTaskU(),
		exec:   function,
		period: period,
		paused: paused,
	}
}

func (t *Task) attemptExec() (error error) {
	return base.Attempt(func() { (*t.exec)(t) })
}

func (t *Task) Cancel() {
	t.cancel = true
}

func (t *Task) Tasker() *Tasking {
	return t.tasker
}
