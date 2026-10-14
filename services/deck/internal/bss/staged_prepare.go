package bss

import (
	"errors"
	"fmt"

	"github.com/Lekuruu/gosu/pkg/beatmaps"
	"github.com/osuTitanic/titanic/internal/performance"
	"github.com/osuTitanic/titanic/internal/state"
)

// PrepareStagedBeatmaps parses & assignes beatmaps from a pre-osz2 submission staged
// in Redis, and derives the set metadata from one of the provided difficulties.
func (submission *SubmissionContext) PrepareStagedBeatmaps(
	transaction *state.Repositories,
	staged *StagedSubmission,
	ppv1 *performance.PPv1Service,
	ppv2 performance.IPPv2Service,
) error {
	if submission.IsCanceled() {
		return submission.Context.Err()
	}
	if submission.User == nil {
		return ErrUserNotAuthenticated
	}
	if submission.Beatmapset == nil {
		return ErrBeatmapsetNotPrepared
	}
	if submission.Access == nil {
		return errors.New("bss: submission access not resolved")
	}

	if staged.UserId != submission.User.Id {
		return fmt.Errorf("bss: staged submission user %d does not match user %d", staged.UserId, submission.User.Id)
	}
	if staged.BeatmapsetId != submission.Beatmapset.Id {
		return fmt.Errorf(
			"bss: staged beatmapset %d does not match prepared beatmapset %d",
			staged.BeatmapsetId,
			submission.Beatmapset.Id,
		)
	}

	if len(staged.Beatmaps) == 0 {
		return errors.New("bss: staged submission has no beatmaps")
	}
	if len(submission.Beatmaps) != 0 {
		return errors.New("bss: beatmaps already prepared")
	}

	requestedIds := make([]int, len(staged.Beatmaps))
	for i, beatmap := range staged.Beatmaps {
		if err := beatmap.Validate(); err != nil {
			return fmt.Errorf("bss: validate staged beatmap %d: %w", i, err)
		}
		requestedIds[i] = beatmap.Id
	}

	// Populate the submission beatmaps w/ existing beatmaps or new beatmaps
	if err := submission.PrepareBeatmaps(transaction, requestedIds); err != nil {
		return fmt.Errorf("bss: prepare staged beatmaps: %w", err)
	}

	var metadataSource *beatmaps.Beatmap

	for i, beatmap := range staged.Beatmaps {
		if submission.IsCanceled() {
			return submission.Context.Err()
		}

		prepared := submission.Beatmaps[i]
		if prepared.Target == nil {
			return fmt.Errorf("bss: staged beatmap %d has no prepared target", i)
		}
		if !submission.Access.CanUpdateBeatmap(prepared.Target.Id) {
			continue
		}

		prepared, err := submission.AssignBeatmapContent(
			prepared.Target.Id,
			beatmap.Filename,
			beatmap.Contents,
		)
		if err != nil {
			return fmt.Errorf("bss: assign staged beatmap %q: %w", beatmap.Filename, err)
		}
		if err := prepared.CalculateDifficulty(ppv1, ppv2); err != nil {
			return fmt.Errorf("bss: calculate staged beatmap %q difficulty: %w", beatmap.Filename, err)
		}
		metadataSource = prepared.Source
	}

	if metadataSource == nil {
		return errors.New("bss: staged submission has no metadata source")
	}
	submission.Metadata = MetadataFromBeatmap(metadataSource)
	submission.Beatmapset.HasVideo = staged.HasVideo
	submission.Beatmapset.HasStoryboard = staged.HasStoryboard || submission.BeatmapsHaveSB()
	return nil
}
