package scoring

import (
	"fmt"

	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

type Processor struct {
	context      *server.Context
	submission   *SubmissionContext
	repositories *state.Repositories
	warnings     []error
}

func (processor *Processor) AddWarning(format string, a ...any) {
	processor.warnings = append(processor.warnings, fmt.Errorf(format, a...))
}

func NewProcessor(ctx *server.Context, submission *SubmissionContext) *Processor {
	return &Processor{
		context:      ctx,
		submission:   submission,
		repositories: ctx.State.Repositories,
	}
}

func (processor *Processor) Process(password string) (result Result, err error) {
	result = Result{
		Type:       ResultAccepted,
		Submission: processor.submission,
	}

	result.Type, err = processor.prepare(password)
	if err != nil {
		return result, fmt.Errorf("prepare score submission: %w", err)
	}
	if result.Rejected() {
		return result, nil
	}

	result.Type, err = processor.validate()
	if err != nil {
		return result, fmt.Errorf("validate score submission: %w", err)
	}
	if result.Rejected() {
		return result, nil
	}

	err = processor.context.State.DatabaseTransaction(func(repos *state.Repositories) error {
		result.Type, err = processor.persist(repos)
		return err // Commit transaction, or rollback if err != nil
	})
	if err != nil {
		return result, err
	}
	processor.finalize()

	result.Warnings = processor.warnings
	return result, nil
}
