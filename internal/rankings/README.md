# Rankings

This module maintains all kinds of player rankings through redis sorted sets. It provides methods to fetch ranks & scores, update rankings when stats change, and remove players from rankings when they become inactive or restricted.

Services should usually access rankings through the state system (`state.Rankings`). A manual setup is only needed for tests, or other tooling that already has a redis client.

## Usage with state system

`state.NewState(...)` creates the redis client & exposes the rankings service through `app.Rankings`.

```go
app, err := state.NewState()
if err != nil {
	return err
}
defer app.Close()

rank, err := app.Rankings.GlobalRank(user.Id, constants.ModeOsu)
if err != nil {
	return err
}
```

Fetch a leaderboard page by rank type:

```go
players, err := app.Rankings.TopPlayers(
	constants.ModeOsu,
	0,
	50,
	"performance",
	nil,
)
if err != nil {
	return err
}
```

Provide a country code for country-specific rankings:

```go
country := "de"
rank, err := app.Rankings.Rank(user.Id, constants.ModeOsu, "performance", &country)
if err != nil {
	return err
}
```

## Usage without state system

Create a service from an existing Redis client:

```go
service := rankings.NewRankingsService(redisClient)

score, err := service.Performance(user.Id, constants.ModeOsu)
if err != nil {
	return err
}
```

## Indexing

Use `Update` whenever stats for a user change. Right now, it updates global and country sorted sets for `performance`, `rscore`, `tscore`, `ppv1`, `acc`, and `clears`.

```go
if err := app.Rankings.Update(stats, user.Country); err != nil {
	return err
}
```

Leader scores and kudosu require repository providers & are updated separately to avoid circular dependencies:

```go
if err := app.Rankings.UpdateLeaderScores(stats, user.Country, app.Repositories.Scores); err != nil {
	return err
}

if err := app.Rankings.UpdateKudosu(user.Id, user.Country, app.Repositories.Modding); err != nil {
	return err
}
```

Remove users from rankings when they become inactive, restricted, or otherwise unlisted:

```go
if err := app.Rankings.Remove(user.Id, user.Country); err != nil {
	return err
}
```

Use `RemoveFromCountry` when only the country-specific entries need to be cleared, e.g. for a country change.

## Helpers

Use `PlayerCount` and `TopPlayers` for leaderboard pagination. These only include entries with scores greater than zero.

Use `TopCountries` to fetch country rankings for a mode. It sums country-wide `performance`, `rscore`, and `tscore`, then sorts countries by total performance.

Use `PlayerAbove` to find the user directly above another player in a rank type & the score difference needed to reach them. It returns `ErrNoPlayerAbove` when the player is missing or already first.
