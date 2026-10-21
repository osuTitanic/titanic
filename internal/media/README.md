# Media

This module provides helpers for resizing images, creating thumbnails & extracting audio snippets. They are used primarily by the beatmap submission system.

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

Use `ResizeAndCropImage` when the image should keep its aspect ratio while filling the target dimensions. It crops excess width / height from the center before resizing:

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

Use `ExtractAudioSnippet` to cut a section from an audio file & encode it as MP3. It accepts the raw audio bytes, offset & duration (in seconds), and a bitrate. Pass a negative offset to choose a default position around 40% into the track.

```go
snippet, err := media.ExtractAudioSnippet(audioData, 0, 10, 128000)
if err != nil {
	return err
}
```

Use `ExtractAudioSnippetFromReader` when the source is already available as an `io.Reader`, or `ExtractAudioSnippetFromPath` for an existing file.

Audio extraction requires `ffmpeg` and `ffprobe` on the process path and uses the [ffmpeg-go](https://github.com/Lekuruu/ffmpeg-go) bindings.
