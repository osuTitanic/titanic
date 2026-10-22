// Derived from https://github.com/roylee0704/gron (MIT License)
// Copyright (c) 2015 Roy Lee <roylee0704@gmail.com>
package scheduler

import (
	"fmt"
	"log/slog"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/discord"
	"github.com/osuTitanic/titanic/internal/state"
)

// Executor is a function that performs the task.
// It receives the application state and a logger as arguments.
type Executor func(state *state.State, logger *slog.Logger) error

// Name returns the name of the function that implements the Executor.
func (execute Executor) Name() string {
	// https://stackoverflow.com/questions/10742749/get-name-of-function-using-reflection
	function := runtime.FuncForPC(reflect.ValueOf(execute).Pointer())
	if function == nil {
		return "Unknown"
	}

	full := function.Name() // returns the full package path + . + function name
	parts := strings.Split(full, ".")
	return parts[len(parts)-1]
}

// Run executes a task with panic / error reports.
func (execute Executor) Run(app *state.State, logger *slog.Logger) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("task panicked: %v", recovered)
			if app.Officer == nil {
				return
			}
			app.Officer.CallPanic(
				recovered,
				discord.Field{Name: "Service", Value: "jobs", Inline: true},
				discord.Field{Name: "Task", Value: execute.Name(), Inline: true},
			)
			return
		}

		if err != nil && app.Officer != nil {
			embed := discord.Embed{Title: new("Task Failed"), Color: new(0xFF0000)}
			embed.AddField("Error", err.Error(), false)
			embed.AddField("Service", "jobs", true)
			embed.AddField("Task", execute.Name(), true)
			app.Officer.Call(discord.OfficerTagGeneral, "", embed)
		}
	}()

	return execute(app, logger)
}

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
