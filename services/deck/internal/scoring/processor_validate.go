package scoring

import (
	"strings"

	"github.com/osuTitanic/titanic/internal/clients"
)

func (processor *Processor) validate() (ResultType, error) {
	if !processor.submission.CanSubmitScores {
		return ResultRejected, nil
	}

	// TODO: Normalize score values
	// TODO: Validate hit counts, total score, combo & mode
	// TODO: Reject non-whitelisted builds, unranked mods, invalid mods, etc.
	// TODO: Check for duplicate scores
	// TODO: For passed scores, validate the replay
	// TODO: Run touchscreen detection
	// TODO: Check pp limit for user

	// TODO: Add a run() wrapper similar to finalize()
	result, err := processor.validateClient()
	if err != nil {
		return result, err
	}
	if result != ResultAccepted {
		return result, nil
	}

	processor.calculatePPv2()
	processor.calculatePPv1()
	return ResultAccepted, nil
}

func (processor *Processor) calculatePPv2() {
	if !processor.context.State.PPv2.Available() {
		return
	}

	pp, err := processor.context.State.PPv2.CalculatePerformance(processor.submission.Score)
	if err != nil {
		processor.AddWarning("calculate ppv2: %w", err)
		return
	}
	processor.submission.PP = pp
}

func (processor *Processor) calculatePPv1() {
	if !processor.submission.Passed {
		return
	}

	pp, err := processor.context.State.PPv1.CalculatePerformance(processor.submission.Score)
	if err != nil {
		processor.AddWarning("calculate ppv1: %w", err)
		return
	}
	processor.submission.PPv1 = pp
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
