package capabilities

import "github.com/osuTitanic/titanic/services/anchor/internal/models"

type Spectator interface {
	DeliverFrames(models.ReplayFrameBundle) error
	DeliverFellowSpectatorJoined(models.UserInfo) error
	DeliverFellowSpectatorLeft(models.UserInfo) error
	DeliverSpectatingFailure(models.UserInfo) error
}

type SpectatorHost interface {
	DeliverSpectatorJoined(models.UserInfo) error
	DeliverSpectatorLeft(models.UserInfo) error
	DeliverSpectatorUnavailable(models.UserInfo) error
}
