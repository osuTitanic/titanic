package templates

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/fs"
	"net/url"
	"path"
	"reflect"
	"strings"
	"sync"

	"github.com/CloudyKit/jet/v6"
)

// staticUrlCache is a cache for static file URLs that appends a content hash to the query string
// e.g. /js/main.js -> /js/main.js?c=<checksum> to allow for long-term caching of static files
type staticUrlCache struct {
	cache map[string]string
	files fs.FS
	mu    sync.RWMutex
}

func newStaticUrlCache(files fs.FS) *staticUrlCache {
	return &staticUrlCache{
		files: files,
		cache: make(map[string]string),
	}
}

// cachedUrl is a jet template function that returns the cached URL for a static file
// e.g. <script src="{{ cachedUrl("/js/main.js") }}"></script>
func (cache *staticUrlCache) cachedUrl(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("cachedUrl", 1, 1)
	urlPath := a.Get(0).String()
	return reflect.ValueOf(cache.resolve(urlPath))
}

func (cache *staticUrlCache) resolve(urlPath string) string {
	if urlPath == "" || cache.files == nil {
		return urlPath
	}

	cache.mu.RLock()
	cached, ok := cache.cache[urlPath]
	cache.mu.RUnlock()
	if ok {
		return cached
	}

	result := cache.resolveUncached(urlPath)

	cache.mu.Lock()
	cache.cache[urlPath] = result
	cache.mu.Unlock()

	return result
}

func (cache *staticUrlCache) resolveUncached(urlPath string) string {
	parsed, err := url.Parse(urlPath)
	if err != nil || parsed.IsAbs() || !strings.HasPrefix(parsed.Path, "/") {
		return urlPath
	}

	filename := strings.TrimPrefix(path.Clean(parsed.Path), "/")
	if filename == "" || filename == "." || strings.HasPrefix(filename, "../") {
		return urlPath
	}

	file, err := cache.files.Open(filename)
	if err != nil {
		return urlPath
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil || info.IsDir() {
		return urlPath
	}

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return urlPath
	}

	query := parsed.Query()
	query.Set("c", hex.EncodeToString(hash.Sum(nil)))
	parsed.RawQuery = query.Encode()
	return parsed.String()
}
