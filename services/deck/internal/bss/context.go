package bss

import (
	"context"
	"io/fs"

	"github.com/Lekuruu/gosu/pkg/beatmaps"
	"github.com/Lekuruu/osz2-go/pkg/osz2"
	"github.com/osuTitanic/titanic/internal/schemas"
)

// SubmissionContext contains the state shared by
// beatmap submission operations during one request.
type SubmissionContext struct {
	Context context.Context

	User       *schemas.User
	Beatmapset *schemas.Beatmapset
	Beatmaps   []*PreparedBeatmap

	// Metadata is the set metadata sourced from an
	// osz2 package or from one of the beatmap files
	Metadata osz2.Metadata

	// FS is the provided read-only beatmap package (osz or osz2)
	FS fs.FS
}

type PreparedBeatmap struct {
	Source   *beatmaps.Beatmap
	Target   *schemas.Beatmap
	Metadata BeatmapMetadata
}

func NewSubmissionContext(ctx context.Context) *SubmissionContext {
	return &SubmissionContext{
		Context:  ctx,
		Beatmaps: make([]*PreparedBeatmap, 0),
	}
}
