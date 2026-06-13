package rankings

import (
	"sort"
	"strings"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/redis/go-redis/v9"
)

type CountryRanking struct {
	Name             string
	TotalPerformance float64
	TotalRscore      float64
	TotalTscore      float64
	TotalUsers       int
	AveragePP        float64
}

type CountryRankingSingle struct {
	Name  string
	Score float64
}

func (service *RankingsService) TopCountries(mode constants.Mode) ([]*CountryRanking, error) {
	if service == nil || service.client == nil {
		return nil, ErrRedisClientNotInitialized
	}
	rankings := make([]*CountryRanking, 0, len(constants.CountryCodes))

	for _, code := range constants.CountryCodes {
		if code == "XX" {
			continue
		}
		country := strings.ToLower(code)

		performance, err := service.countryLeaderboardScores(mode, country, "performance")
		if err != nil {
			return nil, err
		}
		if len(performance) == 0 {
			continue
		}

		rscore, err := service.countryLeaderboardScores(mode, country, "rscore")
		if err != nil {
			return nil, err
		}
		if len(rscore) == 0 {
			continue
		}

		tscore, err := service.countryLeaderboardScores(mode, country, "tscore")
		if err != nil {
			return nil, err
		}
		if len(tscore) == 0 {
			continue
		}

		totalPerformance := sumRedisScores(performance)
		totalRscore := sumRedisScores(rscore)
		totalTscore := sumRedisScores(tscore)
		totalUsers := len(performance)

		rankings = append(rankings, &CountryRanking{
			Name:             country,
			TotalPerformance: totalPerformance,
			TotalRscore:      totalRscore,
			TotalTscore:      totalTscore,
			TotalUsers:       totalUsers,
			AveragePP:        totalPerformance / float64(totalUsers),
		})
	}

	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].TotalPerformance > rankings[j].TotalPerformance
	})
	return rankings, nil
}

func (service *RankingsService) TopCountriesForType(mode constants.Mode, rankType string) ([]*CountryRankingSingle, error) {
	if service == nil || service.client == nil {
		return nil, ErrRedisClientNotInitialized
	}
	rankings := make([]*CountryRankingSingle, 0, len(constants.CountryCodes))

	for _, code := range constants.CountryCodes {
		if code == "XX" {
			continue
		}
		country := strings.ToLower(code)

		scores, err := service.countryLeaderboardScores(mode, country, rankType)
		if err != nil {
			return nil, err
		}
		if len(scores) == 0 {
			continue
		}

		totalScore := sumRedisScores(scores)
		rankings = append(rankings, &CountryRankingSingle{
			Name:  country,
			Score: totalScore,
		})
	}

	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})
	return rankings, nil
}

func (service *RankingsService) countryLeaderboardScores(mode constants.Mode, country string, rankType string) ([]redis.Z, error) {
	key := service.RankingKey(mode, rankType, &country)
	query := &redis.ZRangeBy{Max: "+inf", Min: "1"}
	return service.client.ZRevRangeByScoreWithScores(service.ctx, key, query).Result()
}

func sumRedisScores(entries []redis.Z) float64 {
	total := 0.0
	for _, entry := range entries {
		total += entry.Score
	}
	return total
}
