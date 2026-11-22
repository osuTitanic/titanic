package chio

import (
	"io"

	"github.com/osuTitanic/titanic/services/anchor/internal/models"
)

/*
 * NOTE: This chio module will later be replaced with an external library that actually implements everything
 */

type Packet struct {
	Id   uint16
	Data any
}

type Client interface {
	ReadPacket(r io.Reader) (*Packet, error)
	WritePacket(w io.Writer, packetId uint16, data []byte) error
	Writers
}

type Writers interface {
	WriteLoginReply(w io.Writer, reply int32) error
	WriteMessage(w io.Writer, message models.Message) error
	WritePing(w io.Writer) error
	WriteIrcChangeUsername(w io.Writer, oldName, newName string) error
	WriteUserStats(w io.Writer, info models.UserInfo) error
	WriteUserQuit(w io.Writer, quit models.UserQuit) error
	WriteSpectatorJoined(w io.Writer, userId int32) error
	WriteSpectatorLeft(w io.Writer, userId int32) error
	WriteSpectateFrames(w io.Writer, bundle models.ReplayFrameBundle) error
	WriteVersionUpdate(w io.Writer) error
	WriteSpectatorCantSpectate(w io.Writer, userId int32) error
	WriteGetAttention(w io.Writer) error
	WriteAnnouncement(w io.Writer, message string) error
	WriteMatchUpdate(w io.Writer, match models.Match) error
	WriteMatchNew(w io.Writer, match models.Match) error
	WriteMatchDisband(w io.Writer, matchId int32) error
	WriteLobbyJoin(w io.Writer, userId int32) error
	WriteLobbyPart(w io.Writer, userId int32) error
	WriteMatchJoinSuccess(w io.Writer, match models.Match) error
	WriteMatchJoinFail(w io.Writer) error
	WriteFellowSpectatorJoined(w io.Writer, userId int32) error
	WriteFellowSpectatorLeft(w io.Writer, userId int32) error
	WriteMatchStart(w io.Writer, match models.Match) error
	WriteMatchScoreUpdate(w io.Writer, frame models.ScoreFrame) error
	WriteMatchTransferHost(w io.Writer) error
	WriteMatchAllPlayersLoaded(w io.Writer) error
	WriteMatchPlayerFailed(w io.Writer, slotId uint32) error
	WriteMatchComplete(w io.Writer) error
	WriteMatchSkip(w io.Writer) error
	WriteUnauthorized(w io.Writer) error
	WriteChannelJoinSuccess(w io.Writer, channel string) error
	WriteChannelRevoked(w io.Writer, channel string) error
	WriteChannelAvailable(w io.Writer, channel models.Channel) error
	WriteChannelAvailableAutojoin(w io.Writer, channel models.Channel) error
	WriteBeatmapInfoReply(w io.Writer, reply models.BeatmapInfoReply) error
	WriteLoginPermissions(w io.Writer, permissions uint32) error
	WriteFriendsList(w io.Writer, userIds []int32) error
	WriteProtocolNegotiation(w io.Writer, version int32) error
	WriteTitleUpdate(w io.Writer, update models.TitleUpdate) error
	WriteMonitor(w io.Writer) error
	WriteMatchPlayerSkipped(w io.Writer, slotId int32) error
	WriteUserPresence(w io.Writer, info models.UserInfo) error
	WriteRestart(w io.Writer, retryMs int32) error
	WriteInvite(w io.Writer, message models.Message) error
	WriteChannelInfoComplete(w io.Writer) error
	WriteMatchChangePassword(w io.Writer, password string) error
	WriteSilenceInfo(w io.Writer, timeRemaining int32) error
	WriteUserSilenced(w io.Writer, userId uint32) error
	WriteUserPresenceSingle(w io.Writer, info models.UserInfo) error
	WriteUserPresenceBundle(w io.Writer, infos []models.UserInfo) error
	WriteUserDMsBlocked(w io.Writer, targetName string) error
	WriteTargetIsSilenced(w io.Writer, targetName string) error
	WriteVersionUpdateForced(w io.Writer) error
	WriteSwitchServer(w io.Writer, target int32) error
	WriteAccountRestricted(w io.Writer) error
	WriteRTX(w io.Writer, message string) error
	WriteMatchAbort(w io.Writer) error
	WriteSwitchTournamentServer(w io.Writer, ip string) error
}
