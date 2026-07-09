package authentication

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const DefaultSessionPrefix = "authentication:session:"

var (
	ErrMissingSessionStore = errors.New("authentication: missing session store")
	ErrMissingClaims       = errors.New("authentication: missing token claims")
)

type Session struct {
	TokenId  string      `json:"token_id"`
	UserId   int         `json:"user_id"`
	Source   TokenSource `json:"source"`
	Type     TokenType   `json:"type"`
	IssuedAt int64       `json:"issued_at"`
	Expiry   int64       `json:"expiry"`
}

type SessionStore struct {
	Redis  redis.Cmdable
	Prefix string
}

func NewSessionStore(client redis.Cmdable) *SessionStore {
	return &SessionStore{
		Redis:  client,
		Prefix: DefaultSessionPrefix,
	}
}

func SessionRedisKey(tokenId string) string {
	return DefaultSessionPrefix + tokenId
}

func (store *SessionStore) SessionRedisKey(tokenId string) string {
	prefix := DefaultSessionPrefix
	if store != nil && store.Prefix != "" {
		prefix = store.Prefix
	}
	return prefix + tokenId
}

func NewSessionFromClaims(claims *TokenClaims) (*Session, error) {
	if claims == nil {
		return nil, ErrMissingClaims
	}
	if claims.TokenId == "" {
		return nil, ErrMissingTokenId
	}
	if claims.Type != TokenTypeAccess && claims.Type != TokenTypeRefresh {
		return nil, ErrInvalidTokenType
	}

	return &Session{
		TokenId:  claims.TokenId,
		UserId:   claims.Id,
		Source:   claims.Source,
		Type:     claims.Type,
		IssuedAt: claims.IssuedAt,
		Expiry:   claims.Expiry,
	}, nil
}

func (store *SessionStore) Upsert(ctx context.Context, claims *TokenClaims, now time.Time) (*Session, error) {
	if store == nil || store.Redis == nil {
		return nil, ErrMissingSessionStore
	}
	if now.IsZero() {
		now = time.Now()
	}

	session, err := NewSessionFromClaims(claims)
	if err != nil {
		return nil, err
	}

	ttl := time.Unix(session.Expiry, 0).Sub(now)
	if ttl <= 0 {
		return nil, ErrExpiredToken
	}

	payload, err := json.Marshal(session)
	if err != nil {
		return nil, fmt.Errorf("authentication: marshal session: %w", err)
	}

	if err := store.Redis.Set(ctx, store.SessionRedisKey(session.TokenId), payload, ttl).Err(); err != nil {
		return nil, err
	}

	return session, nil
}

func (store *SessionStore) Get(ctx context.Context, tokenId string) (*Session, error) {
	if store == nil || store.Redis == nil {
		return nil, ErrMissingSessionStore
	}
	if tokenId == "" {
		return nil, ErrMissingTokenId
	}

	payload, err := store.Redis.Get(ctx, store.SessionRedisKey(tokenId)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(payload, &session); err != nil {
		return nil, fmt.Errorf("authentication: decode session: %w", err)
	}

	return &session, nil
}

func (store *SessionStore) Validate(ctx context.Context, claims *TokenClaims, now time.Time) (bool, error) {
	if claims == nil {
		return false, ErrMissingClaims
	}
	if now.IsZero() {
		now = time.Now()
	}
	if now.Unix() > claims.Expiry {
		return false, ErrExpiredToken
	}

	session, err := store.Get(ctx, claims.TokenId)
	if err != nil {
		return false, err
	}
	if session == nil {
		return false, nil
	}

	if session.TokenId != claims.TokenId {
		return false, nil
	}
	if session.UserId != claims.Id {
		return false, nil
	}
	if session.Source != claims.Source {
		return false, nil
	}
	if session.Type != claims.Type {
		return false, nil
	}
	if session.IssuedAt != claims.IssuedAt {
		return false, nil
	}
	if session.Expiry != claims.Expiry {
		return false, nil
	}
	if now.Unix() > session.Expiry {
		return false, nil
	}

	return true, nil
}

func (store *SessionStore) Delete(ctx context.Context, tokenId string) error {
	if store == nil || store.Redis == nil {
		return ErrMissingSessionStore
	}
	if tokenId == "" {
		return ErrMissingTokenId
	}

	return store.Redis.Del(ctx, store.SessionRedisKey(tokenId)).Err()
}

func (store *SessionStore) DeleteClaims(ctx context.Context, claims *TokenClaims) error {
	if claims == nil {
		return ErrMissingClaims
	}
	return store.Delete(ctx, claims.TokenId)
}
