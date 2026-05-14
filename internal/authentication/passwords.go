package authentication

import (
	"crypto/md5"
	"encoding/hex"
	"maps"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

var (
	passwordCache   = map[string]bool{}
	passwordCacheMu sync.RWMutex
)

func PasswordMd5(password string) string {
	sum := md5.Sum([]byte(password))
	return hex.EncodeToString(sum[:])
}

func CreatePasswordHash(password string) (string, error) {
	return CreatePasswordHashFromMd5(PasswordMd5(password))
}

func VerifyPasswordHash(password string, hash string) bool {
	return VerifyPasswordHashFromMd5(PasswordMd5(password), hash)
}

func CreatePasswordHashFromMd5(md5 string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(md5),
		bcrypt.DefaultCost,
	)
	return string(hashedBytes), err
}

func GetPasswordCache() map[string]bool {
	passwordCacheMu.RLock()
	defer passwordCacheMu.RUnlock()

	cloned := make(map[string]bool, len(passwordCache))
	maps.Copy(cloned, passwordCache)
	return cloned
}

func ClearPasswordCache() {
	passwordCacheMu.Lock()
	passwordCache = map[string]bool{}
	passwordCacheMu.Unlock()
}

func VerifyPasswordHashFromMd5(md5Hex string, hash string) bool {
	cacheKey := md5Hex + ":" + hash

	passwordCacheMu.RLock()
	cached, exists := passwordCache[cacheKey]
	passwordCacheMu.RUnlock()
	if exists {
		return cached
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(md5Hex))
	ok := err == nil

	passwordCacheMu.Lock()
	passwordCache[cacheKey] = ok
	passwordCacheMu.Unlock()
	return ok
}
