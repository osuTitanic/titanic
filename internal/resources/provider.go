package resources

import (
	"context"
	"io"
)

// BeatmapResourceProvider provides access to beatmap resources
type BeatmapResourceProvider interface {
	// Setup prepares the provider for use.
	Setup(ctx context.Context) error

	// Osz returns a stream to the osz archive for the given beatmapset.
	// The caller is responsible for closing the returned stream.
	Osz(ctx context.Context, setId int, noVideo bool) (io.ReadCloser, int64, error)

	// Osu returns a stream to a single beatmap file.
	// The caller is responsible for closing the returned stream.
	Osu(ctx context.Context, beatmapId int) (io.ReadCloser, error)

	// Preview returns a stream to the audio preview for the given beatmapset.
	// The caller is responsible for closing the returned stream.
	Preview(ctx context.Context, setId int) (io.ReadCloser, error)

	// Background returns a stream to the thumbnail for the given beatmapset.
	// The caller is responsible for closing the returned stream.
	Background(ctx context.Context, setId int, large bool) (io.ReadCloser, error)
}
