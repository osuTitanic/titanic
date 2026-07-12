package main

import (
	"errors"
	"flag"
	"log/slog"
	"os"

	"github.com/osuTitanic/titanic/internal/state"
)

func main() {
	listFlag := flag.Bool("list", false, "list all tasks")
	fileFlag := flag.String("file", "", "specify a scheduler file")
	nameFlag := flag.String("name", "", "run a specific task by name (task flags follow global flags)")
	intervalFlag := flag.Int("interval", 0, "the interval to run the task in (seconds)")
	intervalAtFlag := flag.String("interval-at", "", "specify the time period at which the task should run (e.g. 15:00)")

	// We want to split the global args from the task-specific args
	// Go doesn't provide us a way of ignoring unregistered args, so this is what we have to do
	globalArgs, taskArgs := splitCommandLineArgs(flag.CommandLine, os.Args[1:])

	// Parse only the currently registered (global) args
	if err := flag.CommandLine.Parse(globalArgs); err != nil {
		os.Exit(2)
	}

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
		err := RunSingleTask(app, *nameFlag, *intervalFlag, *intervalAtFlag, taskArgs)
		if err != nil {
			if errors.Is(err, flag.ErrHelp) {
				os.Exit(0)
			}
			app.Logger.Error("Failed to run task", "error", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	flag.Usage()
	os.Exit(1)
}
