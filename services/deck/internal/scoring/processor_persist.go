package scoring

import "github.com/osuTitanic/titanic/internal/state"

func (processor *Processor) Persist(transaction *state.Repositories) (ResultType, error) {
	// TODO: Take snapshot of user's stats before submission
	// TODO: Load the current pp / score pb's for this beatmap
	// TODO: Capture the previous user rank & beatmap rank
	// TODO: Determine whether the beatmap awards score/pp
	// TODO: Assign new score statuses
	// TODO: Insert the submitted score
	// TODO: Increment the beatmap playcount / passcount
	// TODO: Update user statistics
	// TODO: Commit the transaction (rollback in case of failures)
	return ResultAccepted, nil
}
