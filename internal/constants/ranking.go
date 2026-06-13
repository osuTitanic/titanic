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
