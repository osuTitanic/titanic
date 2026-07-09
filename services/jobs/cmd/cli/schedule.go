package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/jobs/internal/scheduler"
)

type SchedulerConfig struct {
	Name       string  `json:"name"`
	Interval   *int    `json:"interval"` // 0 -> run once on startup, >0 -> recurring interval in seconds
	IntervalAt *string `json:"interval_at"`
}

func RunSchedulerFile(app *state.State, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read scheduler file: %w", err)
	}

	var configs []SchedulerConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return fmt.Errorf("parse scheduler file: %w", err)
	}
	s := scheduler.New()

	for _, config := range configs {
		if config.Interval == nil {
			app.Logger.Error("Interval missing", "name", config.Name)
			continue
		}

		intervalAt := ""
		if config.IntervalAt != nil {
			intervalAt = *config.IntervalAt
		}

		if err := ScheduleTask(app, s, config.Name, *config.Interval, intervalAt); err != nil {
			app.Logger.Error("Failed to schedule task", "name", config.Name, "error", err)
		}
	}

	StartSchedulerAndWait(app, s)
	return nil
}

func ScheduleTask(app *state.State, s *scheduler.Scheduler, name string, interval int, intervalAt string) error {
	taskFunc, ok := availableTasks[name]
	if !ok {
		return fmt.Errorf("unknown task: %s", name)
	}

	schedule, err := scheduleFromConfig(interval, intervalAt)
	if err != nil {
		return err
	}

	task := s.Add(schedule, taskFunc)
	task.SetLogger("tasks/" + name)

	if interval == 0 {
		app.Logger.Info("Scheduled startup task", "name", name)
		return nil
	}

	app.Logger.Info("Scheduled task", "name", name, "interval", interval, "interval_at", intervalAt)
	return nil
}

func scheduleFromConfig(interval int, intervalAt string) (scheduler.Schedule, error) {
	if interval < 0 {
		return nil, fmt.Errorf("interval must be >= 0")
	}

	if interval == 0 {
		// Run task once immediately on startup
		return scheduler.Now(), nil
	}

	period := time.Duration(interval) * time.Second
	schedule := scheduler.Every(period)

	if intervalAt == "" {
		return schedule, nil
	}

	if period < 24*time.Hour {
		return nil, fmt.Errorf("interval_at requires an interval of at least 86400 seconds")
	}
	if _, err := time.Parse("15:04", intervalAt); err != nil {
		return nil, fmt.Errorf("invalid interval_at %q: expected HH:MM", intervalAt)
	}

	return schedule.At(intervalAt), nil
}

func StartSchedulerAndWait(app *state.State, s *scheduler.Scheduler) {
	s.Start(app)
	app.Logger.Info("Scheduler started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	app.Logger.Info("Shutting down...")
	s.Stop()
}
