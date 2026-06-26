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
	Query       url.Values
	Config      *config.Config
	CurrentUser *schemas.User
	Permissions *permissions.Set
	Stats       *Statistics
	CSRFToken   string
	CurrentPath string
	CurrentURI  string
}

func (v DefaultView) IsAuthenticated() bool {
	return v.CurrentUser != nil
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

type BeatmapSearchView struct {
	DefaultView
	Beatmapsets []*schemas.Beatmapset
	SearchSort  string
	SearchOrder string
	Pagination  PaginationView
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

type ScoreView struct {
	DefaultView
	Score      *schemas.Score
	User       *schemas.User
	UserStats  *schemas.Stats
	Beatmap    *schemas.Beatmap
	Beatmapset *schemas.Beatmapset
	ScoreRank  int
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
	PostId      int
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

type BeatmapPacksView struct {
	DefaultView
	BeatmapPacks     []*schemas.BeatmapPack
	Categories       []string
	CategorySelected string
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
