package rankings

import (
	"strconv"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/redis/go-redis/v9"
)

func (service *RankingsService) PlayerAbove(userId int, mode constants.Mode, rankType string) (scoreDifference int64, aboveUserId int, err error) {
	if service == nil || service.client == nil {
		return 0, 0, ErrRedisClientNotInitialized
	}

	key := service.RankingKey(mode, rankType, nil)
	position, err := service.client.ZRevRank(service.ctx, key, strconv.Itoa(userId)).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, 0, ErrNoPlayerAbove
		}
		return 0, 0, err
	}
	if position <= 0 {
		return 0, 0, ErrNoPlayerAbove
	}

	userScore, err := service.client.ZScore(service.ctx, key, strconv.Itoa(userId)).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, 0, ErrNoPlayerAbove
		}
		return 0, 0, err
	}

	aboveEntries, err := service.client.ZRevRangeWithScores(service.ctx, key, position-1, position-1).Result()
	if err != nil {
		return 0, 0, err
	}
	if len(aboveEntries) == 0 {
		return 0, 0, ErrNoPlayerAbove
	}

	aboveId, ok := redisMemberToInt(aboveEntries[0].Member)
	if !ok {
		return 0, 0, ErrNoPlayerAbove
	}

	return int64(aboveEntries[0].Score) - int64(userScore), aboveId, nil
}
