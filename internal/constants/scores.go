package constants

import (
	"fmt"
)

type ScoreStatus int

const (
	ScoreStatusHidden    ScoreStatus = -1
	ScoreStatusFailed    ScoreStatus = 0
	ScoreStatusExited    ScoreStatus = 1
	ScoreStatusSubmitted ScoreStatus = 2
	ScoreStatusBest      ScoreStatus = 3
	ScoreStatusMods      ScoreStatus = 4
)

type Grade string

const (
	GradeXH Grade = "XH"
	GradeSH Grade = "SH"
	GradeX  Grade = "X"
	GradeS  Grade = "S"
	GradeA  Grade = "A"
	GradeB  Grade = "B"
	GradeC  Grade = "C"
	GradeD  Grade = "D"
	GradeF  Grade = "F"
	GradeN  Grade = "N"
)

func (grade Grade) String() string {
	if !grade.IsValid() {
		return fmt.Sprintf("Unknown(%s)", string(grade))
	}
	return string(grade)
}

func (grade Grade) Value() int8 {
	switch grade {
	case GradeXH:
		return 0
	case GradeSH:
		return 1
	case GradeX:
		return 2
	case GradeS:
		return 3
	case GradeA:
		return 4
	case GradeB:
		return 5
	case GradeC:
		return 6
	case GradeD:
		return 7
	case GradeF:
		return 8
	case GradeN:
		return 9
	default:
		return -1
	}
}

func (grade Grade) IsValid() bool {
	switch grade {
	case GradeXH, GradeSH, GradeX, GradeS, GradeA, GradeB, GradeC, GradeD, GradeF, GradeN:
		return true
	default:
		return false
	}
}
