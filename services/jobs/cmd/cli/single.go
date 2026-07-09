package main

import (
	"fmt"

	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/jobs/internal/scheduler"
)

func RunSingleTask(app *state.State, name string, interval int, intervalAt string) error {
	taskFunc, ok := availableTasks[name]
	if !ok {
		return fmt.Errorf("unknown task: %s", name)
	}

	if interval > 0 {
		s := scheduler.New()
		ScheduleTask(app, s, name, interval, intervalAt)
		StartSchedulerAndWait(app, s)
		return nil
	}

	app.Logger.Info("Running task", "name", name)
	if err := taskFunc(app, app.Logger); err != nil {
		return fmt.Errorf("task failed: %w", err)
	}

	app.Logger.Info("Done.")
	return nil
}
