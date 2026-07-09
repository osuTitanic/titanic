package state

import (
	"github.com/osuTitanic/titanic/internal/repositories"
	"gorm.io/gorm"
)

type Repositories struct {
	// Users
	Users         *repositories.UserRepository
	Stats         *repositories.StatsRepository
	Names         *repositories.NameRepository
	Badges        *repositories.BadgeRepository
	Stamps        *repositories.StampRepository
	Reports       *repositories.ReportRepository
	Infringements *repositories.InfringementRepository
	Relationships *repositories.RelationshipRepository
	Verifications *repositories.VerificationRepository
	Notifications *repositories.NotificationRepository
	Activities    *repositories.ActivityRepository
	Achievements  *repositories.AchievementRepository
	Screenshots   *repositories.ScreenshotRepository

	// Groups & Permissions
	Groups      *repositories.GroupRepository
	Permissions *repositories.PermissionsRepository

	// Bancho
	Activity *repositories.BanchoActivityRepository
	Logins   *repositories.LoginRepository
	Channels *repositories.ChannelRepository
	Filters  *repositories.ChatFilterRepository
	Clients  *repositories.ClientRepository
	Messages *repositories.MessageRepository
	Matches  *repositories.MatchRepository

	// Releases
	ReleasesTitanic  *repositories.ReleaseRepository
	ReleasesOfficial *repositories.ReleasesOfficialRepository
	ReleasesModded   *repositories.ModdedReleaseRepository
	ReleasesExtra    *repositories.ExtraReleaseRepository

	// Beatmaps
	Beatmaps        *repositories.BeatmapRepository
	Beatmapsets     *repositories.BeatmapsetRepository
	BeatmapPacks    *repositories.BeatmapPackRepository
	Ratings         *repositories.BeatmapRatingRepository
	Favourites      *repositories.BeatmapFavouriteRepository
	Nominations     *repositories.BeatmapNominationRepository
	Collaborations  *repositories.BeatmapCollaborationRepository
	Modding         *repositories.BeatmapModdingRepository
	Comments        *repositories.BeatmapCommentRepository
	Plays           *repositories.BeatmapPlaysRepository
	ResourceMirrors *repositories.ResourceMirrorRepository

	// Rankings
	Scores     *repositories.ScoreRepository
	Histories  *repositories.HistoryRepository
	Benchmarks *repositories.BenchmarkRepository

	// Wiki
	WikiPages      *repositories.WikiPageRepository
	WikiCategories *repositories.WikiCategoryRepository
	WikiContents   *repositories.WikiContentRepository
	WikiOutlinks   *repositories.WikiOutlinkRepository
	// TODO: I'm considering to merge these

	// Forums
	Forums           *repositories.ForumRepository
	ForumTopics      *repositories.ForumTopicRepository
	ForumPosts       *repositories.ForumPostRepository
	ForumBookmarks   *repositories.ForumBookmarkRepository
	ForumSubscribers *repositories.ForumSubscriberRepository
	ForumIcons       *repositories.ForumIconRepository
	ForumReports     *repositories.ForumReportRepository
	ForumStars       *repositories.ForumStarRepository
	// TODO: I'm considering to merge these
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Users:            repositories.NewUserRepository(db),
		Stats:            repositories.NewStatsRepository(db),
		Relationships:    repositories.NewRelationshipRepository(db),
		Badges:           repositories.NewBadgeRepository(db),
		Stamps:           repositories.NewStampRepository(db),
		Names:            repositories.NewNameRepository(db),
		Infringements:    repositories.NewInfringementRepository(db),
		Reports:          repositories.NewReportRepository(db),
		Verifications:    repositories.NewVerificationRepository(db),
		Groups:           repositories.NewGroupRepository(db),
		Permissions:      repositories.NewPermissionsRepository(db),
		Notifications:    repositories.NewNotificationRepository(db),
		Achievements:     repositories.NewAchievementRepository(db),
		Favourites:       repositories.NewBeatmapFavouriteRepository(db),
		Histories:        repositories.NewHistoryRepository(db),
		Beatmaps:         repositories.NewBeatmapRepository(db),
		Beatmapsets:      repositories.NewBeatmapsetRepository(db),
		Scores:           repositories.NewScoreRepository(db),
		Nominations:      repositories.NewNominationRepository(db),
		Messages:         repositories.NewMessageRepository(db),
		Activities:       repositories.NewActivityRepository(db),
		Benchmarks:       repositories.NewBenchmarkRepository(db),
		Channels:         repositories.NewChannelRepository(db),
		Clients:          repositories.NewClientRepository(db),
		Collaborations:   repositories.NewBeatmapCollaborationRepository(db),
		Comments:         repositories.NewBeatmapCommentRepository(db),
		Filters:          repositories.NewChatFilterRepository(db),
		Logins:           repositories.NewLoginRepository(db),
		Matches:          repositories.NewMatchRepository(db),
		Modding:          repositories.NewBeatmapModdingRepository(db),
		Plays:            repositories.NewBeatmapPlaysRepository(db),
		Ratings:          repositories.NewBeatmapRatingRepository(db),
		ReleasesTitanic:  repositories.NewReleaseRepository(db),
		ReleasesModded:   repositories.NewModdedReleaseRepository(db),
		ReleasesExtra:    repositories.NewExtraReleaseRepository(db),
		ReleasesOfficial: repositories.NewReleasesOfficialRepository(db),
		ResourceMirrors:  repositories.NewResourceMirrorRepository(db),
		Screenshots:      repositories.NewScreenshotRepository(db),
		Activity:         repositories.NewBanchoActivityRepository(db),
		BeatmapPacks:     repositories.NewBeatmapPackRepository(db),
		WikiPages:        repositories.NewWikiPageRepository(db),
		WikiCategories:   repositories.NewWikiCategoryRepository(db),
		WikiContents:     repositories.NewWikiContentRepository(db),
		WikiOutlinks:     repositories.NewWikiOutlinkRepository(db),
		Forums:           repositories.NewForumRepository(db),
		ForumTopics:      repositories.NewForumTopicRepository(db),
		ForumPosts:       repositories.NewForumPostRepository(db),
		ForumIcons:       repositories.NewForumIconRepository(db),
		ForumReports:     repositories.NewForumReportRepository(db),
		ForumStars:       repositories.NewForumStarRepository(db),
		ForumBookmarks:   repositories.NewForumBookmarkRepository(db),
		ForumSubscribers: repositories.NewForumSubscriberRepository(db),
	}
}
