package constants

import (
	"fmt"
	"strings"
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

func (grade Grade) Valid() bool {
	switch grade {
	case GradeXH, GradeSH, GradeX, GradeS, GradeA, GradeB, GradeC, GradeD, GradeF, GradeN:
		return true
	default:
		return false
	}
}

func (grade Grade) String() string {
	if !grade.Valid() {
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

type IntegrityFlags uint16

const IntegrityFlagsNone IntegrityFlags = 0

const (
	FlagFlashlightHackKeyBindingPressed IntegrityFlags = 1 << iota
	FlagSpeedHackDetected
	FlagIncorrectModValue
	FlagMultipleOsuClients
	FlagChecksumFailure
	FlagFlashlightChecksumIncorrect
	FlagOsuExecutableChecksum
	FlagMissingProcessesInList
	FlagFlashlightImageHack
	FlagSpinnerHack
	FlagTransparentWindow
	FlagFastPress
	FlagRawMouseDiscrepancy
	FlagRawKeyboardDiscrepancy
)

const allIntegrityFlags = FlagFlashlightHackKeyBindingPressed |
	FlagSpeedHackDetected |
	FlagIncorrectModValue |
	FlagMultipleOsuClients |
	FlagChecksumFailure |
	FlagFlashlightChecksumIncorrect |
	FlagOsuExecutableChecksum |
	FlagMissingProcessesInList |
	FlagFlashlightImageHack |
	FlagSpinnerHack |
	FlagTransparentWindow |
	FlagFastPress |
	FlagRawMouseDiscrepancy |
	FlagRawKeyboardDiscrepancy

const suspiciousIntegrityFlags = FlagFlashlightImageHack |
	FlagSpinnerHack |
	FlagTransparentWindow |
	FlagFastPress |
	FlagFlashlightChecksumIncorrect |
	FlagChecksumFailure |
	FlagRawMouseDiscrepancy |
	FlagRawKeyboardDiscrepancy

var integrityFlagNames = [...]struct {
	flag IntegrityFlags
	name string
}{
	{FlagFlashlightHackKeyBindingPressed, "FlashlightHackKeyBindingPressed"},
	{FlagSpeedHackDetected, "SpeedHackDetected"},
	{FlagIncorrectModValue, "IncorrectModValue"},
	{FlagMultipleOsuClients, "MultipleOsuClients"},
	{FlagChecksumFailure, "ChecksumFailure"},
	{FlagFlashlightChecksumIncorrect, "FlashlightChecksumIncorrect"},
	{FlagOsuExecutableChecksum, "OsuExecutableChecksum"},
	{FlagMissingProcessesInList, "MissingProcessesInList"},
	{FlagFlashlightImageHack, "FlashlightImageHack"},
	{FlagSpinnerHack, "SpinnerHack"},
	{FlagTransparentWindow, "TransparentWindow"},
	{FlagFastPress, "FastPress"},
	{FlagRawMouseDiscrepancy, "RawMouseDiscrepancy"},
	{FlagRawKeyboardDiscrepancy, "RawKeyboardDiscrepancy"},
}

// Has reports whether any flag in mask is present.
func (flags IntegrityFlags) Has(mask IntegrityFlags) bool {
	return flags&mask != 0
}

// Suspicious reports whether any suspicious flag is present.
func (flags IntegrityFlags) Suspicious() bool {
	return flags.Has(suspiciousIntegrityFlags)
}

// String returns the human-readable names of the enabled flags in bit order.
func (flags IntegrityFlags) String() string {
	if flags == IntegrityFlagsNone {
		return "None"
	}

	// Append the names of known flags to a string slice
	parts := make([]string, 0, len(integrityFlagNames)+1)
	for _, entry := range integrityFlagNames {
		if flags.Has(entry.flag) {
			parts = append(parts, entry.name)
		}
	}

	// Append unknown flags to string too
	if unknown := flags &^ allIntegrityFlags; unknown != IntegrityFlagsNone {
		parts = append(parts, fmt.Sprintf("Unknown(%#x)", uint16(unknown)))
	}
	return strings.Join(parts, ", ")
}
