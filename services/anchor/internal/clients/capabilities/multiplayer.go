package capabilities

import "github.com/osuTitanic/titanic/services/anchor/internal/models"

type Lobby interface {
	DeliverMatchAvailable(models.Match) error
	DeliverMatchRemoved(models.Match) error
	DeliverLobbyUserJoined(models.UserInfo) error
	DeliverLobbyUserLeft(models.UserInfo) error
}

type MatchObserver interface {
	DeliverObservedMatch(models.Match) error
	DeliverObservedMatchClosed(int32) error
}

type MatchInvites interface {
	DeliverMatchInvite(models.Match, models.UserInfo) error
}
