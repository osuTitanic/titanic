package models

type UserInfo interface {
	Id() int32
	Name() string
	UserPresence
	UserStatus
	UserStats
}

type UserPresence interface {
	IsIrc() bool
	Timezone() int8
	Permissions() uint8
	Country() string
	CountryIndex() uint8
	City() string
	Longitude() float32
	Latitude() float32
}

type UserStats interface {
	Rank() int32
	RankedScore() uint64
	TotalScore() uint64
	Accuracy() float64
	Playcount() int32
	PP() uint16
}

type UserStatus interface {
	Action() uint8
	Text() string
	Mods() uint32
	Mode() uint8
	BeatmapChecksum() string
	BeatmapId() int32
}

type UserQuit interface {
	Info() UserInfo
	QuitState() uint8
}

type UserInfoList interface {
	UserCount() int
	User(index int) UserInfo
}
