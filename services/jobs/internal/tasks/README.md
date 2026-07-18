# Tasks

This module contains the background tasks that can be run by the jobs service. Tasks receive the application state / logger, returning an `error` when they fail.

## Adding a Task

Add the task implementation to this directory, preferrably in a file named after the task, e.g. `users_cleanup.go`:

```go
func CleanupUsers(app *state.State, logger *slog.Logger) error {
	logger.Info("Cleaning up users")
	// ...
	return nil
}
```

Then register it in `services/jobs/cmd/cli/tasks.go`:

```go
var availableTasks = TaskList{
	// ...
	"users_cleanup": TaskWithoutArguments(tasks.CleanupUsers),
}
```

The `users_cleanup` key is the task name used by the CLI and schedule file.
Tasks without options/args should use `TaskWithoutArguments`.

Run the task directly with:

```sh
go run ./services/jobs/cmd/cli -name users_cleanup
```

## Task Options

Task options allow tasks to be configured with flags.

Define task options next to the task & validate them before doing any work. They would technically already be validated by the CLI builder, but this ensures that scheduled tasks (through the json) are also validated.

```go
type CleanupUsersOptions struct {
	BatchSize int
}

func (o CleanupUsersOptions) Validate() error {
	if o.BatchSize < 1 {
		return fmt.Errorf("batch size must be greater than zero")
	}
	return nil
}

func CleanupUsers(app *state.State, logger *slog.Logger, options CleanupUsersOptions) error {
	if err := options.Validate(); err != nil {
		return fmt.Errorf("invalid cleanup options: %w", err)
	}
	// ...
	return nil
}
```

In `services/jobs/cmd/cli/tasks.go`, add a builder that parses the task-specific flags and returns a `scheduler.Executor`:

```go
func BuildCleanupUsersTask(args []string) (scheduler.Executor, error) {
	// The `args` will contain the task-specific flags, excluding the
	// global CLI flags (such as `-name` or `-schedule-file`)

	// Create a new flag set to parse the task-specific flags
	flags := flag.NewFlagSet("users_cleanup", flag.ContinueOnError)
	batchSize := flags.Int("batch-size", 100, "number of users processed at once")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	if flags.NArg() > 0 {
		return nil, fmt.Errorf("unexpected arguments: %s", strings.Join(flags.Args(), " "))
	}

	// Create & validate the task options from the parsed flags
	options := tasks.CleanupUsersOptions{BatchSize: *batchSize}
	if err := options.Validate(); err != nil {
		return nil, err
	}

	// Return a `scheduler.Executor` that the scheduler can
	// call to run the task with the parsed options
	return func(app *state.State, logger *slog.Logger) error {
		return tasks.CleanupUsers(app, logger, options)
	}, nil
}
```

Register the builder instead of using `TaskWithoutArguments`:

```go
"users_cleanup": {Build: BuildCleanupUsersTask},
```

### Using Task Options

The service will try its best to parse various flag combinations:

```sh
go run ./services/jobs/cmd/cli -name users_cleanup -batch-size 250
```

If you want a consistent approach, use `--` to separate global flags from task flags:

```sh
go run ./services/jobs/cmd/cli -name users_cleanup -- -batch-size 250
```

For scheduled tasks, pass the same flags through the `args` array:

```json
{
    "name": "users_cleanup",
    "interval": 3600,
    "interval_at": "",
    "args": ["-batch-size", "250"]
}
```

## Future Considerations

It would be nice to eventually have [struct tags](https://stackoverflow.com/questions/10858787/what-are-the-uses-for-struct-tags-in-go) for the task options to replace the manual flag parsing. For example:

```go
type CleanupUsersOptions struct {
    BatchSize int `flag:"batch-size" default:"100" description:"number of users processed at once"`
}
```

This would allow the CLI to automatically generate the flag set and parse the flags into the struct, reducing boilerplate code. However, I am simply just too lazy to implement this right now.
