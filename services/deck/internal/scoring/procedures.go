package scoring

import (
	"context"
)

type Procedure struct {
	Id        ProcedureId
	DependsOn []ProcedureId
	Run       ProcedureFunc
}

type ProcedureId int

type ProcedureFunc func(context.Context, *SubmissionContext) error
