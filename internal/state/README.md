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
- `Extensions` for service-specific dependencies
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

## Extensions

Extensions allow services to attach their own dependencies to `State` without adding every service-specific type to the shared state struct.
They are stored by string key and retrieved with the expected Go type.

Stern uses this for the wiki service.
The service is created during startup in [main.go](https://github.com/osuTitanic/titanic-go/blob/main/services/stern/cmd/web/main.go#L183) and registered under the `"wiki"` key:

```go
wikiService := wiki.NewService(
	app.Config,
	app.Repositories,
	slog.Default().With("component", "wiki"),
)
state.RegisterExtension(app, "wiki", wikiService)
```

Routes can later resolve the typed service from the request state:

```go
service, ok := state.GetExtension[*wiki.Service](ctx.State, "wiki")
if !ok {
	return errors.New("wiki service not available")
}

// ...
```

`GetExtension` returns `false` when the key is missing or the stored value does not match the requested type.
Extensions are only stored on the state object. If an extension owns resources that need cleanup, the service that registers it should close them during shutdown. In the future, I might add a way to do cleanup extensions through a `Close` method or something similar.

<!-- TODO: Document repositories, transactions & lifecycle -->
