package schemas

import (
	"encoding/json"
	"time"

	"github.com/osuTitanic/titanic-go/internal/constants"
)

type Login struct {
	UserId  int       `gorm:"column:user_id;primaryKey"`
	Time    time.Time `gorm:"column:time;primaryKey;autoCreateTime"`
	Ip      string    `gorm:"column:ip"`
	Version string    `gorm:"column:osu_version"`

	User *User `gorm:"foreignKey:UserId;references:Id"`
}

func (Login) TableName() string {
	return "logins"
}

type Channel struct {
	Name             string `gorm:"column:name;primaryKey"`
	Topic            string `gorm:"column:topic"`
	ReadPermissions  int    `gorm:"column:read_permissions;default:1"`
	WritePermissions int    `gorm:"column:write_permissions;default:1"`
}

func (Channel) TableName() string {
	return "channels"
}

type Message struct {
	Id      int       `gorm:"column:id;primaryKey;autoIncrement"`
	Sender  string    `gorm:"column:sender"`
	Target  string    `gorm:"column:target"`
	Message string    `gorm:"column:message"`
	Time    time.Time `gorm:"column:time;autoCreateTime"`
}

func (Message) TableName() string {
	return "messages"
}

type DirectMessage struct {
	Id       int       `gorm:"column:id;primaryKey;autoIncrement"`
	SenderId int       `gorm:"column:sender_id"`
	TargetId int       `gorm:"column:target_id"`
	Message  string    `gorm:"column:message"`
	Time     time.Time `gorm:"column:time;autoCreateTime"`
	Read     bool      `gorm:"column:read;default:false"`

	Sender *User `gorm:"foreignKey:SenderId;references:Id"`
	Target *User `gorm:"foreignKey:TargetId;references:Id"`
}

func (DirectMessage) TableName() string {
	return "direct_messages"
}

type ChatFilter struct {
	Name                 string    `gorm:"column:name;primaryKey"`
	Pattern              string    `gorm:"column:pattern"`
	Replacement          *string   `gorm:"column:replacement"`
	Block                bool      `gorm:"column:block;default:false"`
	BlockTimeoutDuration *int      `gorm:"column:block_timeout_duration"`
	CreatedAt            time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (ChatFilter) TableName() string {
	return "filters"
}

type Activity struct {
	Id     int             `gorm:"column:id;primaryKey;autoIncrement"`
	UserId int             `gorm:"column:user_id"`
	Time   time.Time       `gorm:"column:time;autoCreateTime"`
	Mode   *constants.Mode `gorm:"column:mode"`
	Type   int             `gorm:"column:type;default:0"`
	Data   json.RawMessage `gorm:"column:data;type:jsonb;default:'{}'"`
	Hidden bool            `gorm:"column:hidden;default:false"`

	User *User `gorm:"foreignKey:UserId;references:Id"`
}

func (Activity) TableName() string {
	return "profile_activity"
}

type HardwareInfo struct {
	UserId        int    `gorm:"column:user_id;primaryKey"`
	Executable    string `gorm:"column:executable;primaryKey"`
	Adapters      string `gorm:"column:adapters;primaryKey"`
	UniqueId      string `gorm:"column:unique_id;primaryKey"`
	DiskSignature string `gorm:"column:disk_signature;primaryKey"`
	Banned        bool   `gorm:"column:banned;default:false"`

	User *User `gorm:"foreignKey:UserId;references:Id"`
}

func (HardwareInfo) TableName() string {
	return "clients"
}

type HardwareVerified struct {
	Type int    `gorm:"column:type;primaryKey"`
	Hash string `gorm:"column:hash;primaryKey"`
}

func (HardwareVerified) TableName() string {
	return "clients_verified"
}

type Match struct {
	Id        int        `gorm:"column:id;primaryKey;autoIncrement"`
	BanchoId  int        `gorm:"column:bancho_id"`
	Name      string     `gorm:"column:name"`
	CreatorId int        `gorm:"column:creator_id"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	EndedAt   *time.Time `gorm:"column:ended_at"`

	Creator *User         `gorm:"foreignKey:CreatorId;references:Id"`
	Events  []*MatchEvent `gorm:"foreignKey:MatchId;references:Id"`
}

func (Match) TableName() string {
	return "mp_matches"
}

type MatchEvent struct {
	MatchId int             `gorm:"column:match_id;primaryKey"`
	Time    time.Time       `gorm:"column:time;primaryKey;autoCreateTime"`
	Type    int             `gorm:"column:type"`
	Data    json.RawMessage `gorm:"column:data;type:jsonb"`

	Match *Match `gorm:"foreignKey:MatchId;references:Id"`
}

func (MatchEvent) TableName() string {
	return "mp_events"
}

type BanchoActivity struct {
	Time     time.Time `gorm:"column:time;primaryKey;autoCreateTime"`
	OsuCount int       `gorm:"column:osu_count;default:0"`
	IrcCount int       `gorm:"column:irc_count;default:0"`
	MpCount  int       `gorm:"column:mp_count;default:0"`
}

func (BanchoActivity) TableName() string {
	return "user_activity"
}
