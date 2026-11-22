package capabilities

import "github.com/osuTitanic/titanic/services/anchor/internal/models"

type Chat interface {
	DeliverMessage(models.Message) error
	DeliverAwayNotice(string) error
}

type Channels interface {
	DeliverChannelList([]models.Channel) error
	DeliverChannelUpdate(models.Channel) error
	DeliverChannelJoined(models.Channel) error
	DeliverChannelLeft(models.Channel) error

	DeliverChannelMembers(models.Channel, []models.UserInfo) error
	DeliverChannelMemberJoined(models.Channel, models.UserInfo) error
	DeliverChannelMemberLeft(models.Channel, models.UserInfo) error
}
