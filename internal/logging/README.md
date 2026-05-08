# Logging

This module customizes the default `slog` logger & makes it more usable, at least in my opinion.
It formats log lines with timestamps, components, levels, colors & structured attributes.

## Usage

Use the default `slog` logger after the app state has been created.

```go
app, err := state.NewState()
if err != nil {
	return err
}
defer app.Close()

// use either `slog` or `app.Logger`
slog.Info("Service started", "port", 8080)
slog.Debug("Cache refreshed", "key", "beatmaps")
```

For standalone services, set the default logger directly.

```go
logging.SetDefault("titanic", slog.LevelInfo)
slog.Info("Worker started", "name", "score-processor")
```

This will override the default slog logger to a component logger with the given component name, level & writer(s).
It will always have a console writer by default.

### Custom loggers

Create a component logger when you need a logger instance without changing the process default.

```go
logger := logging.NewComponentLogger("api", slog.LevelDebug, logging.GetConsoleWriter())
logger.Info("Request handled", "status", 200, "path", "/users/123")
```

Another (easier) option is to use the default logger & apply a custom `component`.
This can be seen as the equivilent of python's `logging.getLogger("component")`.

```go
logger := slog.Default().With("component", "database")
logger.Info("Connection established", "host", "localhost", "port", 5432)
```

## Writers

Log output can be sent to stdout, files, or multiple writers.

```go
fileWriter, err := logging.GetFileWriter(".data/titanic.log")
if err != nil {
	return err
}

logging.SetDefault("titanic", slog.LevelInfo, fileWriter)
```

`SetDefault` always includes the console writer & appends any additional writers.

## Output

Log records include the timestamp, component, level, message & attributes.
If no `component` attribute is present, the component defaults to `titanic`.

```text
[2026-05-08 12:00:00] - <titanic> INFO: Scheduler started -> tasks=5
```
