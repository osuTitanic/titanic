package schemas

import (
	"fmt"
	"time"

	"github.com/osuTitanic/titanic-go/internal/constants"
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
	artist := "Unknown Artist"
	title := "Unknown Title"
	if b.Artist != nil {
		artist = *b.Artist
	}
	if b.Title != nil {
		title = *b.Title
	}
	return artist + " - " + title
}

func (b *Beatmapset) Link() string {
	if b.Server == constants.BeatmapServerBancho {
		return fmt.Sprintf("https://osu.ppy.sh/s/%d", b.Id)
	} else {
		return fmt.Sprintf("/s/%d", b.Id)
	}
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
	if b.Beatmapset == nil {
		return fmt.Sprintf("/b/%d", b.Id)
	}
	if b.Beatmapset.Server == constants.BeatmapServerBancho {
		return fmt.Sprintf("https://osu.ppy.sh/b/%d", b.Id)
	}
	return fmt.Sprintf("/b/%d", b.Id)
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
	PostId   int       `gorm:"column:post_id"`
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
