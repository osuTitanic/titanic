package schemas

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"
)

type Screenshot struct {
	Id        int       `gorm:"column:id;primaryKey;autoIncrement"`
	UserId    int       `gorm:"column:user_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	Hidden    bool      `gorm:"column:hidden;default:false"`

	User *User `gorm:"foreignKey:UserId;references:Id"`
}

func (Screenshot) TableName() string {
	return "screenshots"
}

func (s *Screenshot) Checksum() string {
	dateString := s.CreatedAt.Format("2006-01-02 15:04:05")
	checksum := md5.Sum([]byte(dateString))
	return hex.EncodeToString(checksum[:])
}

type Benchmark struct {
	Id         int             `gorm:"column:id;primaryKey;autoIncrement"`
	UserId     int             `gorm:"column:user_id"`
	Smoothness float64         `gorm:"column:smoothness"`
	Framerate  int             `gorm:"column:framerate"`
	Score      int64           `gorm:"column:score"`
	Grade      string          `gorm:"column:grade;default:N"`
	CreatedAt  time.Time       `gorm:"column:created_at;autoCreateTime"`
	Client     string          `gorm:"column:client"`
	Hardware   json.RawMessage `gorm:"column:hardware;type:jsonb"`

	User *User `gorm:"foreignKey:UserId;references:Id"`
}

func (Benchmark) TableName() string {
	return "benchmarks"
}

type Log struct {
	Id      int       `gorm:"column:id;primaryKey;autoIncrement"`
	Level   string    `gorm:"column:level"`
	Type    string    `gorm:"column:type"`
	Message string    `gorm:"column:message"`
	Time    time.Time `gorm:"column:time;autoCreateTime"`
}

func (Log) TableName() string {
	return "logs"
}
