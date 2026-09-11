package scoring

func (processor *Processor) PostProcess() []error {
	// TODO: Synchronize the redis rankings after the database commit
	// TODO: Resolve the global, country, and beatmap ranks
	// TODO: Check & unlock achievements
	// TODO: Fill the charts with all the required data
	// TODO: Upload replay when the score was a pass
	// TODO: Notify bancho that user stats should be reloaded
	// TODO: Create rank, beatmap, and performance highlights / activity
	return nil
}
