package authentication

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic-go/internal/config"
	"github.com/redis/go-redis/v9"
)

const DefaultCSRFTTL = 24 * time.Hour

type CSRFStore struct {
	Redis redis.Cmdable
	TTL   time.Duration
}

func NewCSRFStore(client redis.Cmdable) *CSRFStore {
	return &CSRFStore{
		Redis: client,
		TTL:   DefaultCSRFTTL,
	}
}

func CSRFRedisKey(userId int) string {
	return "csrf:" + strconv.Itoa(userId)
}

func GenerateCSRFToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func RequiresCSRF(cfg *config.Config, origin string) bool {
	if cfg == nil || origin == "" {
		return false
	}

	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}

	return strings.EqualFold(parsed.Hostname(), "osu."+cfg.DomainName)
}

func (store *CSRFStore) Upsert(ctx context.Context, userId int) (string, error) {
	token, err := GenerateCSRFToken()
	if err != nil {
		return "", err
	}

	ttl := store.TTL
	if ttl <= 0 {
		ttl = DefaultCSRFTTL
	}

	if err := store.Redis.Set(ctx, CSRFRedisKey(userId), token, ttl).Err(); err != nil {
		return "", err
	}

	return token, nil
}

func (store *CSRFStore) Get(ctx context.Context, userId int) (string, error) {
	token, err := store.Redis.Get(ctx, CSRFRedisKey(userId)).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return token, nil
}

func (store *CSRFStore) Validate(ctx context.Context, userId int, token string) (bool, error) {
	if token == "" {
		return false, nil
	}

	stored, err := store.Get(ctx, userId)
	if err != nil {
		return false, err
	}

	return stored != "" && stored == token, nil
}
