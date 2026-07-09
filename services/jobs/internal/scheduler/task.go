// Derived from https://github.com/roylee0704/gron (MIT License)
// Copyright (c) 2015 Roy Lee <roylee0704@gmail.com>
package scheduler

import (
	"log/slog"
	"time"

	"github.com/osuTitanic/titanic/internal/state"
)

// Executor is a function that performs the task.
// It receives the application state and a logger as arguments.
type Executor func(state *state.State, logger *slog.Logger) error

// Task consists of a Schedule and an Executor to be executed on that schedule.
type Task struct {
	Schedule Schedule
	Executor Executor
	Logger   *slog.Logger

	// The next time the task will run. This is zero time if Scheduler
	// has not been started or invalid schedule.
	Next time.Time

	// The last time the task was run. This is zero time if the
	// task has not been run.
	Prev time.Time
}

// SetLogger allows the task to use a logger with a specific component name
// which will make it easier to keep track of logs related to this task.
func (t *Task) SetLogger(name string) {
	t.Logger = slog.Default().With("component", name)
}
