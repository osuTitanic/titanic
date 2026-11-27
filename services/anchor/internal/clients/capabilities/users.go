package capabilities

import "github.com/osuTitanic/titanic/services/anchor/internal/models"

type Presence interface {
	DeliverUserPresence([]models.UserInfo) error
	DeliverUserOffline(models.UserInfo) error
}

type Stats interface {
	DeliverStats([]models.UserInfo) error
}

type Friends interface {
	DeliverFriends([]models.UserInfo) error
}

type Notifications interface {
	DeliverAnnouncement(string) error
	DeliverRestrictionNotice(int) error
	DeliverRestartNotice(int) error
}
