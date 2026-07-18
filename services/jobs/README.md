# Jobs

Jobs is the background task runner for Titanic. It either runs one task directly or executes multiple tasks through the scheduler.

- See [internal/tasks](internal/tasks/README.md) for adding / configuring tasks
- See [internal/scheduler](internal/scheduler/README.md) for how the scheduler works

## Usage

Run commands from the repository root so the service can load the `.env`.

List all registered tasks:

```sh
go run ./services/jobs/cmd/cli -list
```

Run a task once:

```sh
go run ./services/jobs/cmd/cli -name stats_website
```

Task-specific options can follow the global flags:

```sh
go run ./services/jobs/cmd/cli -name ranks_index -force -workers 4
```

Use `-interval` to keep running a single task at an interval in seconds. For daily or longer intervals, `-interval-at` sets the time of day to run the task in `HH:MM` format:

```sh
go run ./services/jobs/cmd/cli -name ranks_unlist -interval 86400 -interval-at 11:00
```

## Schedule File

Use `-file` to run multiple tasks from a json schedule file:

```sh
go run ./services/jobs/cmd/cli -file services/jobs/schedule.example.json
```

The schedule is an array where each entry contains a registered task name and its schedule:

```json
[
    {
        "name": "ppv1_updates",
        "interval": 86400,
        "interval_at": "09:00",
        "args": []
    }
]
```

The service continues running until it receives a keyboard interrupt, then stops the scheduler gracefully.  
The Docker image copies `schedule.example.json` to `/app/schedule.json` and runs it by default.
