package tasks

import (
	"log/slog"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
)

const (
	QualifiedToRankedTime  = time.Hour * 24 * 7
	PendingToGraveyardTime = time.Hour * 24 * 14
)

const (
	BeatmapForumIdRanked    = 8
	BeatmapForumIdPending   = 9
	BeatmapForumIdWIP       = 10
	BeatmapForumIdGraveyard = 12
)

// UpdateBeatmapStatuses updates beatmap statuses (graveyard & qualified)
func UpdateBeatmapStatuses(app *state.State, logger *slog.Logger) error {
	logger.Info("Updating qualified beatmaps...")
	qualifiedSets, err := app.Repositories.Beatmapsets.FetchByStatus(constants.BeatmapStatusQualified, "Beatmaps")
	if err != nil {
		return err
	}
	for _, beatmapset := range qualifiedSets {
		handleQualifiedSet(app, logger, beatmapset)
	}

	logger.Info("Updating pending beatmaps...")
	pendingSets, err := app.Repositories.Beatmapsets.FetchByStatus(constants.BeatmapStatusPending, "Beatmaps")
	if err != nil {
		return err
	}
	for _, beatmapset := range pendingSets {
		handlePendingSet(app, logger, beatmapset)
	}

	logger.Info("Updating WIP beatmaps...")
	wipSets, err := app.Repositories.Beatmapsets.FetchByStatus(constants.BeatmapStatusWIP, "Beatmaps")
	if err != nil {
		return err
	}
	for _, beatmapset := range wipSets {
		handlePendingSet(app, logger, beatmapset)
	}

	return nil
}

func handleQualifiedSet(app *state.State, logger *slog.Logger, beatmapset *schemas.Beatmapset) {
	if beatmapset.ApprovedAt == nil {
		return
	}
	approvedTime := time.Since(*beatmapset.ApprovedAt)
	rankingTime := QualifiedToRankedTime
	now := time.Now()

	if approvedTime < rankingTime {
		return
	}

	// why is this a thing actually... bancho should disable this too
	if app.Config.RemoveScoresOnRanked {
		hideScoresForSet(app, beatmapset)
	}

	maxDrain := 0
	for _, beatmap := range beatmapset.Beatmaps {
		if beatmap.DrainLength > maxDrain {
			maxDrain = beatmap.DrainLength
		}
	}

	// Determine status based on drain time
	// Map will be set to "Approved" if drain time exceeds 5 minutes, otherwise "Ranked"
	status := constants.BeatmapStatusRanked
	if maxDrain >= 5*60 {
		status = constants.BeatmapStatusApproved
	}

	updateBeatmapIcon(app, beatmapset, status, beatmapset.Status)
	moveBeatmapTopic(app, beatmapset, status)

	beatmapsetUpdate := &schemas.Beatmapset{
		Id:         beatmapset.Id,
		Status:     status,
		ApprovedAt: &now,
	}
	beatmapUpdate := &schemas.Beatmap{
		SetId:  beatmapset.Id,
		Status: status,
	}
	beatmapUpdateCriteria := map[string]any{
		"set_id": beatmapset.Id,
		"status": constants.BeatmapStatusQualified,
	}

	app.Repositories.Beatmapsets.Update(beatmapsetUpdate, "submission_status", "approved_date")
	app.Repositories.Beatmaps.UpdateByCriteria(beatmapUpdateCriteria, beatmapUpdate, "status")

	logger.Info("Beatmapset was approved.", "id", beatmapset.Id, "name", beatmapset.Name(), "status", status)
}

func handlePendingSet(app *state.State, logger *slog.Logger, beatmapset *schemas.Beatmapset) {
	lastUpdate := time.Since(beatmapset.LastUpdate)
	graveyardTime := PendingToGraveyardTime

	if lastUpdate < graveyardTime {
		return
	}

	updateBeatmapIcon(app, beatmapset, constants.BeatmapStatusGraveyard, beatmapset.Status)
	moveBeatmapTopic(app, beatmapset, constants.BeatmapStatusGraveyard)

	beatmapsetUpdate := &schemas.Beatmapset{
		Id:     beatmapset.Id,
		Status: constants.BeatmapStatusGraveyard,
	}
	beatmapUpdate := &schemas.Beatmap{
		SetId:  beatmapset.Id,
		Status: constants.BeatmapStatusGraveyard,
	}

	app.Repositories.Beatmapsets.Update(beatmapsetUpdate, "submission_status")
	app.Repositories.Beatmaps.UpdateBySetId(beatmapUpdate, "status")
	app.Repositories.Nominations.DeleteAll(beatmapset.Id)

	logger.Info("Beatmapset was sent to the beatmap graveyard.", "id", beatmapset.Id, "name", beatmapset.Name())
}

func moveBeatmapTopic(app *state.State, beatmapset *schemas.Beatmapset, status constants.BeatmapStatus) {
	if beatmapset.TopicId == nil {
		return
	}

	var forumId int
	switch status {
	case constants.BeatmapStatusPending:
		forumId = BeatmapForumIdPending
	case constants.BeatmapStatusWIP:
		forumId = BeatmapForumIdWIP
	case constants.BeatmapStatusGraveyard, constants.BeatmapStatusInactive:
		forumId = BeatmapForumIdGraveyard
	default:
		forumId = BeatmapForumIdRanked
	}

	topicUpdate := &schemas.ForumTopic{Id: *beatmapset.TopicId, ForumId: forumId}
	postUpdate := &schemas.ForumPost{TopicId: *beatmapset.TopicId, ForumId: forumId}
	app.Repositories.ForumTopics.Update(topicUpdate, "forum_id")
	app.Repositories.ForumPosts.UpdateByTopic(postUpdate, "forum_id")
}

func updateBeatmapIcon(app *state.State, beatmapset *schemas.Beatmapset, status constants.BeatmapStatus, previousStatus constants.BeatmapStatus) {
	if beatmapset.TopicId == nil {
		return
	}
	topicUpdate := &schemas.ForumTopic{Id: *beatmapset.TopicId}

	if status == constants.BeatmapStatusRanked || status == constants.BeatmapStatusQualified || status == constants.BeatmapStatusLoved {
		// Set icon to heart for ranked, qualified & loved maps
		topicUpdate.IconId = constants.ForumIconHeart.Pointer()
		app.Repositories.ForumTopics.Update(topicUpdate, "icon")
		return
	}

	if status == constants.BeatmapStatusApproved {
		// Set icon to flame for approved maps
		topicUpdate.IconId = constants.ForumIconFire.Pointer()
		app.Repositories.ForumTopics.Update(topicUpdate, "icon")
		return
	}

	isRankedStatus := func(status constants.BeatmapStatus) bool {
		return status == constants.BeatmapStatusQualified ||
			status == constants.BeatmapStatusApproved ||
			status == constants.BeatmapStatusRanked ||
			status == constants.BeatmapStatusLoved
	}

	if isRankedStatus(previousStatus) {
		// Set icon to broken heart for maps that were approved but are no longer approved
		topicUpdate.IconId = constants.ForumIconHeartPop.Pointer()
		app.Repositories.ForumTopics.Update(topicUpdate, "icon")
		return
	}

	if status == constants.BeatmapStatusGraveyard {
		nominations, err := app.Repositories.Nominations.FetchBySet(beatmapset.Id)
		// Pop the bubble if the map was nominated but got graveyarded
		if err == nil && len(nominations) > 0 {
			topicUpdate.IconId = constants.ForumIconBubblePop.Pointer()
			app.Repositories.ForumTopics.Update(topicUpdate, "icon")
			return
		}
	}

	// Remove icon for all other status changes
	topicUpdate.IconId = nil
	app.Repositories.ForumTopics.Update(topicUpdate, "icon")
}

func hideScoresForSet(app *state.State, beatmapset *schemas.Beatmapset) {
	for _, beatmap := range beatmapset.Beatmaps {
		scoreUpdate := &schemas.Score{
			BeatmapId: beatmap.Id,
			StatusPP:  constants.ScoreStatusHidden,
			Hidden:    true,
		}
		app.Repositories.Scores.UpdateByBeatmapId(scoreUpdate, "status", "hidden")
	}
}
