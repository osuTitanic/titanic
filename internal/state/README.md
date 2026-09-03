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
- `Officer` for staff webhook notifications
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

## Repositories

`State` creates one `Repositories` value for the main database session and exposes it as both `app.Repositories` and embedded repository fields like `app.Users`, `app.Scores`, and so on.
Use these repositories for normal service code instead of calling [GORM](https://gorm.io/) directly.

```go
user, err := app.Repositories.Users.ById(userId, "Stats", "Groups")
if err != nil {
	return err
}

scores, err := app.Scores.FetchBest(user.Id, constants.ModeOsu, true, "Beatmap")
if err != nil {
	return err
}
```

Repository methods usually accept optional preload names for relationships that should be loaded with the query.
For the repository method patterns, see [Repositories](../repositories/README.md).

## Transactions

Use `DatabaseTransaction` when multiple writes need to commit or roll back together.
The callback receives a new `Repositories` value bound to the transaction.

```go
err := app.DatabaseTransaction(func(repos *state.Repositories) error {
	user := &schemas.User{
		Name:     name,
		SafeName: schemas.ResolveSafeName(name),
		Email:    email,
	}
	if err := repos.Users.Create(user); err != nil {
		return err // will perform a rollback
	}

	entry := &schemas.GroupEntry{
		UserId:  user.Id,
		GroupId: constants.GroupPlayers,
	}
	if err := repos.Groups.CreateEntry(entry); err != nil {
		return err // will perform a rollback
	}

	return nil // will commit the transaction
})
if err != nil {
	return err
}
```

Return an error to roll back the transaction and return `nil` to commit it.
Inside the callback, use the `repos` argument instead of `app.Repositories` so every query participates in the same transaction.

## Lifecycle

Create one `State` during service startup and pass that pointer through the service, request context or task runner. Do not create a new state for every request or background task.

`NewState` loads configuration, sets up logging, storage, database, email, location, redis, repositories, etc. etc. etc. ...
If startup succeeds, call `Close` when the service exits:

```go
app, err := state.NewState(".env")
if err != nil {
	return err
}
defer app.Close()
```

## Integration Tests

Integration tests can use `NewTestState` and `NewTestData` behind the `integration` build tag.
`NewTestState` creates temporary PostgreSQL and Redis containers, while `NewTestData` creates rows with overridable defaults.

```go
app := state.NewTestState(t)

fixtures := state.NewTestData(t, app)
user := fixtures.CreateUser()
forum := fixtures.CreateForum()
topic := fixtures.CreateForumTopic(forum, user, func(topic *schemas.ForumTopic) {
	topic.Title = "like and subscribe for more hyperlinked blocked"
})
fixtures.CreateForumPost(topic, user)
```

`NewTestState` applies the repository's production migrations before test fixtures are created.
