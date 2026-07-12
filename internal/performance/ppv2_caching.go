package performance

type PPv2CacheLayer interface {
	ToCache(beatmapId int, data any) bool
	FromCache(beatmapId int) (bool, any)
}
