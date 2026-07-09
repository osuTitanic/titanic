# Replays

This module contains helpers for building osu! replay files.
It allows for serialization of stored replay data into playable `.osr` files and computing replay checksums.

## Usage

Fetch the score with the relationships required by replay serialization, read the stored replay payload, then pass both to `Serialize`.

```go
score, err := app.Repositories.Scores.ById(
	scoreId,
	"User",
	"Beatmap",
)
if err != nil {
	return err
}
if score == nil || score.User == nil || score.Beatmap == nil {
	return errors.New("score is missing metadata")
}

data, err := app.Storage.Read(strconv.FormatInt(score.Id, 10), "replays")
if err != nil {
	return err
}

replay := replays.Serialize(score, data)
```

`Serialize` writes the replay metadata in the format expected by osu! clients.
The replay payload should already be the compressed replay frame data stored for that score.

## Planned Features

Touchscreen detection is one of the planned feature to add in here. The old deck implementation currently lives in [here](https://github.com/osuTitanic/deck/blob/main/app/helpers/replays.py#L97). This would then also include validation of replay payloads.
