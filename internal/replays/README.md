# Replays

This module contains helpers for doing stuff with osu! replays.
It can serialize replay data into `.osr` files, and deserialize compressed replay frames.

## Usage

### Serializing (osr files)

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

### Deserializing (frames)

Use `DeserializeFrames` to decompress a stored replay payload & parse its individual frames.

```go
// Example: read compressed replay frames from storage
data, err := app.Storage.Read(strconv.FormatInt(score.Id, 10), "replays")
if err != nil {
	return err
}

frames, seed, err := replays.DeserializeFrames(data)
if err != nil {
	return err
}

// do stuff ...
```

Each `Frame` contains its delta time, absolute replay time, cursor coordinates & pressed buttons.
Use `ButtonState.Has` to check individual buttons:

```go
if frame.Buttons.Has(replays.Left1) {
	// Left mouse button is pressed
}
```

The returned seed is read from the replay's optional RNG seed frame & is `0` when no matching frame is present. The seed is used for the osu!mania random mod.

## Planned Features

Touchscreen detection is one of the planned feature to add in here. The old deck implementation currently lives in [here](https://github.com/osuTitanic/deck/blob/main/app/helpers/replays.py#L97). This would then also include validation of replay payloads.
