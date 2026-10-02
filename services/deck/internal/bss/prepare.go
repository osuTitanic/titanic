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
