# Storage

This module provides a storage api for files used by all services. The current implementation is a local filesystem backend rooted at `config.Config.DataPath`. S3 support will hopefully come too at some point.

## Usage with state system

Services should use `state.NewState(...)` instead of creating storage directly.
`State` creates the storage backend, calls `Setup()` & exposes it through `app.Storage`.

```go
app, err := state.NewState()
if err != nil {
	return err
}
defer app.Close()

if err := app.Storage.Save("avatar.jpg", "avatars", data); err != nil {
	return err
}
```

## Usage without state system

Create a local file storage backend from a data path & call `Setup()` before using it.

```go
store := storage.NewFileStorage(".data")

if err := store.Setup(); err != nil {
	return err
}

if err := store.Save("avatar.jpg", "avatars", data); err != nil {
	return err
}
```

## Streams

Use `SaveStream` and `ReadStream` for larger files.

```go
if err := store.SaveStream("replay.osr", "replays", reader); err != nil {
	return err
}

stream, err := store.ReadStream("replay.osr", "replays")
if err != nil {
	return err
}
defer stream.Close()
```

## Folders

Storage methods take a folder/directory name and a `key`. The local backend stores files as `<DATA_PATH>/<folder>/<key>`.
`Setup()` creates the required storage directories if they do not already exist.

Additionally, during setup, the storage downloads default assets declared in `DefaultAssetUrls` when they do not already exist.
At the moment this includes the default avatars in the `avatars` folder.
