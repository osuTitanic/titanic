package capabilities

import "github.com/osuTitanic/titanic/services/anchor/internal/models"

type Presence interface {
	DeliverUserPresence([]models.UserPresence) error // TODO: I don't think we can use UserPresence here because it doesn't have the id
	DeliverUserOffline(models.UserPresence) error
}

type Stats interface {
	DeliverStats([]models.UserStats) error // TODO: I don't think we can use UserStats here because it doesn't have the id
}

type Friends interface {
	DeliverFriends([]models.UserInfo) error
}

type Notifications interface {
	DeliverAnnouncement(string) error
	DeliverRestrictionNotice(int) error
	DeliverRestartNotice(int) error
}
