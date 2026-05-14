package authentication

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/osuTitanic/titanic-go/internal/schemas"
)

var (
	ErrEmptyTokenSecret = errors.New("authentication: token secret is empty")
	ErrInvalidToken     = errors.New("authentication: invalid token")
	ErrExpiredToken     = errors.New("authentication: expired token")
	ErrInvalidTokenType = errors.New("authentication: invalid token type")
	ErrMissingTokenId   = errors.New("authentication: missing token id")
)

type TokenSource int

const (
	TokenSourceWeb TokenSource = iota
	TokenSourceAPI
	TokenSourceOAuth2
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type TokenClaims struct {
	Id       int         `json:"id"`
	Name     string      `json:"name"`
	Expiry   int64       `json:"exp"`
	IssuedAt int64       `json:"iat"`
	Source   TokenSource `json:"source"`
	Type     TokenType   `json:"type"`
	TokenId  string      `json:"jti"`
}

type TokenPair struct {
	AccessToken   string
	RefreshToken  string
	AccessExpiry  time.Time
	RefreshExpiry time.Time
}

func GenerateToken(secret string, user *schemas.User, expiry time.Time, source TokenSource, tokenType TokenType) (string, error) {
	if user == nil {
		return "", fmt.Errorf("authentication: user is nil")
	}

	tokenId, err := GenerateTokenId()
	if err != nil {
		return "", err
	}

	return GenerateTokenClaims(secret, TokenClaims{
		Id:       user.Id,
		Name:     user.Name,
		Expiry:   expiry.Unix(),
		IssuedAt: time.Now().Unix(),
		Source:   source,
		Type:     tokenType,
		TokenId:  tokenId,
	})
}

func GenerateTokenClaims(secret string, claims TokenClaims) (string, error) {
	if secret == "" {
		return "", ErrEmptyTokenSecret
	}
	if claims.Type != TokenTypeAccess && claims.Type != TokenTypeRefresh {
		return "", ErrInvalidTokenType
	}
	if claims.TokenId == "" {
		return "", ErrMissingTokenId
	}
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("authentication: marshal token header: %w", err)
	}

	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("authentication: marshal token claims: %w", err)
	}

	headerSegment := base64.RawURLEncoding.EncodeToString(headerBytes)
	claimsSegment := base64.RawURLEncoding.EncodeToString(claimsBytes)
	signingInput := headerSegment + "." + claimsSegment
	signature := signToken(secret, signingInput)

	return signingInput + "." + signature, nil
}

func GenerateTokenPair(
	secret string,
	user *schemas.User,
	now time.Time,
	accessTTL, refreshTTL time.Duration,
	source TokenSource,
) (*TokenPair, error) {
	if now.IsZero() {
		now = time.Now()
	}

	pair := &TokenPair{
		AccessExpiry:  now.Add(accessTTL),
		RefreshExpiry: now.Add(refreshTTL),
	}

	var err error
	pair.AccessToken, err = GenerateToken(secret, user, pair.AccessExpiry, source, TokenTypeAccess)
	if err != nil {
		return nil, err
	}

	pair.RefreshToken, err = GenerateToken(secret, user, pair.RefreshExpiry, source, TokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	return pair, nil
}

func ValidateToken(token string, secret string) (*TokenClaims, error) {
	return ValidateTokenAt(token, secret, time.Now())
}

func ValidateTokenAt(token string, secret string, now time.Time) (*TokenClaims, error) {
	if secret == "" {
		return nil, ErrEmptyTokenSecret
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	signingInput := parts[0] + "." + parts[1]
	expectedSignature := signToken(secret, signingInput)
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return nil, ErrInvalidToken
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("%w: decode claims", ErrInvalidToken)
	}

	var claims TokenClaims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, fmt.Errorf("%w: decode payload", ErrInvalidToken)
	}

	if now.IsZero() {
		now = time.Now()
	}
	if now.Unix() > claims.Expiry {
		return nil, ErrExpiredToken
	}
	if claims.Type != TokenTypeAccess && claims.Type != TokenTypeRefresh {
		return nil, ErrInvalidTokenType
	}
	if claims.TokenId == "" {
		return nil, ErrMissingTokenId
	}

	return &claims, nil
}

func ValidateTokenType(token string, secret string, tokenType TokenType) (*TokenClaims, error) {
	return ValidateTokenTypeAt(token, secret, tokenType, time.Now())
}

func ValidateTokenTypeAt(token string, secret string, tokenType TokenType, now time.Time) (*TokenClaims, error) {
	claims, err := ValidateTokenAt(token, secret, now)
	if err != nil {
		return nil, err
	}
	if claims.Type != tokenType {
		return nil, ErrInvalidTokenType
	}
	return claims, nil
}

func GenerateTokenId() (string, error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("authentication: generate token id: %w", err)
	}
	return hex.EncodeToString(randomBytes), nil
}

func signToken(secret string, signingInput string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(signingInput))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
