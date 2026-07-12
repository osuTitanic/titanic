package main

import (
	"fmt"

	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/jobs/internal/scheduler"
)

func RunSingleTask(app *state.State, name string, interval int, intervalAt string, args []string) error {
	if interval > 0 {
		s := scheduler.New()
		if err := ScheduleTask(app, s, name, interval, intervalAt, args); err != nil {
			return err
		}
		StartSchedulerAndWait(app, s)
		return nil
	}

	taskFunc, err := availableTasks.Build(name, args)
	if err != nil {
		return err
	}

	app.Logger.Info("Running task", "name", name)
	if err = taskFunc(app, app.Logger); err != nil {
		return fmt.Errorf("task failed: %w", err)
	}

	app.Logger.Info("Done.")
	return nil
}
