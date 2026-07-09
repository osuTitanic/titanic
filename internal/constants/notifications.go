package constants

type NotificationType int

const (
	NotificationTypeWelcome NotificationType = iota
	NotificationTypeAchievement
	NotificationTypeChat
	NotificationTypeForum
	NotificationTypeBeatmaps
	NotificationTypeSecurity
)

const (
	WelcomeNotificationHeader  = "Welcome!"
	WelcomeNotificationContent = "Welcome aboard! Get started by downloading one of our game clients [here](%s/download). Enjoy your journey!"
)
