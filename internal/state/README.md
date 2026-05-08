# State

This module wires the shared application state used by services.
Most services should start here instead of manually creating each dependency.

## Components

`State` exposes various shared dependencies:

- `Config` for parsed application configuration
- `Logger` for the default component logger
- `Database` for the underlying `*gorm.DB` session
- `Redis` for the shared Redis client
- `Storage` for app storage
- `Email` for email delivery
- `Repositories` for actual database access
- `Rankings` for ranking services
- `PPv1` for ppv1 score calculation helpers

## Usage

Create the state once during service startup and close it when the service exits.

```go
app, err := state.NewState()
if err != nil {
	return err
}
defer app.Close()

user, err := app.Repositories.Users.ById(userId, "Stats")
if err != nil {
	return err
}

app.Logger.Info("User loaded", "id", user.Id)
```

You can pass dotenv file paths through to the config loader when a service needs a specific environment file.

```go
app, err := state.NewState(".env", ".env.local")
if err != nil {
	return err
}
defer app.Close()
```

<!-- TODO: Document repositories, transactions & lifecycle -->
