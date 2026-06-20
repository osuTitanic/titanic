package resources

import "io"

// BeatmapResourceResolver receives the actual resource data from a service
// (e.g. external mirrors or our own local storage).
type BeatmapResourceResolver interface {
	// Setup prepares the resolver for use.
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

// TODO: we could simplify this by consolidating BeatmapResourceResolver and BeatmapResourceProvider
// 		 but for now we'll keep them separate i think
