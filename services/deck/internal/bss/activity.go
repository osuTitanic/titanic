package bss

import (
	"encoding/json"
	"fmt"

	"github.com/osuTitanic/titanic/internal/activity"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
)

// BroadcastUploadActivity publishes & stores the activity for a beatmap submission.
func BroadcastActivity(isUpdate bool, app *state.State, submission *SubmissionContext) error {
	if isUpdate {
		return BroadcastUploadActivity(app, submission)
	} else {
		return BroadcastUpdateActivity(app, submission)
	}
}

// BroadcastUploadActivity publishes & stores the activity for a newly submitted beatmapset.
func BroadcastUploadActivity(app *state.State, submission *SubmissionContext) error {
	beatmapset, creatorId, mode, err := submissionActivityContext(submission)
	if err != nil {
		return err
	}

	return activity.Submit(
		app,
		creatorId,
		&mode,
		constants.ActivityBeatmapUploaded,
		map[string]any{
			"username":        submission.User.Name,
			"title":           beatmapset.TitleName(),
			"artist":          beatmapset.ArtistName(),
			"beatmapset_id":   beatmapset.Id,
			"beatmapset_name": fmt.Sprintf("%s (%s)", beatmapset.Name(), beatmapset.CreatorName()),
		},
		true,  // should be sent to #announce
		false, // should be visible on user profile
	)
}

// BroadcastUpdateActivity publishes the activity for an updated beatmapset.
// Consecutive submissions are not stored on the profile activity to reduce spam.
func BroadcastUpdateActivity(app *state.State, submission *SubmissionContext) error {
	beatmapset, userId, mode, err := submissionActivityContext(submission)
	if err != nil {
		return err
	}

	latest, err := app.Activities.FetchLatestByUser(userId)
	if err != nil {
		return fmt.Errorf("bss: fetch latest submission activity: %w", err)
	}

	activityIsDuplicate := isDuplicateUpdateActivity(latest, beatmapset.Id)
	activityName := fmt.Sprintf("%s (%s)", beatmapset.Name(), beatmapset.CreatorName())

	return activity.Submit(
		app,
		userId,
		&mode,
		constants.ActivityBeatmapUpdated,
		map[string]any{
			"username":        submission.User.Name,
			"beatmapset_id":   beatmapset.Id,
			"beatmapset_name": activityName,
		},
		false,               // should not be sent to #announce
		activityIsDuplicate, // should only be seen on user profile if not a duplicate
	)
}

func submissionActivityContext(submission *SubmissionContext) (*schemas.Beatmapset, int, constants.Mode, error) {
	// Determine the primary mode for the activity db entry
	var primaryMode constants.Mode
	var primaryCount = 0

	counts := make(map[constants.Mode]int, len(constants.Modes))
	for _, beatmap := range submission.Beatmaps {
		counts[beatmap.Target.Mode]++
	}
	for _, beatmap := range submission.Beatmaps {
		if count := counts[beatmap.Target.Mode]; count > primaryCount {
			primaryMode = beatmap.Target.Mode
			primaryCount = count
		}
	}

	return submission.Beatmapset, submission.User.Id, primaryMode, nil
}

func isDuplicateUpdateActivity(latest *schemas.Activity, setId int) bool {
	if latest == nil {
		return false
	}
	if latest.Type != int(constants.ActivityBeatmapUploaded) && latest.Type != int(constants.ActivityBeatmapUpdated) {
		return false
	}

	var data struct {
		BeatmapsetId int `json:"beatmapset_id"`
	}
	if err := json.Unmarshal(latest.Data, &data); err != nil {
		return false
	}
	return data.BeatmapsetId == setId
}
