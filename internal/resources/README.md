# Resources

This module provides access to beatmap resources (`.osz` archives, `.osu` files, audio previews & beatmap thumbnails), abstracting away *where* a resource is actually fetched from.

## Provider & resolvers

Everything implements a single `BeatmapResourceProvider` interface, providing access to `Osz`, `Osu`, `Preview` & `Background`:

- `MirrorResolver` fetches resources from external mirrors over HTTP, performing requests in a round-robin fashion and checking for rate limits.
- `StorageResolver` reads resources directly from our own storage backend, and is used for beatmaps hosted on titanic.
- `BeatmapProvider` inspects a beatmap's configured download server and delegates the request to the matching per-server implementation, i.e. either `MirrorResolver` or `StorageResolver`.

If we would want to add new download servers, we can be implement a new `BeatmapResourceProvider` and register the implementation in `NewBeatmapProvider`.

## Usage with state system

Services should use `state.NewState(...)` instead of creating the provider directly. `State` creates the provider, calls `Setup()` & exposes it through `app.Resources`.

```go
app, err := state.NewState()
if err != nil {
	return err
}
defer app.Close()

stream, err := app.Resources.Osz(setId, false)
if err != nil {
	return err
}
defer stream.Close()

// do whatever
```
