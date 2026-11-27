package capabilities

import "github.com/osuTitanic/titanic/services/anchor/internal/models"

type Multiplayer interface {
	DeliverMatchAvailable(models.Match) error
	DeliverMatchRemoved(models.Match) error
	DeliverMatchInvite(models.Match, models.UserInfo) error
	DeliverMatchScore(models.ScoreFrame) error

	DeliverLobbyUserJoined(models.UserInfo) error
	DeliverLobbyUserLeft(models.UserInfo) error
}

type MultiplayerObserver interface {
	DeliverObservedMatch(models.Match) error
	DeliverObservedMatchClosed(int32) error
}
