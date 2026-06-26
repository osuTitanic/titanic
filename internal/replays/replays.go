package replays

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/schemas"
)

const (
	// dotNetEpochOffset is the number of seconds between the .NET epoch
	// (0001-01-01) and the Unix epoch (1970-01-01).
	dotNetEpochOffset = 62135596800
	ticksPerSecond    = 10_000_000
)

// Serialize produces a replay file that osu! clients can play back.
func Serialize(score *schemas.Score, replay []byte) []byte {
	mods := score.Mods
	if mods.Has(constants.Nightcore) && !mods.Has(constants.DoubleTime) {
		// Nightcore requires DoubleTime to be present
		mods |= constants.DoubleTime
	}

	buf := &bytes.Buffer{}
	writeU8(buf, uint8(score.Mode))
	writeS32(buf, int32(score.ClientVersion))
	writeString(buf, score.Beatmap.Checksum)
	writeString(buf, score.User.Name)
	writeString(buf, OfflineScoreChecksum(score, mods))
	writeU16(buf, uint16(score.Count300))
	writeU16(buf, uint16(score.Count100))
	writeU16(buf, uint16(score.Count50))
	writeU16(buf, uint16(score.CountGeki))
	writeU16(buf, uint16(score.CountKatu))
	writeU16(buf, uint16(score.CountMiss))
	writeU32(buf, uint32(score.TotalScore))
	writeU16(buf, uint16(score.MaxCombo))
	writeBool(buf, score.Perfect)
	writeU32(buf, uint32(mods))
	writeString(buf, "") // HP Graph
	writeS64(buf, ticks(score.SubmittedAt))
	writeU32(buf, uint32(len(replay)))
	buf.Write(replay)

	if score.ClientVersion >= 20140721 {
		writeU64(buf, uint64(score.Id))
	} else {
		writeU32(buf, uint32(score.Id))
	}

	return buf.Bytes()
}

// OfflineScoreChecksum computes osu!'s replay score checksum.
func OfflineScoreChecksum(score *schemas.Score, mods constants.Mods) string {
	raw := fmt.Sprintf(
		"%dp%do%do%dt%da%sr%de%sy%so%du%s%d%s",
		score.Count100+score.Count300,
		score.Count50,
		score.CountGeki,
		score.CountKatu,
		score.CountMiss,
		score.Beatmap.Checksum,
		score.MaxCombo,
		formattedBool(score.Perfect),
		score.User.Name,
		score.TotalScore,
		string(score.Grade),
		uint32(mods),
		formattedBool(score.Passed()),
	)
	sum := md5.Sum([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func ticks(t time.Time) int64 {
	t = t.UTC()
	seconds := t.Unix() + dotNetEpochOffset
	return seconds*ticksPerSecond + int64(t.Nanosecond())/100
}

func formattedBool(value bool) string {
	if value {
		return "True"
	}
	return "False"
}
