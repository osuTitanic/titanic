package authentication

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const DefaultWebsiteSessionPrefix = "authentication:website:"

type WebsiteSession struct {
	Id        string `json:"id"`
	UserId    int    `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
	ExpiresAt int64  `json:"expires_at"`
}

type WebsiteSessionStore struct {
	Redis  redis.Cmdable
	Prefix string
}

func NewWebsiteSessionStore(client redis.Cmdable) *WebsiteSessionStore {
	return &WebsiteSessionStore{
		Redis:  client,
		Prefix: DefaultWebsiteSessionPrefix,
	}
}

func WebsiteSessionRedisKey(sessionId string) string {
	return DefaultWebsiteSessionPrefix + sessionId
}

func (store *WebsiteSessionStore) WebsiteSessionRedisKey(sessionId string) string {
	prefix := DefaultWebsiteSessionPrefix
	if store != nil && store.Prefix != "" {
		prefix = store.Prefix
	}
	return prefix + sessionId
}

func NewWebsiteSession(userId int, now time.Time, ttl time.Duration) (*WebsiteSession, error) {
	if now.IsZero() {
		now = time.Now()
	}
	if ttl <= 0 {
		return nil, ErrExpiredToken
	}

	sessionId, err := GenerateTokenId()
	if err != nil {
		return nil, err
	}

	return &WebsiteSession{
		Id:        sessionId,
		UserId:    userId,
		CreatedAt: now.Unix(),
		ExpiresAt: now.Add(ttl).Unix(),
	}, nil
}

func (store *WebsiteSessionStore) Create(ctx context.Context, userId int, now time.Time, ttl time.Duration) (*WebsiteSession, error) {
	session, err := NewWebsiteSession(userId, now, ttl)
	if err != nil {
		return nil, err
	}

	if err := store.Upsert(ctx, session, now); err != nil {
		return nil, err
	}

	return session, nil
}

func (store *WebsiteSessionStore) Upsert(ctx context.Context, session *WebsiteSession, now time.Time) error {
	if store == nil || store.Redis == nil {
		return ErrMissingSessionStore
	}
	if session == nil {
		return ErrMissingClaims
	}
	if session.Id == "" {
		return ErrMissingTokenId
	}
	if now.IsZero() {
		now = time.Now()
	}

	ttl := time.Unix(session.ExpiresAt, 0).Sub(now)
	if ttl <= 0 {
		return ErrExpiredToken
	}

	payload, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("authentication: marshal website session: %w", err)
	}

	return store.Redis.Set(ctx, store.WebsiteSessionRedisKey(session.Id), payload, ttl).Err()
}

func (store *WebsiteSessionStore) Get(ctx context.Context, sessionId string) (*WebsiteSession, error) {
	if store == nil || store.Redis == nil {
		return nil, ErrMissingSessionStore
	}
	if sessionId == "" {
		return nil, ErrMissingTokenId
	}

	payload, err := store.Redis.Get(ctx, store.WebsiteSessionRedisKey(sessionId)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var session WebsiteSession
	if err := json.Unmarshal(payload, &session); err != nil {
		return nil, fmt.Errorf("authentication: decode website session: %w", err)
	}

	return &session, nil
}

func (store *WebsiteSessionStore) Validate(ctx context.Context, sessionId string, now time.Time) (*WebsiteSession, error) {
	if now.IsZero() {
		now = time.Now()
	}

	session, err := store.Get(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}
	if now.Unix() > session.ExpiresAt {
		return nil, nil
	}

	return session, nil
}

func (store *WebsiteSessionStore) Delete(ctx context.Context, sessionId string) error {
	if store == nil || store.Redis == nil {
		return ErrMissingSessionStore
	}
	if sessionId == "" {
		return ErrMissingTokenId
	}

	return store.Redis.Del(ctx, store.WebsiteSessionRedisKey(sessionId)).Err()
}
