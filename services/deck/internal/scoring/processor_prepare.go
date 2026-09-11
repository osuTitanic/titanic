package scoring

import (
	"fmt"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

func (processor *Processor) Prepare(password string) (ResultType, error) {
	// Authenticate the user & ensure they're eligible to submit a score
	user, err := processor.Context.AuthenticateUser(
		processor.Submission.Username, password,
		true, // We want them to be online on bancho
	)
	if err != nil {
		switch err {
		case server.ErrUserNotFound:
			return ResultUserNotFound, nil
		case server.ErrInvalidPassword:
			return ResultInvalidPassword, nil
		case server.ErrBanchoPresenceNotFound:
			return ResultBanchoUnavailable, nil
		default:
			return ResultAccepted, err
		}
	}

	if !user.Activated {
		return ResultInactive, nil
	}
	if user.Restricted {
		return ResultBanned, nil
	}
	if user.IsBot {
		return ResultRejected, nil
	}
	processor.Submission.User = user

	// TODO: Check scores.submit permission
	// TODO: Load bancho user presence from redis
	// TODO: Validate client hash (and processes maybe?)
	// TODO: Validate & set client version

	// Resolve the beatmap that the user set a score on
	beatmap, err := processor.Repositories.Beatmaps.ByChecksum(
		processor.Submission.BeatmapChecksum, "Beatmapset",
	)
	if err != nil {
		return ResultAccepted, nil
	}
	if beatmap == nil {
		return ResultBeatmapUnavailable, nil
	}
	if beatmap.Status == constants.BeatmapStatusInactive {
		return ResultBeatmapUnavailable, nil
	}
	if beatmap.Beatmapset.Status == constants.BeatmapStatusInactive {
		return ResultBeatmapUnavailable, nil
	}
	processor.Submission.Beatmap = beatmap

	// Populate charts for the modular submission response
	processor.Submission.Charts.Beatmap.BeatmapId = beatmap.Id
	processor.Submission.Charts.Beatmap.BeatmapSetId = beatmap.SetId
	processor.Submission.Charts.Beatmap.BeatmapPlaycount = beatmap.Playcount
	processor.Submission.Charts.Beatmap.BeatmapPasscount = beatmap.Passcount
	processor.Submission.Charts.Beatmap.ApprovedDate = beatmap.Beatmapset.ApprovedAt

	baseUrl := processor.Context.State.Config.OsuBaseUrl()
	processor.Submission.Charts.Overall.ChartUrl = fmt.Sprintf(
		"%s/p/playerranking/?f=%s",
		baseUrl, user.Name,
	)
	processor.Submission.Charts.Ranking.ChartUrl = fmt.Sprintf(
		"%s/b/%d",
		baseUrl, beatmap.Id,
	)
	return ResultAccepted, nil
}
