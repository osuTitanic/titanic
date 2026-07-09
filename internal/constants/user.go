package constants

import "strings"

type Playstyle uint8

const (
	PlaystyleNotSpecified Playstyle = 0
	PlaystyleMouse        Playstyle = 1 << 0
	PlaystyleTablet       Playstyle = 1 << 1
	PlaystyleKeyboard     Playstyle = 1 << 2
	PlaystyleTouch        Playstyle = 1 << 3
)

func (p Playstyle) Has(flag Playstyle) bool {
	return p&flag != 0
}

func (p Playstyle) String() string {
	if p == PlaystyleNotSpecified {
		return "None"
	}

	parts := make([]string, 0, 4)
	if p.Has(PlaystyleMouse) {
		parts = append(parts, "Mouse")
	}
	if p.Has(PlaystyleTablet) {
		parts = append(parts, "Tablet")
	}
	if p.Has(PlaystyleKeyboard) {
		parts = append(parts, "Keyboard")
	}
	if p.Has(PlaystyleTouch) {
		parts = append(parts, "Touch")
	}

	if len(parts) == 0 {
		return "Unknown"
	}

	return strings.Join(parts, ",")
}

type UserActivity int

const (
	ActivityRanksGained             UserActivity = 1
	ActivityNumberOne               UserActivity = 2
	ActivityBeatmapLeaderboardRank  UserActivity = 3
	ActivityLostFirstPlace          UserActivity = 4
	ActivityPPRecord                UserActivity = 5
	ActivityTopPlay                 UserActivity = 6
	ActivityAchievementUnlocked     UserActivity = 7
	ActivityScoreSubmitted          UserActivity = 8
	ActivityBeatmapUploaded         UserActivity = 9
	ActivityBeatmapUpdated          UserActivity = 10
	ActivityBeatmapRevived          UserActivity = 11
	ActivityBeatmapFavouriteAdded   UserActivity = 12
	ActivityBeatmapFavouriteRemoved UserActivity = 13
	ActivityBeatmapRated            UserActivity = 14
	ActivityBeatmapCommented        UserActivity = 15
	ActivityBeatmapDownloaded       UserActivity = 16
	ActivityBeatmapStatusUpdated    UserActivity = 17
	ActivityBeatmapNominated        UserActivity = 18
	ActivityForumTopicCreated       UserActivity = 19
	ActivityForumPostCreated        UserActivity = 20
	ActivityForumSubscribed         UserActivity = 21
	ActivityForumUnsubscribed       UserActivity = 22
	ActivityForumBookmarked         UserActivity = 23
	ActivityForumUnbookmarked       UserActivity = 24
	ActivityOsuCoinsReceived        UserActivity = 25
	ActivityOsuCoinsUsed            UserActivity = 26
	ActivityFriendAdded             UserActivity = 27
	ActivityFriendRemoved           UserActivity = 28
	ActivityReplayWatched           UserActivity = 29
	ActivityScreenshotUploaded      UserActivity = 30
	ActivityUserRegistration        UserActivity = 31
	ActivityUserLogin               UserActivity = 32
	ActivityUserChatMessage         UserActivity = 33
	ActivityUserMatchCreated        UserActivity = 34
	ActivityUserMatchJoined         UserActivity = 35
	ActivityUserMatchLeft           UserActivity = 36
	ActivityBeatmapNuked            UserActivity = 37
)

var DisallowedUsernameSubstrings = []string{
	"blow job",
	"blowjob",
	"cockhead",
	"cocksucker",
	"cunt",
	"cunts",
	"dildo",
	"fag1t",
	"faget",
	"fagg1t",
	"faggit",
	"faggot",
	"fagit",
	"fags",
	"massterbait",
	"masstrbait",
	"masstrbate",
	"masterbaiter",
	"masterbate",
	"masterbates",
	"n1gr",
	"nigger",
	"nigur",
	"niiger",
	"niigr",
	"penis",
	"pussy",
	"slut",
	"whore",
	"nigga",
}
