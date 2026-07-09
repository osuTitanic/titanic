package schemas

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/lib/pq"
)

type Release struct {
	Name        string          `gorm:"column:name;primaryKey"`
	Version     int             `gorm:"column:version"`
	Description string          `gorm:"column:description;default:"`
	Category    string          `gorm:"column:category;default:Uncategorized"`
	KnownBugs   *string         `gorm:"column:known_bugs"`
	Supported   bool            `gorm:"column:supported;default:true"`
	Preview     bool            `gorm:"column:preview;default:false"`
	Downloads   pq.StringArray  `gorm:"column:downloads;type:text[];default:'{}'"`
	Screenshots pq.StringArray  `gorm:"column:screenshots;type:text[];default:'{}'"`
	Hashes      json.RawMessage `gorm:"column:hashes;type:jsonb;default:'[]'"`
	CreatedAt   time.Time       `gorm:"column:created_at;autoCreateTime"`
}

func (Release) TableName() string {
	return "releases_titanic"
}

func (r *Release) PrimaryDownloadUrl() string {
	if r == nil || len(r.Downloads) == 0 {
		return ""
	}
	return r.Downloads[0]
}

func (r *Release) PrimaryScreenshotUrl() string {
	if r == nil || len(r.Screenshots) == 0 {
		return ""
	}
	return r.Screenshots[0]
}

func (r *Release) IsDisplayable() bool {
	if r == nil {
		return false
	}
	if !r.Supported {
		return false
	}
	if r.Preview {
		return false
	}
	if len(r.Downloads) == 0 {
		return false
	}
	return true
}

type ModdedRelease struct {
	Name            string    `gorm:"column:name;primaryKey"`
	Description     string    `gorm:"column:description"`
	CreatorId       int       `gorm:"column:creator_id"`
	TopicId         int       `gorm:"column:topic_id"`
	ClientVersion   int       `gorm:"column:client_version"`
	ClientExtension string    `gorm:"column:client_extension"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`

	Creator *User       `gorm:"foreignKey:CreatorId;references:Id"`
	Topic   *ForumTopic `gorm:"foreignKey:TopicId;references:Id"`
}

func (ModdedRelease) TableName() string {
	return "releases_modding"
}

type ModdedReleaseEntries struct {
	Id          int       `gorm:"column:id;primaryKey;autoIncrement"`
	ModName     string    `gorm:"column:mod_name"`
	Version     string    `gorm:"column:version"`
	Stream      string    `gorm:"column:stream"`
	Checksum    string    `gorm:"column:checksum"`
	DownloadUrl *string   `gorm:"column:download_url"`
	UpdateUrl   *string   `gorm:"column:update_url"`
	PostId      *int      `gorm:"column:post_id"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`

	Mod  *ModdedRelease `gorm:"foreignKey:ModName;references:Name"`
	Post *ForumPost     `gorm:"foreignKey:PostId;references:Id"`
}

func (ModdedReleaseEntries) TableName() string {
	return "releases_modding_entries"
}

type ModdedReleaseChangelog struct {
	Id        int       `gorm:"column:id;primaryKey;autoIncrement"`
	EntryId   int       `gorm:"column:entry_id"`
	Text      string    `gorm:"column:text"`
	Type      string    `gorm:"column:type"`
	Branch    string    `gorm:"column:branch"`
	Author    string    `gorm:"column:author"`
	AuthorId  *int      `gorm:"column:author_id"`
	Area      *string   `gorm:"column:area"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`

	Entry      *ModdedReleaseEntries `gorm:"foreignKey:EntryId;references:Id"`
	AuthorUser *User                 `gorm:"foreignKey:AuthorId;references:Id"`
}

func (ModdedReleaseChangelog) TableName() string {
	return "releases_modding_changelog"
}

type ExtraRelease struct {
	Name        string `gorm:"column:name;primaryKey"`
	Description string `gorm:"column:description"`
	Download    string `gorm:"column:download"`
	Filename    string `gorm:"column:filename"`
	Md5         string `gorm:"column:md5"`
}

func (ExtraRelease) TableName() string {
	return "releases_extra"
}

type ReleasesOfficial struct {
	Id         int       `gorm:"column:id;primaryKey;autoIncrement"`
	Version    int       `gorm:"column:version"`
	Stream     string    `gorm:"column:stream"`
	Subversion int       `gorm:"column:subversion"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`

	Files []*ReleaseFiles `gorm:"many2many:releases_official_entries;foreignKey:Id;joinForeignKey:ReleaseId;References:Id;joinReferences:FileId"`
}

func (ReleasesOfficial) TableName() string {
	return "releases_official"
}

type ReleasesOfficialEntries struct {
	ReleaseId int `gorm:"column:release_id;primaryKey"`
	FileId    int `gorm:"column:file_id;primaryKey"`

	Release *ReleasesOfficial `gorm:"foreignKey:ReleaseId;references:Id"`
	File    *ReleaseFiles     `gorm:"foreignKey:FileId;references:Id"`
}

func (ReleasesOfficialEntries) TableName() string {
	return "releases_official_entries"
}

type ReleaseFiles struct {
	Id          int       `gorm:"column:id;primaryKey;autoIncrement"`
	Filename    string    `gorm:"column:filename" json:"filename"`
	FileVersion int       `gorm:"column:file_version" json:"file_version,string"`
	FileHash    string    `gorm:"column:file_hash" json:"file_hash"`
	Filesize    int       `gorm:"column:filesize" json:"filesize,string"`
	PatchId     *string   `gorm:"column:patch_id" json:"patch_id,omitempty"`
	UrlFull     string    `gorm:"column:url_full" json:"url_full"`
	UrlPatch    *string   `gorm:"column:url_patch" json:"url_patch,omitempty"`
	Timestamp   Timestamp `gorm:"column:timestamp;autoCreateTime" json:"timestamp"`

	OfficialReleases []*ReleasesOfficial `gorm:"many2many:releases_official_entries;foreignKey:Id;joinForeignKey:FileId;References:Id;joinReferences:ReleaseId"`
}

func (ReleaseFiles) TableName() string {
	return "releases_official_files"
}

type ReleaseChangelog struct {
	Id        int       `gorm:"column:id;primaryKey;autoIncrement"`
	Text      string    `gorm:"column:text"`
	Type      string    `gorm:"column:type"`
	Branch    string    `gorm:"column:branch"`
	Author    string    `gorm:"column:author"`
	Area      *string   `gorm:"column:area"`
	CreatedAt time.Time `gorm:"column:created_at;default:func.current_date()"`
}

func (ReleaseChangelog) TableName() string {
	return "releases_official_changelog"
}

// PrefixedText returns the changelog text prefixed with its area, if set
func (c *ReleaseChangelog) PrefixedText() string {
	if c.Area == nil {
		return c.Text
	}
	area := strings.TrimSpace(*c.Area)
	if area == "" {
		return c.Text
	}
	return area + ": " + c.Text
}

// TypeSymbol returns the plain-text symbol used in the client changelog stream
func (c *ReleaseChangelog) TypeSymbol() string {
	switch strings.ToLower(c.Type) {
	case "add", "important":
		return "+"
	case "remove":
		return "-"
	case "fix":
		return "*"
	default:
		return ""
	}
}

// IconType returns the changelog type normalized to an available icon name
func (c *ReleaseChangelog) IconType() string {
	switch c.Type {
	case "fix", "add", "misc":
		return c.Type
	default:
		return "misc"
	}
}
