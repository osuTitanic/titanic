package models

type Match interface {
	Id() int32
	InProgress() bool
	Type() uint8
	Mods() uint32
	Name() string
	Password() string
	BeatmapText() string
	BeatmapId() int32
	BeatmapChecksum() string
	HostId() int32
	Mode() uint8
	ScoringType() uint8
	TeamType() uint8
	Freemod() bool
	Seed() int32
	SlotCount() int
	Slot(index int) MatchSlot
}

type MatchSlot interface {
	UserId() int32
	Status() uint8
	Team() uint8
	Mods() uint32
}

type MatchJoin interface {
	MatchId() int32
	Password() string
}
