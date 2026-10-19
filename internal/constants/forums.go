package constants

const (
	// Titanic!
	ForumDevelopment     = 3
	ForumGameplay        = 4
	ForumSkinning        = 5
	ForumFeatureRequests = 6
	ForumSupport         = 7
	// Beatmaps
	ForumBeatmapsRanked     = 8
	ForumBeatmapsPending    = 9
	ForumBeatmapsWIP        = 10
	ForumBeatmapsRequests   = 11
	ForumBeatmapsGraveyard  = 12
	ForumBeatmapsDiscussion = 13
	// Other
	ForumGeneralDiscussion = 20
	ForumOffTopic          = 21
	ForumIntroductions     = 22
	ForumClientModding     = 23
	ForumVideoGames        = 24
	ForumArt               = 25
	// Sub-Forums
	ForumBugReports               = 18
	ForumTournamentsAndContests   = 26
	ForumBeatmapRequestsCompleted = 28
	ForumBeatmapRequestsRejected  = 29
)

var BeatmapForumIds = map[int]bool{
	ForumBeatmapsRanked:    true,
	ForumBeatmapsPending:   true,
	ForumBeatmapsWIP:       true,
	ForumBeatmapsGraveyard: true,
}

type ForumIcon int

const (
	ForumIconHeart       ForumIcon = 1
	ForumIconHeartPop    ForumIcon = 2
	ForumIconBubble      ForumIcon = 3
	ForumIconBubblePop   ForumIcon = 4
	ForumIconFire        ForumIcon = 5
	ForumIconStar        ForumIcon = 6
	ForumIconRadioactive ForumIcon = 7
	ForumIconAlert       ForumIcon = 8
	ForumIconInfo        ForumIcon = 9
	ForumIconQuestion    ForumIcon = 10
	ForumIconOsu         ForumIcon = 11
	ForumIconTaiko       ForumIcon = 12
	ForumIconCatch       ForumIcon = 13
	ForumIconMania       ForumIcon = 14
)

func (icon ForumIcon) Pointer() *ForumIcon {
	return &icon
}
