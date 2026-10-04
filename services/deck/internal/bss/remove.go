package bss

import (
	"errors"
	"fmt"

	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
)

var ErrBeatmapRemovalNotAllowed = errors.New("bss: beatmap removal not allowed")

// RemoveBeatmaps deletes difficulties and their dependent records.
func (submission *SubmissionContext) RemoveBeatmaps(transaction *state.Repositories, beatmaps []*schemas.Beatmap) error {
	if err := submission.CheckBeatmapRemovals(beatmaps); err != nil {
		return err
	}

	for _, beatmap := range beatmaps {
		if err := transaction.Plays.DeleteByBeatmapId(beatmap.Id); err != nil {
			return fmt.Errorf("bss: delete plays for beatmap %d: %w", beatmap.Id, err)
		}
		if err := transaction.Beatmaps.Delete(beatmap); err != nil {
			return fmt.Errorf("bss: delete beatmap %d: %w", beatmap.Id, err)
		}
	}
	return nil
}

// CheckBeatmapRemovals ensures that the user is allowed to remove the given difficulties.
// Collaborators should not be able to delete difficulties from a beatmapset.
func (submission *SubmissionContext) CheckBeatmapRemovals(beatmaps []*schemas.Beatmap) error {
	if len(beatmaps) == 0 {
		return nil
	}
	if submission.Access == nil || !submission.Access.Owner {
		return ErrBeatmapRemovalNotAllowed
	}
	for _, beatmap := range beatmaps {
		if beatmap.SetId != submission.Beatmapset.Id {
			return ErrBeatmapRemovalNotAllowed
		}
	}
	return nil
}
