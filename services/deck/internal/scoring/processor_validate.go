package scoring

func (processor *Processor) validate() (ResultType, error) {
	// TODO: Normalize score values & mod combinations
	// TODO: Validate hit counts, total score, combo & mode
	// TODO: Reject non-whitelisted builds, unranked mods, invalid mods, etc.
	// TODO: Compare client hash with bancho
	// TODO: Check for duplicate scores
	// TODO: For passed scores, validate the replay
	// TODO: Run touchscreen detection
	// TODO: Check pp limit for user

	// TODO: Add a run() wrapper similar to finalize()
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
