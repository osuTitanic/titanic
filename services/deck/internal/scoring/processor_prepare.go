package scoring

import (
	"fmt"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

func (processor *Processor) prepare(password string) (ResultType, error) {
	// Authenticate the user & ensure they're eligible to submit a score
	user, err := processor.context.AuthenticateUser(
		processor.submission.Username, password,
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
	processor.submission.User = user

	// TODO: Check scores.submit permission
	// TODO: Load bancho user presence from redis
	// TODO: Validate client hash (and processes maybe?)
	// TODO: Validate & set client version

	// Resolve the beatmap that the user set a score on
	beatmap, err := processor.repositories.Beatmaps.ByChecksum(
		processor.submission.BeatmapChecksum, "Beatmapset",
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
	processor.submission.Beatmap = beatmap

	// Populate charts for the modular submission response
	processor.submission.Charts.Beatmap.BeatmapId = beatmap.Id
	processor.submission.Charts.Beatmap.BeatmapSetId = beatmap.SetId
	processor.submission.Charts.Beatmap.BeatmapPlaycount = beatmap.Playcount
	processor.submission.Charts.Beatmap.BeatmapPasscount = beatmap.Passcount
	processor.submission.Charts.Beatmap.ApprovedDate = beatmap.Beatmapset.ApprovedAt

	baseUrl := processor.context.State.Config.OsuBaseUrl()
	processor.submission.Charts.Overall.ChartUrl = fmt.Sprintf(
		"%s/p/playerranking/?f=%s",
		baseUrl, user.Name,
	)
	processor.submission.Charts.Ranking.ChartUrl = fmt.Sprintf(
		"%s/b/%d",
		baseUrl, beatmap.Id,
	)
	return ResultAccepted, nil
}
