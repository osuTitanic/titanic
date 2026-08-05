package main

import (
	"flag"
	"fmt"
	"log/slog"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/jobs/internal/scheduler"
	"github.com/osuTitanic/titanic/services/jobs/internal/tasks"
)

// TaskDefinition builds a scheduler executor from task-specific arguments
type TaskDefinition struct {
	Build func(args []string) (scheduler.Executor, error)
}

type TaskList map[string]TaskDefinition

var availableTasks = TaskList{
	"ranks_index":         {Build: BuildIndexRanksTask},
	"ranks_sync":          TaskWithoutArguments(tasks.UpdateRanks),
	"ranks_unlist":        TaskWithoutArguments(tasks.UnlistUsers),
	"stats_website":       TaskWithoutArguments(tasks.UpdateWebsiteStats),
	"stats_activity":      TaskWithoutArguments(tasks.UpdateActivityStats),
	"users_notifications": TaskWithoutArguments(tasks.UpdateUnreadDmNotifications),
	"users_historical":    TaskWithoutArguments(tasks.FixHistoricalData),
	"users_autodelete":    TaskWithoutArguments(tasks.AutoDeleteUsers),
	"beatmap_statuses":    TaskWithoutArguments(tasks.UpdateBeatmapStatuses),
	"ppv1_updates":        TaskWithoutArguments(tasks.UpdatePPv1),
	"release_updates":     TaskWithoutArguments(tasks.ReleaseUpdates),
	"scores_status_pp":    {Build: BuildRecalculatePPStatusTask},
	"scores_status_score": {Build: BuildRecalculateScoreStatusTask},
}

func (t TaskList) List() {
	slog.Info("Available tasks:")
	for name := range t {
		slog.Info(fmt.Sprintf(" - %s", name))
	}
}

func (t TaskList) Build(name string, args []string) (scheduler.Executor, error) {
	definition, ok := t[name]
	if !ok {
		return nil, fmt.Errorf("unknown task: %s", name)
	}

	executor, err := definition.Build(args)
	if err != nil {
		return nil, fmt.Errorf("invalid arguments for task %s: %w", name, err)
	}
	return executor, nil
}

func TaskWithoutArguments(executor scheduler.Executor) TaskDefinition {
	return TaskDefinition{
		Build: func(args []string) (scheduler.Executor, error) {
			if len(args) > 0 {
				return nil, fmt.Errorf("task does not accept arguments: %s", strings.Join(args, " "))
			}
			return executor, nil
		},
	}
}

func BuildIndexRanksTask(args []string) (scheduler.Executor, error) {
	options, err := ParseIndexRanksOptions(args)
	if err != nil {
		return nil, err
	}

	return func(app *state.State, logger *slog.Logger) error {
		return tasks.IndexRanks(app, logger, options)
	}, nil
}

func ParseIndexRanksOptions(args []string) (tasks.IndexRanksOptions, error) {
	flags := flag.NewFlagSet("ranks_index", flag.ContinueOnError)
	force := flags.Bool("force", false, "index ranks even when the leaderboard is not empty")
	workers := flags.Int("workers", 0, "number of workers; 0 selects from application configuration")

	if err := flags.Parse(args); err != nil {
		return tasks.IndexRanksOptions{}, err
	}
	if flags.NArg() > 0 {
		return tasks.IndexRanksOptions{}, fmt.Errorf("unexpected arguments: %s", strings.Join(flags.Args(), " "))
	}

	options := tasks.IndexRanksOptions{
		Force:   *force,
		Workers: *workers,
	}
	if err := options.Validate(); err != nil {
		return tasks.IndexRanksOptions{}, err
	}
	return options, nil
}

func BuildRecalculatePPStatusTask(args []string) (scheduler.Executor, error) {
	options, err := ParseStatusRecalculationOptions("scores_status_pp", args)
	if err != nil {
		return nil, err
	}

	return func(app *state.State, logger *slog.Logger) error {
		return tasks.RecalculatePPStatusUser(app, logger, options)
	}, nil
}

func BuildRecalculateScoreStatusTask(args []string) (scheduler.Executor, error) {
	options, err := ParseStatusRecalculationOptions("scores_status_score", args)
	if err != nil {
		return nil, err
	}

	return func(app *state.State, logger *slog.Logger) error {
		return tasks.RecalculateScoreStatusUser(app, logger, options)
	}, nil
}

func ParseStatusRecalculationOptions(taskName string, args []string) (tasks.StatusRecalculationOptions, error) {
	flags := flag.NewFlagSet(taskName, flag.ContinueOnError)
	userId := flags.Int("user-id", 0, "ID of the user whose score statuses should be recalculated")
	mode := flags.Int("mode", 0, "game mode ID")

	if err := flags.Parse(args); err != nil {
		return tasks.StatusRecalculationOptions{}, err
	}
	if flags.NArg() > 0 {
		return tasks.StatusRecalculationOptions{}, fmt.Errorf("unexpected arguments: %s", strings.Join(flags.Args(), " "))
	}

	options := tasks.StatusRecalculationOptions{
		UserId: *userId,
		Mode:   constants.Mode(*mode),
	}
	if err := options.Validate(); err != nil {
		return tasks.StatusRecalculationOptions{}, err
	}
	return options, nil
}
