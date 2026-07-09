package rankings

import (
	"fmt"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/redis/go-redis/v9"
)

func (service *RankingsService) Update(stats *schemas.Stats, country string) error {
	if service == nil || service.client == nil {
		return ErrRedisClientNotInitialized
	}
	if stats == nil {
		return ErrStatsIsNil
	}

	country = strings.ToLower(country)
	clears := stats.Clears()
	mode := stats.Mode

	pipe := service.client.Pipeline()
	pipe.ZAdd(service.ctx, fmt.Sprintf("bancho:performance:%d", mode), redis.Z{Score: stats.PP, Member: stats.UserId})
	pipe.ZAdd(service.ctx, fmt.Sprintf("bancho:rscore:%d", mode), redis.Z{Score: float64(stats.Rscore), Member: stats.UserId})
	pipe.ZAdd(service.ctx, fmt.Sprintf("bancho:tscore:%d", mode), redis.Z{Score: float64(stats.Tscore), Member: stats.UserId})
	pipe.ZAdd(service.ctx, fmt.Sprintf("bancho:ppv1:%d", mode), redis.Z{Score: stats.PPv1, Member: stats.UserId})
	pipe.ZAdd(service.ctx, fmt.Sprintf("bancho:acc:%d", mode), redis.Z{Score: stats.Acc, Member: stats.UserId})
	pipe.ZAdd(service.ctx, fmt.Sprintf("bancho:clears:%d", mode), redis.Z{Score: float64(clears), Member: stats.UserId})

	if len(country) == 2 && country != "xx" {
		pipe.ZAdd(service.ctx, fmt.Sprintf("bancho:performance:%d:%s", mode, country), redis.Z{Score: stats.PP, Member: stats.UserId})
		pipe.ZAdd(service.ctx, fmt.Sprintf("bancho:rscore:%d:%s", mode, country), redis.Z{Score: float64(stats.Rscore), Member: stats.UserId})
		pipe.ZAdd(service.ctx, fmt.Sprintf("bancho:tscore:%d:%s", mode, country), redis.Z{Score: float64(stats.Tscore), Member: stats.UserId})
		pipe.ZAdd(service.ctx, fmt.Sprintf("bancho:ppv1:%d:%s", mode, country), redis.Z{Score: stats.PPv1, Member: stats.UserId})
		pipe.ZAdd(service.ctx, fmt.Sprintf("bancho:acc:%d:%s", mode, country), redis.Z{Score: stats.Acc, Member: stats.UserId})
		pipe.ZAdd(service.ctx, fmt.Sprintf("bancho:clears:%d:%s", mode, country), redis.Z{Score: float64(clears), Member: stats.UserId})
	}
	_, err := pipe.Exec(service.ctx)
	return err
}

func (service *RankingsService) Remove(userId int, country string) error {
	if service == nil || service.client == nil {
		return ErrRedisClientNotInitialized
	}
	country = strings.ToLower(country)
	pipe := service.client.Pipeline()

	for mode := range constants.Modes {
		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:performance:%d", mode), userId)
		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:performance:%d:%s", mode, country), userId)

		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:rscore:%d", mode), userId)
		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:rscore:%d:%s", mode, country), userId)

		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:tscore:%d", mode), userId)
		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:tscore:%d:%s", mode, country), userId)

		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:ppv1:%d", mode), userId)
		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:ppv1:%d:%s", mode, country), userId)

		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:acc:%d", mode), userId)
		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:acc:%d:%s", mode, country), userId)

		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:clears:%d", mode), userId)
		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:clears:%d:%s", mode, country), userId)

		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:leader:%d", mode), userId)
		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:leader:%d:%s", mode, country), userId)
	}

	pipe.ZRem(service.ctx, "bancho:kudosu", userId)
	pipe.ZRem(service.ctx, fmt.Sprintf("bancho:kudosu:%s", country), userId)

	_, err := pipe.Exec(service.ctx)
	return err
}

func (service *RankingsService) RemoveFromCountry(userId int, country string) error {
	if service == nil || service.client == nil {
		return ErrRedisClientNotInitialized
	}
	country = strings.ToLower(country)
	pipe := service.client.Pipeline()

	for mode := range constants.Modes {
		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:performance:%d:%s", mode, country), userId)
		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:rscore:%d:%s", mode, country), userId)
		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:tscore:%d:%s", mode, country), userId)
		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:ppv1:%d:%s", mode, country), userId)
		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:acc:%d:%s", mode, country), userId)
		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:clears:%d:%s", mode, country), userId)
		pipe.ZRem(service.ctx, fmt.Sprintf("bancho:leader:%d:%s", mode, country), userId)
	}

	pipe.ZRem(service.ctx, fmt.Sprintf("bancho:kudosu:%s", country), userId)
	_, err := pipe.Exec(service.ctx)
	return err
}

// ScoresProvider prevents circular dependency between RankingsService and ScoresRepository
type ScoresProvider interface {
	FetchLeaderCount(userId int, mode constants.Mode) (int, error)
}

func (service *RankingsService) UpdateLeaderScores(stats *schemas.Stats, country string, scoreRepo ScoresProvider) error {
	if service == nil || service.client == nil {
		return ErrRedisClientNotInitialized
	}
	if stats == nil {
		return ErrStatsIsNil
	}

	leaderCount, err := scoreRepo.FetchLeaderCount(stats.UserId, stats.Mode)
	if err != nil {
		return err
	}

	entry := redis.Z{Score: float64(leaderCount), Member: stats.UserId}
	country = strings.ToLower(country)

	pipe := service.client.Pipeline()
	pipe.ZAdd(service.ctx, fmt.Sprintf("bancho:leader:%d", stats.Mode), entry)

	if len(country) == 2 && country != "xx" {
		pipe.ZAdd(service.ctx, fmt.Sprintf("bancho:leader:%d:%s", stats.Mode, country), entry)
	}
	_, err = pipe.Exec(service.ctx)
	return err
}

// KudosuProvider prevents circular dependency between RankingsService and BeatmapModdingRepository
type KudosuProvider interface {
	TotalKudosuByUser(userId int) (int, error)
}

// TODO: Deprecate cached kudosu rankings and use db queries instead

func (service *RankingsService) UpdateKudosu(userId int, country string, moddingRepo KudosuProvider) error {
	if service == nil || service.client == nil {
		return ErrRedisClientNotInitialized
	}

	kudosu, err := moddingRepo.TotalKudosuByUser(userId)
	if err != nil {
		return err
	}

	entry := redis.Z{Score: float64(kudosu), Member: userId}
	country = strings.ToLower(country)

	pipe := service.client.Pipeline()
	pipe.ZAdd(service.ctx, "bancho:kudosu", entry)

	if len(country) == 2 && country != "xx" {
		pipe.ZAdd(service.ctx, fmt.Sprintf("bancho:kudosu:%s", country), entry)
	}
	_, err = pipe.Exec(service.ctx)
	return err
}
