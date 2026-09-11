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
		Context:      ctx,
		Submission:   submission,
		Repositories: ctx.State.Repositories,
	}
}

func (processor *Processor) Process(password string) (result Result, err error) {
	result = Result{
		Type:       ResultAccepted,
		Submission: processor.Submission,
		Warnings:   []error{},
	}

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

	err = processor.Context.State.DatabaseTransaction(func(repos *state.Repositories) error {
		result.Type, err = processor.Persist(repos)
		return err // Commit transaction, or rollback if err != nil
	})
	if err != nil {
		return result, err
	}
	if result.Rejected() {
		return result, nil
	}

	result.Warnings = processor.PostProcess()
	return result, err
}
