package bss

import (
	"errors"
	"fmt"
	"strings"

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

	// Keep the resolved beatmap IDs in the same order as the uploaded files
	requestedIds := make([]int, len(staged.Beatmaps))

	// We want to keep difficulty names, filenames, and checksums unique across all staged beatmaps
	usedIds := make(map[int]bool, len(staged.Beatmaps))
	versions := make(map[string]bool, len(staged.Beatmaps))
	filenames := make(map[string]bool, len(staged.Beatmaps))
	checksums := make(map[string]bool, len(staged.Beatmaps))

	for i := range staged.Beatmaps {
		beatmap := &staged.Beatmaps[i]
		if err := beatmap.Validate(); err != nil {
			return fmt.Errorf("bss: validate staged beatmap %d: %w", i, err)
		}
		source, err := beatmap.Parse()
		if err != nil {
			return fmt.Errorf("bss: parse staged beatmap %q: %w", beatmap.Filename, err)
		}

		filename := strings.ToLower(beatmap.Filename)
		checksum := strings.ToLower(beatmap.Checksum)
		version := source.Version
		if version == "" {
			version = "Normal"
		}

		// Reject repeated difficulties or files within the same submission
		if versions[version] || filenames[filename] || checksums[checksum] {
			return ErrDuplicateBeatmap
		}
		versions[version] = true
		filenames[filename] = true
		checksums[checksum] = true

		// Find an existing difficulty match by name
		// If there is no match, `PrepareBeatmaps` will create a new difficulty
		beatmap.Id = -1

		for _, existing := range submission.Beatmapset.Beatmaps {
			if existing.Version != version || usedIds[existing.Id] {
				// Difficulty doesn't match or the ID has already been used
				continue
			}
			beatmap.Id = existing.Id
			usedIds[existing.Id] = true
			break
		}
		requestedIds[i] = beatmap.Id
	}

	// Prepare an existing / new beatmap from db for every uploaded difficulty
	if err := submission.PrepareBeatmaps(transaction, requestedIds); err != nil {
		return fmt.Errorf("bss: prepare staged beatmaps: %w", err)
	}
	for i, requestedId := range requestedIds {
		if requestedId <= 0 {
			submission.Beatmaps[i].Target.Status = submission.Beatmapset.Status
		}
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

	// Derive set-wide metadata from one of the beatmaps
	submission.Metadata = MetadataFromBeatmap(metadataSource)
	submission.Beatmapset.HasVideo = staged.HasVideo
	submission.Beatmapset.HasStoryboard = staged.HasStoryboard || submission.BeatmapsHaveSB()
	return nil
}
