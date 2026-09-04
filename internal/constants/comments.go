package constants

type CommentTarget string

const (
	CommentTargetMap    CommentTarget = "map"
	CommentTargetReplay CommentTarget = "replay"
	CommentTargetSong   CommentTarget = "song"
)

func (target CommentTarget) Valid() bool {
	switch target {
	case CommentTargetMap, CommentTargetReplay, CommentTargetSong:
		return true
	default:
		return false
	}
}

func (target CommentTarget) String() string {
	switch target {
	case CommentTargetMap:
		return "Map"
	case CommentTargetReplay:
		return "Replay"
	case CommentTargetSong:
		return "Song"
	default:
		return "Unknown"
	}
}
