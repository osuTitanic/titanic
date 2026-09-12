package scoring

func (processor *Processor) validate() (ResultType, error) {
	// TODO: Normalize score values & mod combinations
	// TODO: Validate hit counts, total score, combo & mode
	// TODO: Reject non-whitelisted builds, unranked mods, invalid mods, etc.
	// TODO: Compare client hash with bancho
	// TODO: Check for duplicate scores
	// TODO: For passed scores, validate the replay
	// TODO: Run touchscreen detection
	// TODO: Calculate ppv2 / ppv2
	// TODO: Check pp limit for user
	return ResultAccepted, nil
}
