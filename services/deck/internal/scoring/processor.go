package scoring

import (
	"fmt"

	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

type Processor struct {
	Context      *server.Context
	Submission   *SubmissionContext
	Repositories *state.Repositories
}

func NewProcessor(ctx *server.Context, submission *SubmissionContext) *Processor {
	return &Processor{
		Context:    ctx,
		Submission: submission,
	}
}

func (processor *Processor) Process(password string) (result Result, err error) {
	result = Result{
		Type:       ResultAccepted,
		Submission: processor.Submission,
		Warnings:   []error{},
	}

	err = processor.Context.State.DatabaseTransaction(func(repos *state.Repositories) error {
		processor.Repositories = repos

		result.Type, err = processor.Prepare(password)
		if err != nil {
			return fmt.Errorf("prepare score submission: %w", err)
		}
		if result.Rejected() {
			return nil
		}

		result.Type, err = processor.Validate()
		if err != nil {
			return fmt.Errorf("validate score submission: %w", err)
		}
		if result.Rejected() {
			return nil
		}

		result.Type, err = processor.Persist()
		if err != nil {
			return fmt.Errorf("persist score submission: %w", err)
		}
		if result.Rejected() {
			return nil
		}

		result.Warnings = processor.PostProcess()
		return nil
	})
	return result, err
}
