# Task Scheduler

This module allows you to schedule tasks to run immediately, after a delay, periodically, or at specific times of the day, derived from [gron](https://github.com/roylee0704/gron).

## Usage

Create a new scheduler, start it, and ensure it gets stopped gracefully when your application shuts down.

```go
func main() {
	s := scheduler.New()

	// NOTE: Usually you would pass in the application state here, but you don't have to
	s.Start(nil)
	defer s.Stop()

	// add tasks here ...
}
```

You can add tasks using the `Add` method, which takes a Schedule definition and an Executor function.

1. Use `Now()` to execute a task as soon as possible.
```go
s.Add(scheduler.Now(), func(app *state.State, logger *slog.Logger) error {
	logger.Info("Running immediately!")
	return nil
})
```

2. Use `Once(duration)` to execute a task exactly once after the specified duration.
```go
s.Add(scheduler.Once(5*time.Second), func(app *state.State, logger *slog.Logger) error {
	logger.Info("Running once after 5 seconds!")
	return nil
})
```

3. Use `Every(duration)` to execute a task repeatedly at a fixed interval.
```go
s.Add(scheduler.Every(10*time.Second), func(app *state.State, logger *slog.Logger) error {
	logger.Info("Running every 10 seconds!")
	return nil
})
```

4. Use `At(timeString)` to schedule a task to run daily at a specific time of day. Format should be `"15:04"`.
```go
s.Add(scheduler.At("14:30"), func(app *state.State, logger *slog.Logger) error {
	logger.Info("Running every day at 14:30!")
	return nil
})
```
