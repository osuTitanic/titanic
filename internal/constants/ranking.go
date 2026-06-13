package constants

import "fmt"

type RankingType string

const (
	RankingTypePerformance RankingType = "performance"
	RankingTypeCountry     RankingType = "country"
	RankingTypeTotalScore  RankingType = "tscore"
	RankingTypeRankedScore RankingType = "rscore"
	RankingTypePPv1        RankingType = "ppv1"
	RankingTypeClears      RankingType = "clears"
	RankingTypeLeader      RankingType = "leader"
)

func (rankingType RankingType) Alias() string {
	return string(rankingType)
}

func (rankingType RankingType) String() string {
	switch rankingType {
	case RankingTypePerformance:
		return "Performance"
	case RankingTypeCountry:
		return "Country"
	case RankingTypeTotalScore:
		return "Total Score"
	case RankingTypeRankedScore:
		return "Ranked Score"
	case RankingTypePPv1:
		return "PPv1"
	case RankingTypeClears:
		return "Clears"
	case RankingTypeLeader:
		return "First Places"
	default:
		return fmt.Sprintf("Unknown (%s)", string(rankingType))
	}
}

func NewRankingTypeFromAlias(alias string) (RankingType, bool) {
	switch alias {
	case "global":
		return RankingTypePerformance, true
	case "performance":
		return RankingTypePerformance, true
	case "country":
		return RankingTypeCountry, true
	case "tscore":
		return RankingTypeTotalScore, true
	case "rscore":
		return RankingTypeRankedScore, true
	case "ppv1":
		return RankingTypePPv1, true
	case "clears":
		return RankingTypeClears, true
	case "leader":
		return RankingTypeLeader, true
	default:
		return "", false
	}
}
