package models

type BeatmapInfo interface {
	Index() int16
	BeatmapId() int32
	BeatmapSetId() int32
	ThreadId() int32
	RankedStatus() int8
	OsuRank() int8
	TaikoRank() int8
	FruitsRank() int8
	ManiaRank() int8
	Checksum() string
}

type BeatmapInfoReply interface {
	BeatmapCount() int
	Beatmap(index int) BeatmapInfo
}

type BeatmapInfoRequest interface {
	FilenameCount() int
	Filename(index int) string
	IdCount() int
	Id(index int) int32
}
