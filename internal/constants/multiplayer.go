package constants

import "fmt"

type MatchEventType int8

const (
	MatchEventJoin    MatchEventType = 0
	MatchEventLeave   MatchEventType = 1
	MatchEventKick    MatchEventType = 2
	MatchEventHost    MatchEventType = 3
	MatchEventDisband MatchEventType = 4
	MatchEventStart   MatchEventType = 5
	MatchEventResult  MatchEventType = 6
	MatchEventAbort   MatchEventType = 7
)

type MatchScoringType int8

const (
	MatchScoringScore    MatchScoringType = 0
	MatchScoringAccuracy MatchScoringType = 1
	MatchScoringCombo    MatchScoringType = 2
)

func (scoringType MatchScoringType) String() string {
	switch scoringType {
	case MatchScoringScore:
		return "Score"
	case MatchScoringAccuracy:
		return "Accuracy"
	case MatchScoringCombo:
		return "Combo"
	default:
		return fmt.Sprintf("%d", scoringType)
	}
}

type MatchTeamType int8

const (
	MatchTeamHeadToHead MatchTeamType = 0
	MatchTeamTagCoop    MatchTeamType = 1
	MatchTeamTeamVs     MatchTeamType = 2
	MatchTeamTagTeamVs  MatchTeamType = 3
)

func (teamType MatchTeamType) String() string {
	switch teamType {
	case MatchTeamHeadToHead:
		return "Head to Head"
	case MatchTeamTagCoop:
		return "Tag Co-op"
	case MatchTeamTeamVs:
		return "Team VS"
	case MatchTeamTagTeamVs:
		return "Tag Team VS"
	default:
		return fmt.Sprintf("%d", teamType)
	}
}

func (teamType MatchTeamType) HasTeams() bool {
	return teamType == MatchTeamTeamVs || teamType == MatchTeamTagTeamVs
}

type SlotTeam int8

const (
	SlotTeamNeutral SlotTeam = 0
	SlotTeamBlue    SlotTeam = 1
	SlotTeamRed     SlotTeam = 2
)

func (team SlotTeam) String() string {
	switch team {
	case SlotTeamBlue:
		return "Blue"
	case SlotTeamRed:
		return "Red"
	default:
		return "None"
	}
}
