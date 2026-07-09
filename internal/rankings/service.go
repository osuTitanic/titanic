package rankings

import (
	"context"
	"fmt"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/redis/go-redis/v9"
)

type RankingsService struct {
	ctx    context.Context
	client *redis.Client
}

func (service *RankingsService) RankingKey(mode constants.Mode, rankType string, country *string) string {
	countrySuffix := ""
	if country != nil {
		countrySuffix = ":" + strings.ToLower(*country)
	}
	return fmt.Sprintf("bancho:%s:%d%s", rankType, mode, countrySuffix)
}

func (service *RankingsService) RankingKeyNoMode(rankType string, country *string) string {
	countrySuffix := ""
	if country != nil {
		countrySuffix = ":" + strings.ToLower(*country)
	}
	return fmt.Sprintf("bancho:%s%s", rankType, countrySuffix)
}

func NewRankingsService(client *redis.Client) *RankingsService {
	return &RankingsService{
		ctx:    context.Background(),
		client: client,
	}
}
