# Media

This module provides helpers for processing media files: resizing images, creating thumbnails & extracting audio snippets. These are especially useful for the beatmap submission system.

## Images

Use `ResizeImage` to scale an image down to a given resolution. It accepts JPEG or PNG data, picks a scaler suited for downscaling & returns the encoded result in the original format.

```go
var small bytes.Buffer
if err := media.ResizeImage(source, &small, 80, 60); err != nil {
	return err
}
```

When the encoded image is already available as a byte slice, use the convenience methods:

```go
small, err := media.ResizeImageBytes(data, 80, 60)
png, err := media.ResizeImageBytesToPng(data, 256, 256)
```

Use `ResizeAndCropImage` when the image should keep its aspect ratio while filling the image. It crops excess width / height from the center before resizing:

```go
var thumbnail bytes.Buffer
if err := media.ResizeAndCropImage(source, &thumbnail, 160, 120); err != nil {
	return err
}
```

For more control, configure and use `ImageGenerator` directly:

```go
generator := &media.ImageGenerator{
	Width:     160,
	Height:    120,
	Scaler:    "CatmullRom",
	CropToFit: true,
}

if err := generator.CreateThumbnail(source, target, "png"); err != nil {
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
