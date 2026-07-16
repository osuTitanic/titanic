package schemas

import (
	"fmt"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
)

type Beatmapset struct {
	Id                 int                       `gorm:"column:id;primaryKey;autoIncrement"`
	Title              *string                   `gorm:"column:title"`
	TitleUnicode       *string                   `gorm:"column:title_unicode"`
	Artist             *string                   `gorm:"column:artist"`
	ArtistUnicode      *string                   `gorm:"column:artist_unicode"`
	Source             *string                   `gorm:"column:source"`
	SourceUnicode      *string                   `gorm:"column:source_unicode"`
	Creator            *string                   `gorm:"column:creator"`
	DisplayTitle       *string                   `gorm:"column:display_title"`
	Description        *string                   `gorm:"column:description"`
	Tags               *string                   `gorm:"column:tags;default:"`
	Status             constants.BeatmapStatus   `gorm:"column:submission_status;default:3"`
	HasVideo           bool                      `gorm:"column:has_video;default:false"`
	HasStoryboard      bool                      `gorm:"column:has_storyboard;default:false"`
	Server             constants.BeatmapServer   `gorm:"column:server;default:0"`
	DownloadServer     constants.BeatmapServer   `gorm:"column:download_server;default:0"`
	TopicId            *int                      `gorm:"column:topic_id"`
	CreatorId          *int                      `gorm:"column:creator_id"`
	Available          bool                      `gorm:"column:available;default:true"`
	Enhanced           bool                      `gorm:"column:enhanced;default:false"`
	Explicit           bool                      `gorm:"column:explicit;default:false"`
	CreatedAt          time.Time                 `gorm:"column:submission_date;autoCreateTime"`
	ApprovedAt         *time.Time                `gorm:"column:approved_date"`
	ApprovedBy         *int                      `gorm:"column:approved_by"`
	LastUpdate         time.Time                 `gorm:"column:last_updated;autoCreateTime"`
	AddedAt            *time.Time                `gorm:"column:added_at;autoCreateTime"`
	TotalPlaycount     int64                     `gorm:"column:total_playcount;default:0;->"`
	MaxDiff            float64                   `gorm:"column:max_diff;default:0.0;->"`
	RatingAverage      float64                   `gorm:"column:rating_average;default:0.0;->"`
	RatingCount        int                       `gorm:"column:rating_count;default:0;->"`
	FavouriteCount     int                       `gorm:"column:favourite_count;default:0;->"`
	OszFilesize        int                       `gorm:"column:osz_filesize;default:0"`
	OszFilesizeNovideo int                       `gorm:"column:osz_filesize_novideo;default:0"`
	LanguageId         constants.BeatmapLanguage `gorm:"column:language_id;default:1"`
	GenreId            constants.BeatmapGenre    `gorm:"column:genre_id;default:1"`
	StarPriority       int                       `gorm:"column:star_priority;default:0"`
	Offset             int                       `gorm:"column:offset;default:0"`
	MetaHash           *string                   `gorm:"column:meta_hash"`
	InfoHash           *string                   `gorm:"column:info_hash"`
	BodyHash           *string                   `gorm:"column:body_hash"`
	Search             string                    `gorm:"column:search;type:tsvector;->"`

	CreatorUser    *User                `gorm:"foreignKey:CreatorId;references:Id"`
	ApprovedByUser *User                `gorm:"foreignKey:ApprovedBy;references:Id"`
	Beatmaps       []*Beatmap           `gorm:"foreignKey:SetId;references:Id"`
	Nominations    []*BeatmapNomination `gorm:"foreignKey:SetId;references:Id"`
	Modding        []*BeatmapModding    `gorm:"foreignKey:SetId;references:Id"`
	Favourites     []*BeatmapFavourite  `gorm:"foreignKey:SetId;references:Id"`
	Ratings        []*BeatmapRating     `gorm:"foreignKey:SetId;references:Id"`
	Plays          []*BeatmapPlays      `gorm:"foreignKey:SetId;references:Id"`
	PackEntries    []*BeatmapPackEntry  `gorm:"foreignKey:BeatmapsetId;references:Id"`
}

func (Beatmapset) TableName() string {
	return "beatmapsets"
}

func (b *Beatmapset) Name() string {
	return b.ArtistName() + " - " + b.TitleName()
}

func (b *Beatmapset) ArtistName() string {
	if b.Artist != nil {
		return *b.Artist
	}
	return "Unknown Artist"
}

func (b *Beatmapset) TitleName() string {
	if b.Title != nil {
		return *b.Title
	}
	return "Unknown Title"
}

func (b *Beatmapset) SourceName() string {
	if b.Source != nil {
		return *b.Source
	}
	return ""
}

func (b *Beatmapset) Link() string {
	return fmt.Sprintf("/s/%d", b.Id)
}

func (b *Beatmapset) CommentLink() string {
	if b.Server == constants.BeatmapServerBancho {
		return fmt.Sprintf("https://osu.ppy.sh/beatmapsets/%d#comments", b.Id)
	}
	if b.TopicId != nil {
		return fmt.Sprintf("/forum/t/%d", *b.TopicId)
	}
	return "/forum/t/0"
}

func (b *Beatmapset) CreatorName() string {
	if b.Creator != nil {
		return *b.Creator
	}
	return "Unknown"
}

func (b *Beatmapset) CreatorLink() (creatorLink string) {
	if b.CreatorId != nil {
		creatorLink = fmt.Sprintf("/u/%d", *b.CreatorId)
	} else if b.Creator != nil {
		creatorLink = fmt.Sprintf("/u/%s", *b.Creator)
	} else {
		creatorLink = "/u/0" // :ehbounce:
	}

	if b.Server == constants.BeatmapServerBancho {
		creatorLink = "https://osu.ppy.sh" + creatorLink
	}
	return creatorLink
}

func (b *Beatmapset) ThumbnailUrl() string {
	lastModified := b.LastUpdate.Unix()
	return fmt.Sprintf("/mt/%d?c=%d", b.Id, lastModified)
}

func (b *Beatmapset) LargeThumbnailUrl() string {
	lastModified := b.LastUpdate.Unix()
	return fmt.Sprintf("/mt/%dl?c=%d", b.Id, lastModified)
}

func (b *Beatmapset) AudioPreviewUrl() string {
	lastModified := b.LastUpdate.Unix()
	return fmt.Sprintf("/mp3/preview/%d?c=%d", b.Id, lastModified)
}

func (b *Beatmapset) DescriptionText() string {
	if b.Description != nil {
		return *b.Description
	}
	return ""
}

func (b *Beatmapset) DisplayTitleText() string {
	if b.DisplayTitle != nil {
		return *b.DisplayTitle
	}
	return ""
}

func (b *Beatmapset) TagsText() string {
	if b.Tags != nil {
		return *b.Tags
	}
	return ""
}

func (b *Beatmapset) IsApproved() bool {
	return b.Status > constants.BeatmapStatusPending
}

// RequiredNominations returns the amount of nominations a beatmapset needs to
// be qualified: 2, plus one extra for every additional game mode it contains.
func (b *Beatmapset) RequiredNominations() int {
	modes := make(map[constants.Mode]struct{}, len(b.Beatmaps))
	for _, beatmap := range b.Beatmaps {
		modes[beatmap.Mode] = struct{}{}
	}
	// i wish go had sets ... :(

	return 2 + max(len(modes)-1, 0)
}

// RankDate returns the date a qualified beatmapset is expected to be ranked,
// i.e. seven days after it was approved.
func (b *Beatmapset) RankDate() time.Time {
	if b.ApprovedAt != nil {
		return b.ApprovedAt.Add(7 * 24 * time.Hour)
	}
	return b.LastUpdate.Add(7 * 24 * time.Hour)
}

// RankingStatus returns the status a qualified beatmapset will receive after
// its qualification period.
// Beatmapsets with a five-minute drain time are approved instead of ranked.
func (b *Beatmapset) RankingStatus() constants.BeatmapStatus {
	for _, beatmap := range b.Beatmaps {
		if beatmap.DrainLength >= 5*60 {
			return constants.BeatmapStatusApproved
		}
	}
	return constants.BeatmapStatusRanked
}

func (b *Beatmapset) DisplayDate() time.Time {
	if b.ApprovedAt != nil {
		return *b.ApprovedAt
	}
	return b.LastUpdate
}

func (b *Beatmapset) DisplayDateTitle() string {
	if b.ApprovedAt != nil {
		return "Approved date"
	}
	return "Last update"
}

type Beatmap struct {
	Id               int                     `gorm:"column:id;primaryKey;autoIncrement"`
	SetId            int                     `gorm:"column:set_id"`
	Mode             constants.Mode          `gorm:"column:mode;default:0"`
	Status           constants.BeatmapStatus `gorm:"column:status;default:2"`
	Checksum         string                  `gorm:"column:md5"`
	Version          string                  `gorm:"column:version"`
	Filename         string                  `gorm:"column:filename"`
	CreatedAt        time.Time               `gorm:"column:submission_date;autoCreateTime"`
	LastUpdate       time.Time               `gorm:"column:last_updated;autoCreateTime"`
	Playcount        int64                   `gorm:"column:playcount;default:0"`
	Passcount        int64                   `gorm:"column:passcount;default:0"`
	TotalLength      int                     `gorm:"column:total_length"`
	DrainLength      int                     `gorm:"column:drain_length;default:0"`
	CountNormal      int                     `gorm:"column:count_normal;default:0"`
	CountSlider      int                     `gorm:"column:count_slider;default:0"`
	CountSpinner     int                     `gorm:"column:count_spinner;default:0"`
	MaxCombo         int                     `gorm:"column:max_combo"`
	BPM              float64                 `gorm:"column:bpm;default:0.0"`
	CS               float64                 `gorm:"column:cs;default:0.0"`
	AR               float64                 `gorm:"column:ar;default:0.0"`
	OD               float64                 `gorm:"column:od;default:0.0"`
	HP               float64                 `gorm:"column:hp;default:0.0"`
	Diff             float64                 `gorm:"column:diff;default:0.0"`
	DiffEyup         float64                 `gorm:"column:diff_eyup;default:0.0"`
	SliderMultiplier float64                 `gorm:"column:slider_multiplier;default:0.0"`
	Search           string                  `gorm:"column:search;type:tsvector;->"`

	Beatmapset            *Beatmapset                    `gorm:"foreignKey:SetId;references:Id"`
	CollaborationRequests []*BeatmapCollaborationRequest `gorm:"foreignKey:BeatmapId;references:Id"`
	Collaborations        []*BeatmapCollaboration        `gorm:"foreignKey:BeatmapId;references:Id"`
	Ratings               []*BeatmapRating               `gorm:"foreignKey:MapChecksum;references:Checksum"`
	Plays                 []*BeatmapPlays                `gorm:"foreignKey:BeatmapId;references:Id"`
}

func (Beatmap) TableName() string {
	return "beatmaps"
}

func (b *Beatmap) Name() string {
	if b.Beatmapset != nil {
		return b.Beatmapset.Name() + " [" + b.Version + "]"
	}
	return b.Version
}

func (b *Beatmap) Link() string {
	return fmt.Sprintf("/b/%d", b.Id)
}

func (b *Beatmap) DifficultyAlias() string {
	difficulty := "expert"
	if b.Diff < 2 {
		difficulty = "easy"
	} else if b.Diff < 2.7 {
		difficulty = "normal"
	} else if b.Diff < 4 {
		difficulty = "hard"
	} else if b.Diff < 5.3 {
		difficulty = "insane"
	}
	return difficulty
}

type BeatmapCollaboration struct {
	UserId               int       `gorm:"column:user_id;primaryKey"`
	BeatmapId            int       `gorm:"column:beatmap_id;primaryKey"`
	IsBeatmapAuthor      bool      `gorm:"column:is_beatmap_author;default:false"`
	AllowResourceUpdates bool      `gorm:"column:allow_resource_updates;default:false"`
	CreatedAt            time.Time `gorm:"column:created_at;autoCreateTime"`

	User    *User    `gorm:"foreignKey:UserId;references:Id"`
	Beatmap *Beatmap `gorm:"foreignKey:BeatmapId;references:Id"`
}

func (BeatmapCollaboration) TableName() string {
	return "beatmap_collaboration"
}

type BeatmapCollaborationRequest struct {
	Id                   int       `gorm:"column:id;primaryKey;autoIncrement"`
	UserId               int       `gorm:"column:user_id"`
	TargetId             int       `gorm:"column:target_id"`
	BeatmapId            int       `gorm:"column:beatmap_id"`
	AllowResourceUpdates bool      `gorm:"column:allow_resource_updates;default:false"`
	CreatedAt            time.Time `gorm:"column:created_at;autoCreateTime"`

	User    *User    `gorm:"foreignKey:UserId;references:Id"`
	Target  *User    `gorm:"foreignKey:TargetId;references:Id"`
	Beatmap *Beatmap `gorm:"foreignKey:BeatmapId;references:Id"`
}

func (BeatmapCollaborationRequest) TableName() string {
	return "beatmap_collaboration_requests"
}

type BeatmapCollaborationBlacklist struct {
	UserId    int       `gorm:"column:user_id;primaryKey"`
	TargetId  int       `gorm:"column:target_id;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`

	User   *User `gorm:"foreignKey:UserId;references:Id"`
	Target *User `gorm:"foreignKey:TargetId;references:Id"`
}

func (BeatmapCollaborationBlacklist) TableName() string {
	return "beatmap_collaboration_blacklist"
}

type BeatmapNomination struct {
	UserId int       `gorm:"column:user_id;primaryKey"`
	SetId  int       `gorm:"column:set_id;primaryKey"`
	Time   time.Time `gorm:"column:time;autoCreateTime"`

	User       *User       `gorm:"foreignKey:UserId;references:Id"`
	Beatmapset *Beatmapset `gorm:"foreignKey:SetId;references:Id"`
}

func (BeatmapNomination) TableName() string {
	return "beatmap_nominations"
}

type BeatmapModding struct {
	Id       int       `gorm:"column:id;primaryKey;autoIncrement"`
	TargetId int       `gorm:"column:target_id"`
	SenderId int       `gorm:"column:sender_id"`
	SetId    int       `gorm:"column:set_id"`
	PostId   int64     `gorm:"column:post_id"`
	Amount   int       `gorm:"column:amount;default:0"`
	Time     time.Time `gorm:"column:time;autoCreateTime"`

	Beatmapset *Beatmapset `gorm:"foreignKey:SetId;references:Id"`
	Post       *ForumPost  `gorm:"foreignKey:PostId;references:Id"`
	Target     *User       `gorm:"foreignKey:TargetId;references:Id"`
	Sender     *User       `gorm:"foreignKey:SenderId;references:Id"`
}

func (BeatmapModding) TableName() string {
	return "beatmap_modding"
}

type BeatmapPack struct {
	Id           int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name         string    `gorm:"column:name"`
	Category     string    `gorm:"column:category"`
	DownloadLink string    `gorm:"column:download_link"`
	Description  string    `gorm:"column:description;default:"`
	CreatorId    int       `gorm:"column:creator_id"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoCreateTime"`

	Entries []*BeatmapPackEntry `gorm:"foreignKey:PackId;references:Id"`
	Creator *User               `gorm:"foreignKey:CreatorId;references:Id"`
}

func (BeatmapPack) TableName() string {
	return "beatmap_packs"
}

type BeatmapPackEntry struct {
	PackId       int       `gorm:"column:pack_id;primaryKey"`
	BeatmapsetId int       `gorm:"column:beatmapset_id;primaryKey"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`

	Pack       *BeatmapPack `gorm:"foreignKey:PackId;references:Id"`
	Beatmapset *Beatmapset  `gorm:"foreignKey:BeatmapsetId;references:Id"`
}

func (BeatmapPackEntry) TableName() string {
	return "beatmap_pack_entries"
}

type BeatmapPlays struct {
	UserId      int    `gorm:"column:user_id;primaryKey"`
	BeatmapId   int    `gorm:"column:beatmap_id;primaryKey"`
	SetId       int    `gorm:"column:set_id"`
	Count       int    `gorm:"column:count"`
	BeatmapFile string `gorm:"column:beatmap_file"`

	User       *User       `gorm:"foreignKey:UserId;references:Id"`
	Beatmap    *Beatmap    `gorm:"foreignKey:BeatmapId;references:Id"`
	Beatmapset *Beatmapset `gorm:"foreignKey:SetId;references:Id"`
}

func (BeatmapPlays) TableName() string {
	return "plays"
}

type BeatmapFavourite struct {
	UserId    int       `gorm:"column:user_id;primaryKey"`
	SetId     int       `gorm:"column:set_id;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`

	User       *User       `gorm:"foreignKey:UserId;references:Id"`
	Beatmapset *Beatmapset `gorm:"foreignKey:SetId;references:Id"`
}

func (BeatmapFavourite) TableName() string {
	return "favourites"
}

type BeatmapRating struct {
	UserId      int    `gorm:"column:user_id;primaryKey"`
	SetId       int    `gorm:"column:set_id"`
	MapChecksum string `gorm:"column:map_checksum;primaryKey"`
	Rating      int    `gorm:"column:rating"`

	User       *User       `gorm:"foreignKey:UserId;references:Id"`
	Beatmap    *Beatmap    `gorm:"foreignKey:MapChecksum;references:Checksum"`
	Beatmapset *Beatmapset `gorm:"foreignKey:SetId;references:Id"`
}

func (BeatmapRating) TableName() string {
	return "ratings"
}

type BeatmapComment struct {
	Id         int            `gorm:"column:id;primaryKey;autoIncrement"`
	TargetId   int            `gorm:"column:target_id"`
	TargetType string         `gorm:"column:target_type"`
	UserId     int            `gorm:"column:user_id"`
	Mode       constants.Mode `gorm:"column:mode;default:0"`
	Time       time.Time      `gorm:"column:time;autoCreateTime"`
	Comment    string         `gorm:"column:comment"`
	Format     *string        `gorm:"column:format"`
	Color      *string        `gorm:"column:color"`

	User *User `gorm:"foreignKey:UserId;references:Id"`
}

func (BeatmapComment) TableName() string {
	return "comments"
}

type BeatmapMirror struct {
	Url      string                        `gorm:"column:url;primaryKey"`
	Server   constants.BeatmapServer       `gorm:"column:server"`
	Type     constants.BeatmapResourceType `gorm:"column:type"`
	Priority int                           `gorm:"column:priority;default:0"`
}

func (BeatmapMirror) TableName() string {
	return "resource_mirrors"
}
