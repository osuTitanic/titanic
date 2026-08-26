// Derived from https://github.com/roylee0704/gron (MIT License)
// Copyright (c) 2015 Roy Lee <roylee0704@gmail.com>
package scheduler

import (
	"log/slog"
	"slices"
	"time"

	"github.com/osuTitanic/titanic/internal/state"
)

// Scheduler provides a convenient interface for scheduling tasks. It keeps track
// of any number of entries, invoking the associated func as specified by the schedule.
// It may also be started, stopped and the entries may be inspected.
type Scheduler struct {
	logger  *slog.Logger
	entries []*Task
	running bool
	add     chan *Task
	stop    chan struct{}
}

// New instantiates new Scheduler
func New() *Scheduler {
	return &Scheduler{
		logger: slog.Default().With("component", "scheduler"),
		stop:   make(chan struct{}),
		add:    make(chan *Task),
	}
}

// Start the scheduler in its own go-routine.
func (c *Scheduler) Start(app *state.State) {
	if c.running {
		return
	}
	c.running = true
	go c.run(app)
}

// Stop the scheduler.
func (c *Scheduler) Stop() {
	if !c.running {
		return
	}
	c.running = false
	c.stop <- struct{}{}
}

// Add adds a new task to the scheduler.
func (c *Scheduler) Add(schedule Schedule, executor Executor) *Task {
	task := &Task{
		Schedule: schedule,
		Executor: executor,
	}
	c.AddTask(task)
	return task
}

// AddTask adds a new task to the scheduler.
func (c *Scheduler) AddTask(task *Task) {
	if !c.running {
		c.entries = append(c.entries, task)
		return
	}
	c.add <- task
}

func (c *Scheduler) run(app *state.State) {
	sortByTime := func(a, b *Task) int {
		aIsZero := a.Next.IsZero()
		bIsZero := b.Next.IsZero()

		switch {
		case aIsZero && bIsZero:
			return 0
		case aIsZero:
			return 1
		case bIsZero:
			return -1
		default:
			return a.Next.Compare(b.Next)
		}
	}
	now := time.Now()

	for _, entry := range c.entries {
		entry.Next = entry.Schedule.Next(now)
	}

	for {
		slices.SortFunc(c.entries, sortByTime)

		// Set the timer to the next task's scheduled time, or a
		// default of 1 hour if there are no tasks
		delay := time.Hour
		if len(c.entries) > 0 && !c.entries[0].Next.IsZero() {
			delay = time.Until(c.entries[0].Next)
		}
		timer := time.NewTimer(delay)

		select {
		case now = <-timer.C:
			for _, entry := range c.entries {
				if entry.Next.After(now) || entry.Next.IsZero() {
					break
				}
				go c.runTask(app, entry)
				entry.Prev = entry.Next
				entry.Next = entry.Schedule.Next(now)
			}
		case task := <-c.add:
			timer.Stop()
			task.Next = task.Schedule.Next(time.Now())
			c.entries = append(c.entries, task)
		case <-c.stop:
			timer.Stop()
			return
		}
	}
}

func (c *Scheduler) runTask(app *state.State, task *Task) {
	if task.Logger == nil {
		task.Logger = slog.Default().With("component", "tasks")
	}

	defer func() {
		if r := recover(); r != nil {
			task.Logger.Error("task panicked", slog.Any("panic", r))
		}
	}()

	err := task.Executor(app, task.Logger)
	if err != nil {
		task.Logger.Error("task failed", slog.String("error", err.Error()))
	}
}
