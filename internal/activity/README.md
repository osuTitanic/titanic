# Activity

This module contains helpers for publishing and rendering user activity. Activity entries are used for profile activity timelines, live bancho events, and (eventually) other formats such as `#announce` or discord embeds.

## Submission

Use `Submit` when a service wants to announce something a user did. It publishes the activity to the redis event channel and, unless hidden, stores it in `profile_activity`.

```go
err := activity.Submit(
	app,
	user.Id,
	&mode,
	constants.ActivityBeatmapLeaderboardRank,
	map[string]any{
		"username":     user.Name,
		"beatmap_id":   beatmap.Id,
		"beatmap":      beatmap.Name(),
		"beatmap_rank": rank,
		"mode":         mode.Short(),
		"mods":         mods.String(),
		"pp":           pp,
	},
	true,
	false,
)
if err != nil {
	return err
}
```

The `mode` argument can be `nil` for activity that is not tied to a specific game mode.
`isAnnouncement` determines whether the activity should be broadcast to #announce and discord.
`isHidden` means the event is only broadcast and will not be stored in the database, e.g. for user logins.

## Rendering

### HTML

Use `RenderHtml` to render a stored activity entry for the website.

```go
html := activity.RenderHtml(entry)
if html == "" {
	// Unknown activity type or invalid activity data
}
```

Stern exposes this to templates through the `formatActivity` template filter, then renders it in the user profile activity table.

### Other

`RenderText` and `RenderDiscord` exist as placeholders, but they are not implemented yet.

<!-- TODO: Text & Discord Renderers -->
