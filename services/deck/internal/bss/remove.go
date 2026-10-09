package bss

import (
	"fmt"

	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
)

// RemoveInactiveBeatmapsets removes all inactive (placeholder or deleted) beatmapsets owned by the user.
// The "placeholder" state means that a set was reserved by the user but was never populated with metadata.
func (submission *SubmissionContext) RemoveInactiveBeatmapsets(transaction *state.Repositories) ([]*schemas.Beatmapset, error) {
	beatmapsets, err := transaction.Beatmapsets.FetchInactiveByCreatorWithLock(submission.User.Id)
	if err != nil {
		return nil, fmt.Errorf("bss: fetch inactive beatmapsets for user %d: %w", submission.User.Id, err)
	}

	for _, beatmapset := range beatmapsets {
		beatmaps, err := transaction.Beatmaps.ManyBySetIdWithLock(beatmapset.Id)
		if err != nil {
			return nil, fmt.Errorf("bss: fetch inactive beatmaps for set %d: %w", beatmapset.Id, err)
		}
		beatmapset.Beatmaps = beatmaps

		if err := transaction.Modding.DeleteBySetId(beatmapset.Id); err != nil {
			return nil, fmt.Errorf("bss: delete modding for inactive set %d: %w", beatmapset.Id, err)
		}
		if err := transaction.Ratings.DeleteBySetId(beatmapset.Id); err != nil {
			return nil, fmt.Errorf("bss: delete ratings for inactive set %d: %w", beatmapset.Id, err)
		}
		if err := transaction.Plays.DeleteBySetId(beatmapset.Id); err != nil {
			return nil, fmt.Errorf("bss: delete plays for inactive set %d: %w", beatmapset.Id, err)
		}
		if err := transaction.Nominations.DeleteAll(beatmapset.Id); err != nil {
			return nil, fmt.Errorf("bss: delete nominations for inactive set %d: %w", beatmapset.Id, err)
		}
		if err := transaction.Favourites.DeleteBySetId(beatmapset.Id); err != nil {
			return nil, fmt.Errorf("bss: delete favourites for inactive set %d: %w", beatmapset.Id, err)
		}

		for _, beatmap := range beatmaps {
			if err := transaction.Beatmaps.Delete(beatmap); err != nil {
				return nil, fmt.Errorf("bss: delete inactive beatmap %d: %w", beatmap.Id, err)
			}
		}
		if err := transaction.Beatmapsets.Delete(beatmapset); err != nil {
			return nil, fmt.Errorf("bss: delete inactive beatmapset %d: %w", beatmapset.Id, err)
		}
	}

	// TODO: Should this handle storage too?
	return beatmapsets, nil
}

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
