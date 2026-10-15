package constants

type RelationshipStatus int

const (
	RelationshipStatusFriend RelationshipStatus = 0
	RelationshipStatusFoe    RelationshipStatus = 1
)

func (status RelationshipStatus) Valid() bool {
	switch status {
	case RelationshipStatusFriend, RelationshipStatusFoe:
		return true
	default:
		return false
	}
}

func (status RelationshipStatus) Value() int {
	return int(status)
}

func (status RelationshipStatus) String() string {
	switch status {
	case RelationshipStatusFriend:
		return "Friend"
	case RelationshipStatusFoe:
		return "Foe"
	default:
		return "Unknown"
	}
}
