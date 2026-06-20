# Media

This module provides helpers for processing media files: resizing images, creating thumbnails & extracting audio snippets. These are especially useful for the beatmap submission system.

## Images

Use `ResizeImage` to scale an image down to a given resolution. It accepts JPEG or PNG data, picks a scaler suited for downscaling & returns the encoded result in the original format.

```go
small, err := media.ResizeImage(data, 80, 60)
if err != nil {
	return err
}
```

For more control, use the `ImageGenerator` directly. Build an `Image` from a byte array or a file, then create the thumbnail:

```go
generator := &media.ImageGenerator{
	Width:  160,
	Height: 120,
	Scaler: "CatmullRom",
}

image, err := generator.NewImageFromByteArray(data)
if err != nil {
	return err
}

thumbnail, err := generator.CreateThumbnail(image)
if err != nil {
	return err
}
```

`Scaler` selects the interpolation algorithm: `NearestNeighbor`, `ApproxBiLinear`, `BiLinear` or `CatmullRom`. `CatmullRom` gives the best results when downscaling.

The image code is adapted from Matthew Jorgensen's [go-thumbnail](https://github.com/prplecake/go-thumbnail) (BSD 3-Clause).

## Audio

Use `ExtractAudioSnippet` to cut a snippet out of an audio file & re-encode it as MP3. It takes the raw audio bytes, an offset & duration (in seconds) and a bitrate, returning the encoded snippet.

```go
snippet, err := media.ExtractAudioSnippet(audioData, 0, 10, 128000)
if err != nil {
	return err
}
```

This relies on `ffmpeg` being available, through the [ffmpeg-go](https://github.com/Lekuruu/ffmpeg-go) bindings. Glory to the allmighty ffmpeg.
