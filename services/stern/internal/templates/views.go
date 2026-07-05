package templates

import (
	"net/url"
	"time"

	"github.com/osuTitanic/titanic-go/internal/config"
	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/permissions"
	"github.com/osuTitanic/titanic-go/internal/rankings"
	"github.com/osuTitanic/titanic-go/internal/schemas"
)

type Statistics struct {
	TotalUsers     int
	TotalScores    int
	OnlineUsersOsu int
	OnlineUsersIrc int
}

func (stats *Statistics) OnlineUsers() int {
	return stats.OnlineUsersOsu + stats.OnlineUsersIrc
}

type DefaultView struct {
	Query             url.Values
	Config            *config.Config
	CurrentUser       *schemas.User
	Permissions       *permissions.Set
	Stats             *Statistics
	CSRFToken         string
	CurrentPath       string
	CurrentURI        string
	NotificationCount int
	IsModernBrowser   bool
	IsIE              bool
}

func (v DefaultView) IsAuthenticated() bool {
	return v.CurrentUser != nil
}

func (v DefaultView) CurrentUserId() int {
	if v.CurrentUser == nil {
		return 0
	}
	return v.CurrentUser.Id
}

func (v DefaultView) PreferredModeAlias() string {
	if v.CurrentUser == nil {
		return constants.ModeOsu.Alias()
	}
	return v.CurrentUser.PreferredMode.Alias()
}

type ErrorMessageView struct {
	DefaultView
	Title   string
	Heading string
	Message string
}

type HomeView struct {
	DefaultView
	News               []*schemas.ForumPost
	ChatMessages       []*schemas.Message
	MostPlayedBeatmaps map[int]*schemas.Beatmap
}

type LoginView struct {
	DefaultView
	Redirect     string
	ErrorMessage string
}

type RegisterView struct {
	DefaultView
	ErrorMessage     string
	RecaptchaEnabled bool
	RecaptchaSiteKey string
}

type ResetView struct {
	DefaultView
	Redirect     string
	ErrorMessage string
}

type VerificationView struct {
	DefaultView
	Verification *schemas.Verification
	Success      bool
	Reset        bool
	InfoMessage  string
	ErrorMessage string
}

type UserProfileView struct {
	DefaultView
	User          *schemas.User
	Mode          constants.Mode
	IsOnline      bool
	Followers     int
	TotalPosts    int
	PPRank        int
	PPRankCountry int
	CurrentAdded  bool // current user friended this profile
	TargetAdded   bool // profile friended current user
	IsBlocked     bool
	SuperFriendly bool
	General       *UserGeneralTab
}

func (v UserProfileView) IsOwnProfile() bool {
	return v.CurrentUser != nil && v.CurrentUser.Id == v.User.Id
}

func (v UserProfileView) IsOtherProfile() bool {
	return v.CurrentUser != nil && v.CurrentUser.Id != v.User.Id
}

type UserGeneralTab struct {
	User           *schemas.User
	Mode           constants.Mode
	Stats          *schemas.Stats
	PPRank         int
	PPRankCountry  int
	ScoreRank      int
	TotalScoreRank int
	PPv1Rank       int
	TotalKudosu    int
	Activity       *UserActivityPage
}

// HasStats checks if the user has stats worth rendering
// in the general tab for the selected mode.
func (t *UserGeneralTab) HasStats() bool {
	return t.Stats != nil && t.Stats.Playcount > 0 && !t.User.Restricted
}

type UserActivityPage struct {
	UserId     int
	Mode       constants.Mode
	Rows       []*schemas.Activity
	Offset     int
	NextOffset int
	HasMore    bool
}

func (p *UserActivityPage) IsFirstPage() bool {
	return p.Offset == 0
}

type UserTopPlaysTab struct {
	UserId     int
	Mode       constants.Mode
	IsOwner    bool
	FirstsRank int
	Pinned     *UserScorePage
	Best       *UserScorePage
	First      *UserScorePage
}

type UserScorePage struct {
	UserId          int
	Mode            constants.Mode
	Section         string // "pinned" | "best" | "first"
	Scores          []*schemas.Score
	Offset          int
	NextOffset      int
	HasMore         bool
	Total           int
	IsOwner         bool
	ApprovedRewards bool
}

func (p *UserScorePage) IsFirstPage() bool {
	return p.Offset == 0
}

func (p *UserScorePage) ShowWeight() bool {
	return p.Section == "best"
}

type UserHistoryTab struct {
	UserId     int
	Mode       constants.Mode
	MostPlayed []*schemas.BeatmapPlays
	Recent     []*schemas.Score
}

type UserKudosuTab struct {
	UserId      int
	TotalKudosu int
	Entries     []*UserKudosuEntry
}

type UserKudosuEntry struct {
	Time        time.Time
	Status      string // "received" | "gave" | "revoked"
	Preposition string // "from" | "to"
	Amount      int
	ActorId     int
	ActorName   string
	OtherId     int
	OtherName   string
	PostId      int64
	TopicTitle  string
}

type UserAchievementsTab struct {
	UserId        int
	UnlockedCount int
	Categories    []*UserAchievementCategory
}

type UserAchievementCategory struct {
	Name         string
	Achievements []*UserAchievement
}

type UserAchievement struct {
	Name       string
	Unlocked   bool
	Filename   string
	UnlockedAt time.Time
}

type UserBeatmapsTab struct {
	UserId         int
	IsOwner        bool
	Favourites     []*schemas.BeatmapFavourite
	Created        []*UserBeatmapGroup
	Collaborations []*schemas.BeatmapCollaboration
	Nominations    []*schemas.BeatmapNomination
}

type UserBeatmapGroup struct {
	Name        string // e.g. "Ranked", "Graveyarded"
	Key         string // e.g. "ranked", "graveyarded"
	CanEdit     bool
	Revivable   bool
	Beatmapsets []*schemas.Beatmapset
}

type GroupView struct {
	DefaultView
	Group *schemas.Group
	Users []*schemas.User
}

type ForumHomeView struct {
	DefaultView
	Sections []*ForumSection
}

type ForumSection struct {
	Forum     *schemas.Forum
	Subforums []*ForumSubforum
}

type ForumSubforum struct {
	Forum         *schemas.Forum
	Recent        *schemas.ForumPost
	CurrentUserId int
}

type ForumView struct {
	DefaultView
	Forum          *schemas.Forum
	Parents        []*schemas.Forum
	Subforums      []*schemas.Forum
	SubforumRecent map[int]*schemas.ForumPost
	Announcements  []*ForumTopicPreview
	Topics         []*ForumTopicPreview
	ActiveUsers    []*ForumActiveUser
	HasCustomIcons bool
	TopicCount     int
	CanCreateTopic bool
	Pagination     PaginationView
}

func (v ForumView) HasTopics() bool {
	return len(v.Announcements) > 0 || len(v.Topics) > 0
}

type ForumTopicPreview struct {
	Topic          *schemas.ForumTopic
	LastPost       *schemas.ForumPost
	StatusIcon     string
	PageCount      int
	Index          int
	ForumId        int
	HasCustomIcons bool
	CurrentUserId  int
}

func (p ForumTopicPreview) PreviewTruncated() bool {
	// This will determine if we show the
	// ellipsis (...) after the first page
	return p.PageCount > 4
}

func (p ForumTopicPreview) PreviewPages() []int {
	if p.PageCount <= 4 {
		// For shorter topics, show all pages
		// e.g. [1, 2, 3, 4]
		pages := make([]int, 0, p.PageCount)
		for page := 1; page <= p.PageCount; page++ {
			pages = append(pages, page)
		}
		return pages
	}

	// For longer topics, show the first page and the last three
	// e.g. if a topic has 10 pages, this will return [1, 8, 9, 10]
	return []int{1, p.PageCount - 2, p.PageCount - 1, p.PageCount}
}

type ForumActiveUser struct {
	Id   int
	Name string
}

type ForumTopicView struct {
	DefaultView
	Forum           *schemas.Forum
	Topic           *schemas.ForumTopic
	Parents         []*schemas.Forum
	Posts           []*ForumPostPreview
	Pagination      PaginationView
	ActiveUsers     []*ForumActiveUser
	Beatmapset      *schemas.Beatmapset
	PostCount       int
	IsSubscribed    bool
	IsBookmarked    bool
	CanCreatePosts  bool
	CanReply        bool
	ReplyLocked     bool
	MetaDescription string
	MetaImage       string
}

func (v ForumTopicView) TopicLocked() bool {
	return v.Topic.LockedAt != nil
}

func (v ForumTopicView) HasBeatmapset() bool {
	return v.Beatmapset != nil
}

type ForumPostPreview struct {
	Post        *schemas.ForumPost
	Icon        *schemas.ForumIcon
	AuthorTitle string
	PostCount   int
	CanEdit     bool
	CanDelete   bool
	CanQuote    bool

	BeatmapsetId    int
	ShowKudosuBox   bool
	CanManageKudosu bool
	CanResetKudosu  bool
	CanRevokeKudosu bool
	KudosuTotal     int
	LatestKudosu    *schemas.BeatmapModding
}

func (p ForumPostPreview) HasKudosuExcludedIcon() bool {
	if p.Icon == nil {
		return false
	}

	// Exclude kudosu awards for bubble / ranking / qualification
	switch p.Icon.Id {
	case constants.ForumIconHeart, constants.ForumIconBubble, constants.ForumIconFire:
		return true
	default:
		return false
	}
}

func (p ForumPostPreview) KudosuStatusColor() string {
	switch {
	case p.KudosuTotal > 0:
		return "green"
	case p.KudosuTotal == 0:
		return "black"
	default:
		return "red"
	}
}

func (p ForumPostPreview) AbsoluteKudosuTotal() int {
	if p.KudosuTotal < 0 {
		return -p.KudosuTotal
	}
	return p.KudosuTotal
}

type ForumCreateTopicView struct {
	DefaultView
	Forum   *schemas.Forum
	Parents []*schemas.Forum
	Editor  ForumEditorContext
}

type ForumPostEditorView struct {
	DefaultView
	Forum    *schemas.Forum
	Topic    *schemas.ForumTopic
	Parents  []*schemas.Forum
	Editor   ForumEditorContext
	Action   string
	ActionId int64
}

type ForumEditorIcon struct {
	Id       int
	Name     string
	Location string
	Selected bool
}

// this is a very verbose struct. i'll have to see if it can be simplified later, but for now it works
// as i wanted to outline the editor template first, before writing any handler code

type ForumEditorContext struct {
	Content    string
	SubmitText string
	CancelUrl  string // if empty -> no cancel link
	DraftUrl   string // if empty -> no save-draft button
	FocusBody  bool   // autofocus the body instead of the subject

	ShowSubject bool
	Subject     string

	ShowIcons        bool
	NoneIconSelected bool
	Icons            []*ForumEditorIcon

	ShowControls    bool
	ShowStatusInput bool
	StatusText      string
	NotifyChecked   bool
	ShowLockTopic   bool
	TopicLocked     bool
	ShowLockPost    bool
	PostLocked      bool
	ShowTopicTypes  bool
	TopicType       string // "global" | "pinned" | "announcement"

	ShowKudosuHint     bool
	KudosuReward       int
	ShowKudosuIconNote bool
	BeatmapsetId       int
}

func (v ForumEditorContext) HasContent() bool {
	return v.Content != ""
}

type DownloadView struct {
	DefaultView
	SelectedCategory string
	Categories       []*DownloadCategory
	Clients          []*schemas.Release
}

type DownloadCategory struct {
	Name     string
	Url      string
	Selected bool
}

type BeatmapView struct {
	DefaultView
	Beatmap               *schemas.Beatmap
	Beatmapset            *schemas.Beatmapset
	Mode                  constants.Mode
	Mods                  string
	Scores                []*schemas.Score
	PersonalBest          *schemas.Score
	PersonalBestRank      int
	Favourites            []*schemas.BeatmapFavourite
	FavouritesCount       int
	Favourited            bool
	Collaborations        []*schemas.BeatmapCollaboration
	Nominations           []*schemas.BeatmapNomination
	Friends               map[int]bool
	CollaborationRequests []*schemas.BeatmapCollaborationRequest
	Invite                *schemas.BeatmapCollaborationRequest
	IsBeatmapAuthor       bool
	BatNominated          bool
}

type BeatmapPacksView struct {
	DefaultView
	BeatmapPacks     []*schemas.BeatmapPack
	Categories       []string
	CategorySelected string
}

type BeatmapSearchView struct {
	DefaultView
	Beatmapsets []*schemas.Beatmapset
	SearchSort  string
	SearchOrder string
	Pagination  PaginationView
}

type ChangelogOsumeView struct {
	Config *config.Config
	Groups []*ChangelogGroup
}

type ChangelogGroup struct {
	Date    string // formatted as "Jan 2, 2006"
	Entries []*schemas.ReleaseChangelog
}

type ScoreView struct {
	DefaultView
	Score      *schemas.Score
	User       *schemas.User
	UserStats  *schemas.Stats
	Beatmap    *schemas.Beatmap
	Beatmapset *schemas.Beatmapset
	ScoreRank  int
}

type MatchView struct {
	DefaultView
	Match  *schemas.Match
	Events []*MatchEventView
}

func (v MatchView) IsOngoing() bool {
	// Check if match is still running on bancho right now
	return v.Match.EndedAt == nil
}

type MatchEventView struct {
	Time time.Time
	Type constants.MatchEventType
	User *MatchEventUser // affected user for join/leave/kick/host events
	Game *MatchGameView  // populated for result events
}

type MatchEventUser struct {
	Id   int
	Name string
}

type MatchGameView struct {
	BeatmapId   int
	BeatmapText string
	Mode        constants.Mode
	TeamMode    constants.MatchTeamType
	ScoringMode constants.MatchScoringType
	Duration    string
	Results     []*MatchGameResult
	TeamResult  *MatchTeamResult // populated for team vs. games
}

type MatchGameResult struct {
	Place     int
	UserId    int
	Username  string
	Country   string
	Team      constants.SlotTeam
	Mods      constants.Mods
	Score     int
	Accuracy  float64
	MaxCombo  int
	Count300  int
	Count100  int
	Count50   int
	CountMiss int
	Failed    bool
}

type MatchTeamResult struct {
	Blue   string // formatted team value, e.g. "1,234,567" or "98.32%"
	Red    string
	Winner constants.SlotTeam
	Margin string // formatted difference, e.g. "12,345 points"
}

func (r MatchTeamResult) WinnerName() string {
	if r.Winner == constants.SlotTeamNeutral {
		return "Draw"
	}
	return r.Winner.String()
}

type RankingsView struct {
	DefaultView
	Pagination    PaginationView
	Country       string
	CountryName   string
	Location      string
	Mode          constants.Mode
	Type          constants.RankingType
	Entries       []*RankingEntry
	TopCountries  []string
	JumpTo        string
	TotalBeatmaps int
}

type RankingEntry struct {
	User     *schemas.User
	Score    int
	Rank     int
	IsFriend bool
}

type CountryRankingsView struct {
	DefaultView
	Pagination PaginationView
	Country    string
	Mode       constants.Mode
	Type       constants.RankingType
	Entries    []*rankings.CountryRanking
}

type KudosuView struct {
	DefaultView
	Pagination PaginationView
	Entries    []*KudosuEntry
	JumpTo     string
}

type KudosuEntry struct {
	User     *schemas.User
	Kudosu   int64
	Rank     int
	IsFriend bool
}
