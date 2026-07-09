package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/jobs/internal/tasks"
)

type TaskList map[string]func(*state.State, *slog.Logger) error

var availableTasks = TaskList{
	"stats_website":       tasks.UpdateWebsiteStats,
	"stats_activity":      tasks.UpdateActivityStats,
	"users_notifications": tasks.UpdateUnreadDmNotifications,
	"users_historical":    tasks.FixHistoricalData,
	"users_autodelete":    tasks.AutoDeleteUsers,
	"ranks_sync":          tasks.UpdateRanks,
	"ranks_index":         tasks.IndexRanks,
	"ranks_unlist":        tasks.UnlistUsers,
	"beatmap_statuses":    tasks.UpdateBeatmapStatuses,
	"ppv1_updates":        tasks.UpdatePPv1,
	"release_updates":     tasks.ReleaseUpdates,
}

func (t *TaskList) List() {
	slog.Info("Available tasks:")
	for name := range availableTasks {
		slog.Info(fmt.Sprintf(" - %s", name))
	}
}

func main() {
	listFlag := flag.Bool("list", false, "list all tasks")
	fileFlag := flag.String("file", "", "specify a scheduler file")
	nameFlag := flag.String("name", "", "run a specific task by name")
	intervalFlag := flag.Int("interval", 0, "the interval to run the task in (seconds)")
	intervalAtFlag := flag.String("interval-at", "", "specify the time period at which the task should run (e.g. 15:00)")
	flag.Parse()

	if *listFlag {
		availableTasks.List()
		os.Exit(0)
	}

	app, err := state.NewState(".env")
	if err != nil {
		slog.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}
	defer app.Close()

	if *fileFlag != "" {
		if err := RunSchedulerFile(app, *fileFlag); err != nil {
			app.Logger.Error("Failed to run tasks from file", "error", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *nameFlag != "" {
		err := RunSingleTask(app, *nameFlag, *intervalFlag, *intervalAtFlag)
		if err != nil {
			app.Logger.Error("Failed to run task", "error", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	flag.Usage()
	os.Exit(1)
}
