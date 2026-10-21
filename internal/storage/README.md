# Storage

This module provides a common file-storage API for Titanic services. It includes two backends:

- `FileStorage` stores files below `config.Config.DataPath`
- `S3Storage` stores objects in an S3-compatible bucket through MinIO's client library

Both implement the `Storage` interface.

## Usage with state system

Services should use `state.NewState(...)` instead of creating storage directly.
`State` selects the backend from the configuration, calls `Setup`, and exposes it through `app.Storage`.

```go
app, err := state.NewState()
if err != nil {
	return err
}
defer app.Close()

ctx := context.Background()
if err := app.Storage.Save(ctx, "avatar.jpg", "avatars", data); err != nil {
	return err
}
```

## Usage without state system

Create a local file storage backend from a data path & call `Setup(...)` before using it.

```go
store := storage.NewFileStorage(".data")
ctx := context.Background()

if err := store.Setup(ctx); err != nil {
	return err
}

if err := store.Save(ctx, "avatar.jpg", "avatars", data); err != nil {
	return err
}
```

Create an S3-compatible backend:

```go
store, err := storage.NewS3Storage(storage.S3Config{
	Endpoint:        "s3.example.com",
	AccessKeyID:     accessKey,
	SecretAccessKey: secretKey,
	BucketName:      "osutitanic",
	Region:          "eu-central-1",
	UseSSL:          true,
})
if err != nil {
	return err
}
if err := store.Setup(ctx); err != nil {
	return err
}
```

## Streams

Use `SaveStream` and `ReadStream` for larger files.

```go
if err := store.SaveStream(ctx, "replay.osr", "replays", reader); err != nil {
	return err
}

stream, err := store.ReadStream(ctx, "replay.osr", "replays")
if err != nil {
	return err
}
defer stream.Close()
```

## Folders

Storage methods take a `key` and then a folder name. The local backend stores files as `<DATA_PATH>/<folder>/<key>`.
`FileStorage.Setup` creates the required directories when they do not exist.

Additionally, during setup, the storage downloads default assets declared in `DefaultAssetUrls` when they do not already exist.
At the moment this includes the default avatars in the `avatars` folder.
