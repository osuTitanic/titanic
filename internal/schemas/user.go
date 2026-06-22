package schemas

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/osuTitanic/titanic-go/internal/constants"
)

type User struct {
	Id               int                   `gorm:"column:id;primaryKey;autoIncrement"`
	Name             string                `gorm:"column:name;unique"`
	SafeName         string                `gorm:"column:safe_name;unique"`
	Email            string                `gorm:"column:email;unique"`
	DiscordId        *int64                `gorm:"column:discord_id;unique"`
	Bcrypt           string                `gorm:"column:pw"`
	IsBot            bool                  `gorm:"column:bot;default:false"`
	Country          string                `gorm:"column:country"`
	SilenceEnd       *time.Time            `gorm:"column:silence_end"`
	CreatedAt        time.Time             `gorm:"column:created_at;autoCreateTime"`
	LatestActivity   time.Time             `gorm:"column:latest_activity;autoCreateTime"`
	Restricted       bool                  `gorm:"column:restricted;default:false"`
	Activated        bool                  `gorm:"column:activated;default:false"`
	PreferredMode    constants.Mode        `gorm:"column:preferred_mode;default:0"`
	PreferredRanking constants.RankingType `gorm:"column:preferred_ranking;default:global"`
	Playstyle        constants.Playstyle   `gorm:"column:playstyle;default:0"`
	IrcToken         string                `gorm:"column:irc_token;default:encode(gen_random_bytes(5),'hex')"`
	AvatarHash       *string               `gorm:"column:avatar_hash"`
	AvatarLastUpdate time.Time             `gorm:"column:avatar_last_changed;autoCreateTime"`
	FriendOnlyDMs    bool                  `gorm:"column:friendonly_dms;default:false"`

	Userpage  *string `gorm:"column:userpage_about"`
	Signature *string `gorm:"column:userpage_signature"`
	Title     *string `gorm:"column:userpage_title"`
	Banner    *string `gorm:"column:userpage_banner"`
	Website   *string `gorm:"column:userpage_website"`
	Discord   *string `gorm:"column:userpage_discord"`
	Twitter   *string `gorm:"column:userpage_twitter"`
	Location  *string `gorm:"column:userpage_location"`
	Interests *string `gorm:"column:userpage_interests"`

	TargetRelationships []*Relationship         `gorm:"foreignKey:TargetId;references:Id"`
	Relationships       []*Relationship         `gorm:"foreignKey:UserId;references:Id"`
	Collaborations      []*BeatmapCollaboration `gorm:"foreignKey:UserId;references:Id"`
	Nominations         []*BeatmapNomination    `gorm:"foreignKey:UserId;references:Id"`
	BookmarkedTopics    []*ForumBookmark        `gorm:"foreignKey:UserId;references:Id"`
	SubscribedTopics    []*ForumSubscriber      `gorm:"foreignKey:UserId;references:Id"`
	Notifications       []*Notification         `gorm:"foreignKey:UserId;references:Id"`
	Permissions         []*UserPermission       `gorm:"foreignKey:UserId;references:Id"`
	Achievements        []*Achievement          `gorm:"foreignKey:UserId;references:Id"`
	Screenshots         []*Screenshot           `gorm:"foreignKey:UserId;references:Id"`
	Favourites          []*BeatmapFavourite     `gorm:"foreignKey:UserId;references:Id"`
	Groups              []*GroupEntry           `gorm:"foreignKey:UserId;references:Id"`
	Badges              []*Badge                `gorm:"foreignKey:UserId;references:Id"`
	Stats               []*Stats                `gorm:"foreignKey:UserId;references:Id"`
	Names               []*Name                 `gorm:"foreignKey:UserId;references:Id"`
	Infringements       []*Infringement         `gorm:"foreignKey:UserId;references:Id"`
}

func (User) TableName() string {
	return "users"
}

func (user *User) Age() time.Duration {
	return time.Since(user.CreatedAt)
}

func (user *User) AgeDays() int {
	return int(user.Age().Hours() / 24)
}

func (user *User) SortStats() {
	sort.Slice(user.Stats, func(i, j int) bool {
		return user.Stats[i].Mode < user.Stats[j].Mode
	})
}

func (user *User) AvatarUrl() string {
	if user.AvatarHash == nil {
		return fmt.Sprintf("/a/%d", user.Id)
	}
	return fmt.Sprintf("/a/%d?c=%s", user.Id, *user.AvatarHash)
}

func (user *User) UserpageText() string {
	if user.Userpage == nil {
		return ""
	}
	return *user.Userpage
}

// DisplayColor returns the username color derived from the user's groups.
func (user *User) DisplayColor() string {
	for _, group := range user.SortedGroups() {
		if group.Color != "#000000" {
			return group.Color
		}
	}
	return ""
}

// VisibleGroups returns the user's groups excluding
// the donator group & otherwise hidden groups
func (user *User) VisibleGroups() []*Group {
	visible := make([]*Group, 0, len(user.Groups))
	for _, group := range user.SortedGroups() {
		if group.Id == constants.GroupDonator {
			continue
		}
		if group.Hidden {
			continue
		}
		visible = append(visible, group)
	}
	return visible
}

// SortedGroups returns the user's groups in ascending order
func (user *User) SortedGroups() []*Group {
	groups := make([]*Group, 0, len(user.Groups))
	for _, entry := range user.Groups {
		if entry.Group == nil {
			continue
		}
		groups = append(groups, entry.Group)
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Id < groups[j].Id
	})
	return groups
}

// IsDonator returns whether the user is a member of the Donator group.
func (user *User) IsDonator() bool {
	for _, entry := range user.Groups {
		if entry.GroupId == constants.GroupDonator {
			return true
		}
	}
	return false
}

// PreviousNames returns the past usernames of the user.
func (user *User) PreviousNames() []string {
	seen := make(map[string]bool)
	names := make([]string, 0, len(user.Names))
	// hey golang, please add sets, ty

	for _, name := range user.Names {
		if name.Name == user.Name || seen[name.Name] {
			continue
		}
		seen[name.Name] = true
		names = append(names, name.Name)
	}
	return names
}

/* Helper functions for playstyles */

func (user *User) PlaysWithMouse() bool    { return user.Playstyle.Has(constants.PlaystyleMouse) }
func (user *User) PlaysWithKeyboard() bool { return user.Playstyle.Has(constants.PlaystyleKeyboard) }
func (user *User) PlaysWithTablet() bool   { return user.Playstyle.Has(constants.PlaystyleTablet) }
func (user *User) PlaysWithTouch() bool    { return user.Playstyle.Has(constants.PlaystyleTouch) }

type Stats struct {
	UserId      int            `gorm:"column:id;primaryKey"`
	Mode        constants.Mode `gorm:"column:mode;primaryKey"`
	Rank        int            `gorm:"column:rank;default:0"`
	PeakRank    int            `gorm:"column:peak_rank;default:0"`
	Tscore      int64          `gorm:"column:tscore;default:0"`
	Rscore      int64          `gorm:"column:rscore;default:0"`
	PP          float64        `gorm:"column:pp;default:0.0"`
	PPv1        float64        `gorm:"column:ppv1;default:0.0"`
	Playcount   int64          `gorm:"column:playcount;default:0"`
	Playtime    int            `gorm:"column:playtime;default:0"`
	Acc         float64        `gorm:"column:acc;default:0.0"`
	MaxCombo    int            `gorm:"column:max_combo;default:0"`
	TotalHits   int            `gorm:"column:total_hits;default:0"`
	ReplayViews int            `gorm:"column:replay_views;default:0"`
	CountXH     int            `gorm:"column:xh_count;default:0"`
	CountX      int            `gorm:"column:x_count;default:0"`
	CountSH     int            `gorm:"column:sh_count;default:0"`
	CountS      int            `gorm:"column:s_count;default:0"`
	CountA      int            `gorm:"column:a_count;default:0"`
	CountB      int            `gorm:"column:b_count;default:0"`
	CountC      int            `gorm:"column:c_count;default:0"`
	CountD      int            `gorm:"column:d_count;default:0"`

	User *User `gorm:"foreignKey:UserId;references:Id"`
}

func (Stats) TableName() string {
	return "stats"
}

func (stats *Stats) Accuracy() float64 {
	return stats.Acc * 100
}

func (stats *Stats) PlaytimeHours() float64 {
	return float64(stats.Playtime) / 60 / 60
}

func (stats *Stats) Level() int {
	return int(constants.GetLevel(stats.Tscore))
}

// LevelProgress returns the percentage of progress towards the next level
func (stats *Stats) LevelProgress() int {
	precise := constants.GetLevel(stats.Tscore)
	return int(math.Round((precise - math.Floor(precise)) * 100))
}

// LevelBarWidth returns the width (in pixels) of the level progress bar
func (stats *Stats) LevelBarWidth() int {
	return min(stats.LevelProgress(), 100) * 3
}

func (stats *Stats) Clears() int {
	return stats.CountXH +
		stats.CountX +
		stats.CountSH +
		stats.CountS +
		stats.CountA +
		stats.CountB +
		stats.CountC +
		stats.CountD
}

type Relationship struct {
	UserId   int `gorm:"column:user_id;primaryKey"`
	TargetId int `gorm:"column:target_id;primaryKey"`
	Status   int `gorm:"column:status"` // TODO: Add enum for this (0 = friend, 1 = blocked)

	User   *User `gorm:"foreignKey:UserId;references:Id"`
	Target *User `gorm:"foreignKey:TargetId;references:Id"`
}

func (Relationship) TableName() string {
	return "relationships"
}

type Badge struct {
	Id          int       `gorm:"column:id;primaryKey;autoIncrement"`
	UserId      int       `gorm:"column:user_id"`
	Icon        string    `gorm:"column:badge_icon"`
	Url         *string   `gorm:"column:badge_url"`
	Description *string   `gorm:"column:badge_description"`
	Created     time.Time `gorm:"column:created;autoCreateTime"`

	User *User `gorm:"foreignKey:UserId;references:Id"`
}

func (Badge) TableName() string {
	return "profile_badges"
}

type Name struct {
	Id        int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name"`
	UserId    int       `gorm:"column:user_id"`
	Reserved  bool      `gorm:"column:reserved;default:true"`
	ChangedAt time.Time `gorm:"column:changed_at;autoCreateTime"`

	User *User `gorm:"foreignKey:UserId;references:Id"`
}

func (Name) TableName() string {
	return "name_history"
}

type Infringement struct {
	Id          int        `gorm:"column:id;primaryKey;autoIncrement"`
	UserId      int        `gorm:"column:user_id"`
	Time        time.Time  `gorm:"column:time;primaryKey;autoCreateTime"`
	Action      int        `gorm:"column:action;default:0"`
	Length      *time.Time `gorm:"column:length"`
	IsPermanent bool       `gorm:"column:is_permanent;default:false"`
	Description *string    `gorm:"column:description"`

	User *User `gorm:"foreignKey:UserId;references:Id"`
}

func (Infringement) TableName() string {
	return "infringements"
}

type Report struct {
	Id       int       `gorm:"column:id;primaryKey;autoIncrement"`
	TargetId int       `gorm:"column:target_id"`
	SenderId int       `gorm:"column:sender_id"`
	Time     time.Time `gorm:"column:time;autoCreateTime"`
	Reason   *string   `gorm:"column:reason"`
	Resolved bool      `gorm:"column:resolved;default:false"`

	Target *User `gorm:"foreignKey:TargetId;references:Id"`
	Sender *User `gorm:"foreignKey:SenderId;references:Id"`
}

func (Report) TableName() string {
	return "reports"
}

type Verification struct {
	Id     int                        `gorm:"column:id;primaryKey;autoIncrement"`
	Token  string                     `gorm:"column:token"`
	UserId int                        `gorm:"column:user_id"`
	SentAt time.Time                  `gorm:"column:sent_at;autoCreateTime"`
	Type   constants.VerificationType `gorm:"column:type;default:0"`

	User *User `gorm:"foreignKey:UserId;references:Id"`
}

func (Verification) TableName() string {
	return "verifications"
}

func (v *Verification) IsRecent() bool {
	return time.Since(v.SentAt) < 2*time.Minute
}

func (v *Verification) Username() string {
	// smol helper function to make code no messy
	if v.User == nil {
		return ""
	}
	return v.User.Name
}

func (v *Verification) IsActivation() bool {
	return v.Type == constants.VerificationTypeActivation
}

func (v *Verification) IsPassword() bool {
	return v.Type == constants.VerificationTypePassword
}

type Group struct {
	Id                int     `gorm:"column:id;primaryKey;autoIncrement"`
	Name              string  `gorm:"column:name"`
	ShortName         string  `gorm:"column:short_name"`
	Description       *string `gorm:"column:description"`
	Color             string  `gorm:"column:color"`
	BanchoPermissions *int    `gorm:"column:bancho_permissions;default:0"`
	Hidden            bool    `gorm:"column:hidden;default:false"`

	Permissions []*GroupPermission `gorm:"foreignKey:GroupId;references:Id"`
	Entries     []*GroupEntry      `gorm:"foreignKey:GroupId;references:Id"`
}

func (Group) TableName() string {
	return "groups"
}

type GroupEntry struct {
	GroupId int `gorm:"column:group_id;primaryKey"`
	UserId  int `gorm:"column:user_id;primaryKey"`

	Group *Group `gorm:"foreignKey:GroupId;references:Id"`
	User  *User  `gorm:"foreignKey:UserId;references:Id"`
}

func (GroupEntry) TableName() string {
	return "groups_entries"
}

type UserPermission struct {
	Id         int       `gorm:"column:id;primaryKey;autoIncrement"`
	UserId     int       `gorm:"column:user_id;primaryKey"`
	Permission string    `gorm:"column:permission"`
	Rejected   bool      `gorm:"column:rejected;default:false"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoCreateTime"`

	User *User `gorm:"foreignKey:UserId;references:Id"`
}

func (UserPermission) TableName() string {
	return "user_permissions"
}

type GroupPermission struct {
	Id         int       `gorm:"column:id;primaryKey;autoIncrement"`
	GroupId    int       `gorm:"column:group_id;primaryKey"`
	Permission string    `gorm:"column:permission"`
	Rejected   bool      `gorm:"column:rejected;default:false"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoCreateTime"`

	Group *Group `gorm:"foreignKey:GroupId;references:Id"`
}

func (GroupPermission) TableName() string {
	return "group_permissions"
}

type Notification struct {
	Id      int64                      `gorm:"column:id;primaryKey;autoIncrement"`
	UserId  int                        `gorm:"column:user_id;primaryKey"`
	Type    constants.NotificationType `gorm:"column:type"`
	Header  string                     `gorm:"column:header"`
	Content string                     `gorm:"column:content"`
	Link    string                     `gorm:"column:link"`
	Read    bool                       `gorm:"column:read;default:false"`
	Time    time.Time                  `gorm:"column:time;autoCreateTime"`

	User *User `gorm:"foreignKey:UserId;references:Id"`
}

func (Notification) TableName() string {
	return "notifications"
}
