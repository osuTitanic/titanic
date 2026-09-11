package scoring

import (
	"fmt"

	"github.com/osuTitanic/titanic/internal/state"
)

type Processor struct {
	app        *state.State
	submission *SubmissionContext
}

func NewProcessor(app *state.State, submission *SubmissionContext) *Processor {
	return &Processor{app: app, submission: submission}
}

func (processor *Processor) Process(password string) (result Result, err error) {
	result.Type = ResultAccepted
	result.Submission = processor.submission

	result.Type, err = processor.Prepare(password)
	if err != nil {
		return result, fmt.Errorf("prepare score submission: %w", err)
	}
	if result.Rejected() {
		return result, nil
	}

	result.Type, err = processor.Validate()
	if err != nil {
		return result, fmt.Errorf("validate score submission: %w", err)
	}
	if result.Rejected() {
		return result, nil
	}

	// TODO: We should pass in a db transaction here
	result.Type, err = processor.Persist()
	if err != nil {
		return result, fmt.Errorf("persist score submission: %w", err)
	}
	if result.Rejected() {
		return result, nil
	}

	result.Warnings = processor.PostProcess()
	return result, nil
}
