# Bancho

This module provides access to bancho user statuses & events.
`state.NewState(...)` initializes both services and exposes them as `BanchoUsers` & `BanchoEvents`.

## User statuses

Use `BanchoUsers.Get` to read the current status of an online user:

```go
status, err := app.BanchoUsers.Get(ctx, userId)
if err != nil {
	return err
}
if status == nil {
	// The user is offline / has no status entry
}
```

Use `Exists` when only the user's online presence is needed:

```go
online, err := app.BanchoUsers.Exists(ctx, userId)
```

## Events

Use the helper functions to publish common events:

```go
err := app.BanchoEvents.UserUpdate(
	ctx,
	userId,
	mode,
)
```

```go
err := app.BanchoEvents.Activity(
	ctx,
	userId,
	&mode,
	constants.ActivityAchievementUnlocked,
	map[string]any{
		"achievement": achievement.Name,
	},
	true,
)
```

Events are encoded in json, published to the `bancho:events` redis channel, and processed by the bancho service.
Use `Dispatch` directly for an event that does not yet have a dedicated helper.
