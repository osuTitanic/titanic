package bss

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"

	"github.com/Lekuruu/gosu/pkg/beatmaps"
	"github.com/osuTitanic/titanic/internal/media"
)

type ResourceSelector func(beatmap *beatmaps.Beatmap) (string, bool)

const (
	thumbnailWidth  = 160
	thumbnailHeight = 120

	audioPreviewDuration = 10
	audioPreviewBitrate  = 128_000
)

// GenerateThumbnail creates the beatmapset thumbnail
// from the first background referenced by a difficulty.
func (submission *SubmissionContext) GenerateThumbnail(source fs.FS) ([]byte, error) {
	// Return the first valid background we encounter on a beatmap
	pickBackground := func(beatmap *beatmaps.Beatmap) (string, bool) {
		backgrounds := beatmap.Events.Backgrounds()
		if len(backgrounds) == 0 {
			return "", false
		}
		if backgrounds[0].Filename == "" {
			return "", false
		}
		return backgrounds[0].Filename, true
	}

	resource, _, err := submission.findBeatmapResource(source, pickBackground)
	if err != nil || resource == nil {
		return nil, err
	}
	defer resource.Close()

	// Thumbnails should fit comfortably within memory
	var thumbnail bytes.Buffer

	err = media.ResizeAndCropImage(
		resource, &thumbnail,
		thumbnailWidth, thumbnailHeight,
	)
	if err != nil {
		return nil, fmt.Errorf("bss: resize beatmap thumbnail: %w", err)
	}

	return thumbnail.Bytes(), nil
}

// GenerateAudioPreview creates an audio preview from
// the first audio file referenced by a difficulty.
func (submission *SubmissionContext) GenerateAudioPreview(source fs.FS) ([]byte, error) {
	// Pick the first valid audio file
	pickAudio := func(beatmap *beatmaps.Beatmap) (string, bool) {
		return beatmap.Audio, beatmap.Audio != ""
	}

	resource, beatmap, err := submission.findBeatmapResource(source, pickAudio)
	if err != nil || resource == nil {
		return nil, err
	}
	defer resource.Close()

	// Beatmaps with an unspecified preview time
	// will receive one around 40% into the song
	previewTime := -1
	if beatmap.PreviewTime >= 0 {
		previewTime = int(beatmap.PreviewTime / 1000)
	}

	preview, err := media.ExtractAudioSnippetFromReader(
		resource,
		previewTime,
		audioPreviewDuration,
		audioPreviewBitrate,
	)
	if err != nil {
		return nil, fmt.Errorf("bss: extract beatmap audio preview: %w", err)
	}
	return preview, nil
}

func (submission *SubmissionContext) findBeatmapResource(source fs.FS, selectResource ResourceSelector) (fs.File, *beatmaps.Beatmap, error) {
	if source == nil {
		return nil, nil, errors.New("bss: filesystem not set")
	}

	for _, prepared := range submission.Beatmaps {
		filename, selected := selectResource(prepared.Source)
		if !selected {
			continue
		}

		resource, err := readFile(
			source, filename,
			MaxPackageFileSize,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("bss: read submission resource %q: %w", filename, err)
		}
		return resource, prepared.Source, nil
	}

	return nil, nil, nil
}
