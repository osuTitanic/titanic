package replays

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
)

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

func formattedBool(value bool) string {
	if value {
		return "True"
	}
	return "False"
}
