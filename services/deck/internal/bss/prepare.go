package bss

import (
	"fmt"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
)

// PrepareBeatmapset locks & reloads an existing beatmapset and all of its difficulties.
func PrepareBeatmapset(submission *SubmissionContext, transaction *state.Repositories, beatmapsetId int) error {
	beatmapset, err := transaction.Beatmapsets.ByIdWithLock(beatmapsetId)
	if err != nil {
		return fmt.Errorf("bss: prepare beatmapset %d: %w", beatmapsetId, err)
	}
	if beatmapset == nil {
		return fmt.Errorf("bss: beatmapset not found: %d", beatmapsetId)
	}

	difficulties, err := transaction.Beatmaps.ManyBySetIdWithLock(beatmapset.Id)
	if err != nil {
		return fmt.Errorf("bss: prepare beatmaps for set %d: %w", beatmapset.Id, err)
	}

	beatmapset.Beatmaps = difficulties
	submission.Beatmapset = beatmapset
	return nil
}

// PrepareNewBeatmapset creates an inactive beatmapset owned by the authenticated user.
func PrepareNewBeatmapset(submission *SubmissionContext, transaction *state.Repositories) error {
	beatmapset := &schemas.Beatmapset{
		Creator:        &submission.User.Name,
		DisplayTitle:   new(""),
		Tags:           new(""),
		Status:         constants.BeatmapStatusInactive,
		Server:         constants.BeatmapServerTitanic,
		DownloadServer: constants.BeatmapServerTitanic,
		CreatorId:      &submission.User.Id,
		Available:      true,
		Beatmaps:       make([]*schemas.Beatmap, 0),
	}
	if err := transaction.Beatmapsets.Create(beatmapset); err != nil {
		return fmt.Errorf("bss: create beatmapset: %w", err)
	}
	submission.Beatmapset = beatmapset
	return nil
}

// PrepareBeatmaps prepares the requested beatmap IDs in request order.
// This keeps every assigned ID at the same position as its submitted difficulty, which is
// important because the client expects the order to be the same when processing the repsonse.
func PrepareBeatmaps(submission *SubmissionContext, transaction *state.Repositories, requestedIds []int) error {
	for _, requestedId := range requestedIds {
		if _, err := PrepareBeatmap(submission, transaction, requestedId); err != nil {
			return err
		}
	}
	return nil
}

// PrepareBeatmap assigns the requested ID a matching beatmap, i.e.
// a map that is already part of the set or a placeholder map for new beatmaps.
func PrepareBeatmap(submission *SubmissionContext, transaction *state.Repositories, requestedId int) (*PreparedBeatmap, error) {
	target := reusableBeatmap(submission, requestedId)
	if target == nil {
		if submission.Access == nil || !submission.Access.Owner {
			return nil, ErrBeatmapCreationNotAllowed
		}

		target = &schemas.Beatmap{
			SetId:  submission.Beatmapset.Id,
			Status: constants.BeatmapStatusInactive,
		}
		if err := transaction.Beatmaps.Create(target); err != nil {
			return nil, fmt.Errorf("bss: create beatmap for set %d: %w", submission.Beatmapset.Id, err)
		}
	}

	beatmap := &PreparedBeatmap{Target: target}
	submission.Beatmaps = append(submission.Beatmaps, beatmap)
	return beatmap, nil
}

// RemovedBeatmaps returns existing difficulties which were not prepared for this submission.
// They will eventually be removed from the set when the submission is processed.
func RemovedBeatmaps(submission *SubmissionContext) []*schemas.Beatmap {
	prepared := make(map[int]struct{}, len(submission.Beatmaps))
	for _, beatmap := range submission.Beatmaps {
		prepared[beatmap.Target.Id] = struct{}{}
	}

	removed := make([]*schemas.Beatmap, 0)
	for _, beatmap := range submission.Beatmapset.Beatmaps {
		if _, ok := prepared[beatmap.Id]; !ok {
			removed = append(removed, beatmap)
		}
	}
	return removed
}

func reusableBeatmap(submission *SubmissionContext, requestedId int) *schemas.Beatmap {
	if requestedId <= 0 {
		// new beatmap
		return nil
	}

	for _, beatmap := range submission.Beatmapset.Beatmaps {
		if beatmap.Id != requestedId {
			continue
		}
		// we found a match, check if it's already prepared
		for _, prepared := range submission.Beatmaps {
			if prepared.Target.Id == requestedId {
				return nil
			}
		}
		return beatmap
	}

	// no match found
	return nil
}
