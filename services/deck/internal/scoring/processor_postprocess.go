package scoring

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
)

// PostProcess runs submission tasks after the response has been written.
// This includes things like uploading the replay, broadcasting announcements, etc.
func (processor *Processor) PostProcess() {
	if !processor.committed {
		return
	}

	run := func(name string, action func() error) {
		start := time.Now()
		err := action()
		if err != nil {
			processor.context.Logger.Warn(
				"Score submission post-response task failed",
				"step", name, "error", err,
			)
		}

		processor.context.Logger.Debug(
			"Score submission post-response task completed",
			"step", name, "took", time.Since(start),
		)
	}

	run("bancho user update", processor.banchoUserUpdate)
	run("update rank history", processor.updateRankHistory)
	run("broadcast rank highlight", processor.broadcastRankHighlight)
	run("broadcast beatmap highlight", processor.broadcastBeatmapHighlight)
	run("broadcast performance highlight", processor.broadcastPerformanceHighlight)
	run("upload replay", processor.uploadReplay)
}

func (processor *Processor) updateRankHistory() error {
	if processor.context.State.Config.FrozenRankUpdates {
		return nil
	}
	_, err := processor.repositories.Histories.UpdateRank(
		processor.submission.CurrentStats,
		processor.submission.User.Country,
		processor.context.State.Rankings,
	)
	return err
}

func (processor *Processor) uploadReplay() error {
	submission := processor.submission
	if !submission.Passed || len(submission.Replay) == 0 || submission.Id <= 0 {
		return nil
	}
	if submission.StatusPP <= constants.ScoreStatusExited {
		return nil
	}
	if len(submission.Replay) > MaxReplaySize {
		return fmt.Errorf("replay exceeds 10 MB")
	}

	// TODO: Use NewBeatmapRank to determine if replay should be uploaded
	// TODO: Cache replay if not in rank range

	return processor.context.State.Storage.Save(
		strconv.FormatInt(submission.Id, 10),
		"replays",
		submission.Replay,
	)
}

func (processor *Processor) banchoUserUpdate() error {
	return processor.context.State.BanchoEvents.UserUpdate(
		context.WithoutCancel(processor.context.Request.Context()),
		processor.submission.UserId,
		processor.submission.Mode,
	)
}
