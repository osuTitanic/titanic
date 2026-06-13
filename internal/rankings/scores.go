package rankings

import (
	"math"
	"strconv"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/redis/go-redis/v9"
)

type PlayerScore struct {
	UserId int
	Score  float64
}

func (service *RankingsService) ScoreByKey(key string, userId int) (float64, error) {
	if service == nil || service.client == nil {
		return 0, ErrRedisClientNotInitialized
	}

	value, err := service.client.ZScore(service.ctx, key, strconv.Itoa(userId)).Result()
	if err == nil {
		return value, nil
	}
	if err == redis.Nil {
		return 0, nil
	}
	return 0, err
}

func (service *RankingsService) Performance(userId int, mode constants.Mode) (float64, error) {
	value, err := service.ScoreByKey(service.RankingKey(mode, "performance", nil), userId)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (service *RankingsService) Score(userId int, mode constants.Mode) (int64, error) {
	value, err := service.ScoreByKey(service.RankingKey(mode, "rscore", nil), userId)
	if err != nil {
		return 0, err
	}
	return int64(math.Round(value)), nil
}

func (service *RankingsService) TotalScore(userId int, mode constants.Mode) (int64, error) {
	value, err := service.ScoreByKey(service.RankingKey(mode, "tscore", nil), userId)
	if err != nil {
		return 0, err
	}
	return int64(math.Round(value)), nil
}

func (service *RankingsService) Accuracy(userId int, mode constants.Mode) (float64, error) {
	value, err := service.ScoreByKey(service.RankingKey(mode, "acc", nil), userId)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (service *RankingsService) PPv1(userId int, mode constants.Mode) (float64, error) {
	value, err := service.ScoreByKey(service.RankingKey(mode, "ppv1", nil), userId)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (service *RankingsService) Clears(userId int, mode constants.Mode) (int64, error) {
	value, err := service.ScoreByKey(service.RankingKey(mode, "clears", nil), userId)
	if err != nil {
		return 0, err
	}
	return int64(math.Round(value)), nil
}

func (service *RankingsService) LeaderScores(userId int, mode constants.Mode) (int64, error) {
	value, err := service.ScoreByKey(service.RankingKey(mode, "leader", nil), userId)
	if err != nil {
		return 0, err
	}
	return int64(math.Round(value)), nil
}

func (service *RankingsService) Kudosu(userId int, country *string) (int64, error) {
	// TODO: Deprecate cached kudosu rankings and use db queries instead
	value, err := service.ScoreByKey(service.RankingKeyNoMode("kudosu", country), userId)
	if err != nil {
		return 0, err
	}
	return int64(math.Round(value)), nil
}

func (service *RankingsService) PlayerCount(mode constants.Mode, rankType string, country *string) (int64, error) {
	return service.playerCountByKey(service.RankingKey(mode, rankType, country))
}

func (service *RankingsService) PlayerCountKudosu() (int64, error) {
	return service.playerCountByKey(service.RankingKeyNoMode("kudosu", nil))
}

func (service *RankingsService) playerCountByKey(key string) (int64, error) {
	if service == nil || service.client == nil {
		return 0, ErrRedisClientNotInitialized
	}
	return service.client.ZCount(service.ctx, key, "1", "+inf").Result()
}

func (service *RankingsService) TopPlayers(mode constants.Mode, offset int64, count int64, rankType string, country *string) ([]*PlayerScore, error) {
	return service.topPlayersByKey(service.RankingKey(mode, rankType, country), offset, count)
}

func (service *RankingsService) TopKudosu(offset int64, count int64) ([]*PlayerScore, error) {
	// TODO: Deprecate cached kudosu rankings and use db queries instead
	return service.topPlayersByKey(service.RankingKeyNoMode("kudosu", nil), offset, count)
}

func (service *RankingsService) topPlayersByKey(key string, offset int64, count int64) ([]*PlayerScore, error) {
	if service == nil || service.client == nil {
		return nil, ErrRedisClientNotInitialized
	}

	query := &redis.ZRangeBy{
		Max:    "+inf",
		Min:    "1",
		Offset: offset,
		Count:  count,
	}

	response := service.client.ZRevRangeByScoreWithScores(service.ctx, key, query)
	players, err := response.Result()
	if err != nil {
		return nil, err
	}

	result := make([]*PlayerScore, 0, len(players))
	for _, player := range players {
		// Unfortunately, redis returns an interface here, so
		// we need to painfully convert this back
		parsedId, ok := redisMemberToInt(player.Member)
		if !ok {
			continue
		}

		player := &PlayerScore{UserId: parsedId, Score: player.Score}
		result = append(result, player)
	}

	return result, nil
}

func redisMemberToInt(member any) (int, bool) {
	switch value := member.(type) {
	case int:
		return value, true
	case int64:
		return int(value), true
	case uint64:
		return int(value), true
	case float64:
		return int(value), true
	case string:
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return 0, false
		}
		return parsed, true
	case []byte:
		parsed, err := strconv.Atoi(string(value))
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}
