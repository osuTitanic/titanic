package scoring

import (
	"fmt"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/clients"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/replays"
	"github.com/osuTitanic/titanic/internal/schemas"
)

type Validator func() (ResultType, error)

func (processor *Processor) validate() (ResultType, error) {
	if !processor.submission.CanSubmitScores {
		return ResultRejected, nil
	}
	run := func(name string, action Validator) (ResultType, error) {
		start := time.Now()
		result, err := action()

		processor.context.Logger.Debug(
			"Score submission task completed",
			"step", name, "took", time.Since(start),
		)
		return result, err
	}

	if result, err := run("validate mods", processor.validateMods); wasRejected(result, err) {
		return result, err
	}
	if result, err := run("validate client", processor.validateClient); wasRejected(result, err) {
		return result, err
	}
	if result, err := run("check duplicate score", processor.checkDuplicateScore); wasRejected(result, err) {
		return result, err
	}
	if result, err := run("validate score", processor.validateScore); wasRejected(result, err) {
		return result, err
	}
	if result, err := run("check integrity flags", processor.checkIntegrityFlags); wasRejected(result, err) {
		return result, err
	}
	if result, err := run("validate replay", processor.validateReplay); wasRejected(result, err) {
		return result, err
	}
	if result, err := run("calculate ppv2", processor.calculatePPv2); wasRejected(result, err) {
		return result, err
	}
	if result, err := run("calculate ppv1", processor.calculatePPv1); wasRejected(result, err) {
		return result, err
	}
	if result, err := run("check pp limit", processor.checkPPLimit); wasRejected(result, err) {
		return result, err
	}
	if result, err := run("normalize relax score", processor.normalizeRelaxScore); wasRejected(result, err) {
		return result, err
	}
	return ResultAccepted, nil
}

func (processor *Processor) validateMods() (ResultType, error) {
	mods := processor.submission.Mods
	if !mods.Valid() {
		processor.AddWarning("unknown mod flags: %d", mods)
		return ResultRejected, nil
	}
	if !mods.ValidCombination() {
		processor.AddWarning("invalid mod combination: %s", mods)
		return ResultRejected, nil
	}
	if mods.Unranked() {
		return ResultRejected, nil
	}
	// TODO: Check mania mods in non-mania modes
	return ResultAccepted, nil
}

func (processor *Processor) calculatePPv2() (ResultType, error) {
	if !processor.context.State.PPv2.Available() {
		return ResultAccepted, nil
	}

	pp, err := processor.context.State.PPv2.CalculatePerformance(processor.submission.Score)
	if err != nil {
		processor.AddWarning("calculate ppv2: %w", err)
		return ResultAccepted, nil
	}
	processor.submission.PP = pp
	return ResultAccepted, nil
}

func (processor *Processor) calculatePPv1() (ResultType, error) {
	if !processor.submission.Passed {
		return ResultAccepted, nil
	}

	pp, err := processor.context.State.PPv1.CalculatePerformance(processor.submission.Score)
	if err != nil {
		processor.AddWarning("calculate ppv1: %w", err)
		return ResultAccepted, nil
	}
	processor.submission.PPv1 = pp
	return ResultAccepted, nil
}

func (processor *Processor) validateClient() (ResultType, error) {
	app := processor.context.State
	submission := processor.submission

	status, err := app.BanchoUsers.Get(
		processor.context.Request.Context(),
		submission.UserId,
	)
	if err != nil {
		return ResultAccepted, err
	}
	if status == nil {
		// User must be online on bancho for score submission
		return ResultBanchoUnavailable, nil
	}

	// Validate the client hash, if sent by client
	submission.ClientHash = strings.TrimSuffix(
		strings.TrimSpace(submission.ClientHash),
		":",
	)
	clientHashExists := submission.ClientHash != "" && status.ClientHash != ""
	clientHashMatching := strings.HasPrefix(status.ClientHash, submission.ClientHash)

	if clientHashExists && !clientHashMatching {
		processor.AddWarning(
			"client hash does not match: got `%s`, want `%s`",
			submission.ClientHash, status.ClientHash,
		)
	}

	// Ensure the client version matches, if sent by client
	clientVersionExists := submission.ClientVersion > 0
	clientVersionMatching := submission.ClientVersion == status.ClientVersion

	if clientVersionExists && !clientVersionMatching {
		processor.AddWarning(
			"client version does not match: got %d, want %d",
			submission.ClientVersion, status.ClientVersion,
		)
	}

	// Override the client version with data from bancho
	submission.ClientVersion = status.ClientVersion
	submission.ClientString = status.ClientString

	// Validate the client version itself now
	version, ok := clients.ParseVersion(submission.ClientString)
	if !ok {
		return ResultRejected, nil
	}

	if app.Config.DisableClientVerification {
		return ResultAccepted, nil
	}
	if submission.CanBypassClientValidation {
		return ResultAccepted, nil
	}
	if app.Config.BanchoClientCutoff > 0 && version.Date > app.Config.BanchoClientCutoff {
		return ResultRejected, nil
	}

	executableHash, _, _ := strings.Cut(status.ClientHash, ":")
	if executableHash == "" {
		return ResultRejected, nil
	}

	valid, err := clients.IsValidChecksum(app, version, executableHash)
	if err != nil {
		return ResultAccepted, err
	}
	if !valid {
		return ResultRejected, nil
	}
	return ResultAccepted, nil
}

func (processor *Processor) validateReplay() (ResultType, error) {
	submission := processor.submission
	if !submission.Passed {
		// Failed scores should not have replays
		return ResultAccepted, nil
	}
	if len(submission.Replay) == 0 {
		// hmmm... suspicious...
		processor.AddWarning("passed score has no replay")
		return ResultRejected, nil
	}

	// Decompress & parse the submitted replay frames to check their validity
	frames, _, err := replays.DeserializeFrames(submission.Replay)
	if err != nil {
		processor.AddWarning("invalid replay: %w", err)
		return ResultRejected, nil
	}
	if len(frames) < 100 {
		processor.AddWarning(
			"invalid replay: got %d frames, want at least %d",
			len(frames), 100,
		)
		return ResultRejected, nil
	}

	if submission.Mode != constants.ModeOsu {
		return ResultAccepted, nil
	}

	// Modern osu! clients have client-side touchscreen detection,
	// but on Titanic, we're talking pre-2017 clients, where this
	// kind of special technology did not exist.
	// (https://osu.ppy.sh/community/forums/topics/665986)

	// So, unfortunatly, we have to guess from the replay itself,
	// which is better than nothing, I guess.

	detected, score := replays.DetectTouchscreenUsage(frames, 0.8)
	submission.Touchscreen = detected

	if detected {
		processor.context.Logger.Debug(
			"touchscreen usage detected",
			"score", score,
			"checksum", submission.Checksum,
		)
	}
	return ResultAccepted, nil
}

func (processor *Processor) validateScore() (ResultType, error) {
	score := processor.submission.Score

	hitCounts := []struct {
		name  string
		value int
	}{
		// Check if any of those are < 0
		{"300", score.Count300},
		{"100", score.Count100},
		{"50", score.Count50},
		{"miss", score.CountMiss},
		{"geki", score.CountGeki},
		{"katu", score.CountKatu},
	}
	for _, hitCount := range hitCounts {
		if hitCount.value < 0 {
			processor.AddWarning("negative %s count: %d", hitCount.name, hitCount.value)
			return ResultRejected, nil
		}
	}

	if score.TotalObjects() <= 0 {
		processor.AddWarning("score did not pass any objects")
		return ResultRejected, nil
	}
	if score.TotalScore <= 0 {
		processor.AddWarning("invalid total score: %d", score.TotalScore)
		return ResultRejected, nil
	}
	if score.MaxCombo <= 0 {
		processor.AddWarning("invalid max combo: %d", score.MaxCombo)
		return ResultRejected, nil
	}
	if !score.Mode.Valid() {
		processor.AddWarning("invalid mode: %d", score.Mode)
		return ResultRejected, nil
	}

	if score.ClientVersion <= 504 && score.Mode == constants.ModeCatch {
		// b504 and below let players change the catcher size through their skin,
		// giving an unfair advantage when using a very wide catcher
		processor.AddWarning("unsupported catch submission: %s (%d)", score.ClientString, score.ClientVersion)
		return ResultRejected, nil
	}
	if score.ClientVersion < 452 && score.Mods.Has(constants.Nightcore) {
		// Prevent "Taiko" mod (later used as Nightcore) scores from being submitted
		processor.AddWarning("unsupported nightcore / taiko submission: %s (%d)", score.ClientString, score.ClientVersion)
		return ResultRejected, nil
	}

	// Check converts from "minigames" ;) to standard
	if score.Beatmap.Mode != constants.ModeOsu && score.Mode == constants.ModeOsu {
		processor.AddWarning(
			"invalid mode conversion: %s score on %s beatmap",
			score.Mode, score.Beatmap.Mode,
		)
		return ResultRejected, nil
	}
	return ResultAccepted, nil
}

func (processor *Processor) checkIntegrityFlags() (ResultType, error) {
	flags := processor.submission.Flags
	if flags.Suspicious() {
		processor.AddWarning("suspicious score integrity flags: %s", flags)
	}
	return ResultAccepted, nil
}

func (processor *Processor) checkPPLimit() (ResultType, error) {
	score := processor.submission.Score
	if !score.Beatmap.AwardsPP() {
		return ResultAccepted, nil
	}

	// PP limit will scale up based on account age
	limit := ppLimit(
		score.User.CreatedAt,
		time.Now(),
	)
	if score.PP >= limit {
		processor.AddWarning(
			"pp limit exceeded: got %.2f pp, limit %.2f pp",
			score.PP, limit,
		)
	}
	return ResultAccepted, nil
}

func (processor *Processor) checkDuplicateScore() (ResultType, error) {
	submission := processor.submission
	if !submission.Passed {
		return ResultAccepted, nil
	}
	if submission.ReplayMd5 == nil || *submission.ReplayMd5 == "" {
		return ResultAccepted, nil
	}

	duplicate, err := processor.repositories.Scores.ByReplayChecksum(*submission.ReplayMd5)
	if err != nil {
		return ResultAccepted, fmt.Errorf("find duplicate replay: %w", err)
	}

	result, warning := evaluateDuplicateScore(duplicate, submission.UserId)
	if warning != "" {
		processor.AddWarning("%s (checksum: %s)", warning, *submission.ReplayMd5)
	}
	return result, nil
}

func (processor *Processor) normalizeRelaxScore() (ResultType, error) {
	submission := processor.submission
	if !submission.Relaxing() {
		return ResultAccepted, nil
	}

	totalScore, err := EstimateTotalScore(submission.Score)
	if err != nil {
		return ResultAccepted, err
	}
	submission.TotalScore = totalScore
	return ResultAccepted, nil
}

func evaluateDuplicateScore(duplicate *schemas.Score, userId int) (ResultType, string) {
	if duplicate == nil {
		return ResultAccepted, ""
	}
	if duplicate.UserId == userId {
		return ResultRejected, fmt.Sprintf(
			"duplicate replay from same user: score %d",
			duplicate.Id,
		)
	}
	return ResultRejected, fmt.Sprintf(
		"duplicate replay from another user: score %d, user %d",
		duplicate.Id, duplicate.UserId,
	)
}

func ppLimit(createdAt, now time.Time) float64 {
	accountAgeSeconds := now.Sub(createdAt).Seconds()
	return min(1500.0, max(750.0, accountAgeSeconds/8))
}

func wasRejected(result ResultType, err error) bool {
	return result != ResultAccepted || err != nil
}
