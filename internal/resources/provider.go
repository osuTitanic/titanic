package resources

import "io"

// BeatmapResourceProvider provides access to beatmap resources
type BeatmapResourceProvider interface {
	// Setup prepares the provider & all of its resolvers
	Setup() error

	// Osz returns a stream to the osz archive for the given beatmapset.
	// The caller is responsible for closing the returned stream.
	Osz(setId int, noVideo bool) (io.ReadCloser, error)

	// Osu returns a stream to a single beatmap file.
	// The caller is responsible for closing the returned stream.
	Osu(beatmapId int) (io.ReadCloser, error)

	// Preview returns a stream to the audio preview for the given beatmapset.
	// The caller is responsible for closing the returned stream.
	Preview(setId int) (io.ReadCloser, error)

	// Background returns a stream to the thumbnail for the given beatmapset.
	// The caller is responsible for closing the returned stream.
	Background(setId int, large bool) (io.ReadCloser, error)
}
