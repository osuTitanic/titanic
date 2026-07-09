package rankings

import (
	"sort"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
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

	// Queue every country's leaderboards onto a single pipeline, so the
	// whole listing only costs one round-trip instead of one per country.
	pipe := service.client.Pipeline()
	countries := make([]string, 0, len(constants.CountryCodes))
	performanceCmds := make([]*redis.ZSliceCmd, 0, len(constants.CountryCodes))
	rscoreCmds := make([]*redis.ZSliceCmd, 0, len(constants.CountryCodes))
	tscoreCmds := make([]*redis.ZSliceCmd, 0, len(constants.CountryCodes))

	for _, code := range constants.CountryCodes {
		if code == "XX" {
			continue
		}
		country := strings.ToLower(code)
		countries = append(countries, country)
		performanceCmds = append(performanceCmds, service.queueCountryLeaderboardScores(pipe, mode, country, "performance"))
		rscoreCmds = append(rscoreCmds, service.queueCountryLeaderboardScores(pipe, mode, country, "rscore"))
		tscoreCmds = append(tscoreCmds, service.queueCountryLeaderboardScores(pipe, mode, country, "tscore"))
	}

	if _, err := pipe.Exec(service.ctx); err != nil {
		return nil, err
	}

	rankings := make([]*CountryRanking, 0, len(countries))
	for i, country := range countries {
		performance, err := performanceCmds[i].Result()
		if err != nil {
			return nil, err
		}
		if len(performance) == 0 {
			continue
		}

		rscore, err := rscoreCmds[i].Result()
		if err != nil {
			return nil, err
		}
		if len(rscore) == 0 {
			continue
		}

		tscore, err := tscoreCmds[i].Result()
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

	// Queue every country's leaderboard onto a single pipeline, so the
	// whole listing only costs one round-trip
	pipe := service.client.Pipeline()
	countries := make([]string, 0, len(constants.CountryCodes))
	commands := make([]*redis.ZSliceCmd, 0, len(constants.CountryCodes))

	for _, code := range constants.CountryCodes {
		if code == "XX" {
			continue
		}
		country := strings.ToLower(code)
		countries = append(countries, country)
		commands = append(commands, service.queueCountryLeaderboardScores(pipe, mode, country, rankType))
	}

	if _, err := pipe.Exec(service.ctx); err != nil {
		return nil, err
	}

	rankings := make([]*CountryRankingSingle, 0, len(countries))
	for i, country := range countries {
		scores, err := commands[i].Result()
		if err != nil {
			return nil, err
		}
		if len(scores) == 0 {
			continue
		}

		rankings = append(rankings, &CountryRankingSingle{
			Name:  country,
			Score: sumRedisScores(scores),
		})
	}

	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})
	return rankings, nil
}

func (service *RankingsService) queueCountryLeaderboardScores(pipe redis.Pipeliner, mode constants.Mode, country string, rankType string) *redis.ZSliceCmd {
	key := service.RankingKey(mode, rankType, &country)
	query := &redis.ZRangeBy{Max: "+inf", Min: "1"}
	return pipe.ZRevRangeByScoreWithScores(service.ctx, key, query)
}

func sumRedisScores(entries []redis.Z) float64 {
	total := 0.0
	for _, entry := range entries {
		total += entry.Score
	}
	return total
}
