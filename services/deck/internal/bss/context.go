package bss

import (
	"context"
	"io/fs"

	"github.com/Lekuruu/gosu/pkg/beatmaps"
	"github.com/Lekuruu/osz2-go/pkg/osz2"
	"github.com/osuTitanic/titanic/internal/schemas"
)

// PreparedBeatmap holds a beatmap source and target,
// used for beatmap submission operations.
type PreparedBeatmap struct {
	Source *beatmaps.Beatmap
	Target *schemas.Beatmap
}

// SubmissionContext contains the state shared by
// beatmap submission operations during one request.
type SubmissionContext struct {
	Context context.Context

	User       *schemas.User
	Beatmapset *schemas.Beatmapset
	Beatmaps   []*PreparedBeatmap
	Access     *SubmissionAccess

	// Metadata is the set metadata sourced from an
	// osz2 package or from one of the beatmap files
	Metadata osz2.Metadata

	// FS is the provided read-only beatmap package (osz or osz2)
	FS fs.FS
}

func NewSubmissionContext(ctx context.Context) *SubmissionContext {
	return &SubmissionContext{
		Context:  ctx,
		Beatmaps: make([]*PreparedBeatmap, 0),
	}
}

func (s *SubmissionContext) IsCanceled() bool {
	return s != nil &&
		s.Context != nil &&
		s.Context.Err() != nil
}

func (s *SubmissionContext) BeatmapById(beatmapId int) *PreparedBeatmap {
	for _, beatmap := range s.Beatmaps {
		if beatmap.Target.Id == beatmapId {
			return beatmap
		}
	}
	return nil
}
