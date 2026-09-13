package constants

type ClientStatus int8

const (
	ClientStatusIdle ClientStatus = iota
	ClientStatusAfk
	ClientStatusPlaying
	ClientStatusEditing
	ClientStatusModding
	ClientStatusMultiplayer
	ClientStatusWatching
	ClientStatusUnknown
	ClientStatusTesting
	ClientStatusSubmitting
	ClientStatusPaused
	ClientStatusLobby
	ClientStatusMultiplaying
	ClientStatusOsuDirect
)

func (status ClientStatus) Valid() bool {
	return status >= ClientStatusIdle && status <= ClientStatusOsuDirect
}
