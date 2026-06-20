package resources

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic-go/internal/config"
	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/repositories"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/redis/go-redis/v9"
)

// MirrorResolver receives beatmap resources from external mirrors over HTTP.
// It uses a round-robin system to rotate between mirrors & handles rate limits.
type MirrorResolver struct {
	config  *config.Config
	logger  *slog.Logger
	cache   *redis.Client
	mirrors *repositories.ResourceMirrorRepository
	server  constants.BeatmapServer
	session *httpSession
}

func NewMirrorResolver(
	server constants.BeatmapServer,
	cfg *config.Config,
	cache *redis.Client,
	mirrors *repositories.ResourceMirrorRepository,
) *MirrorResolver {
	return &MirrorResolver{
		logger:  slog.Default().With("component", "BeatmapMirror"),
		config:  cfg,
		cache:   cache,
		mirrors: mirrors,
		server:  server,
	}
}

func (resolver *MirrorResolver) Setup() error {
	userAgent := fmt.Sprintf("osuTitanic (%s)", resolver.config.DomainName)
	resolver.session = createHttpSession(userAgent)
	return nil
}

func (resolver *MirrorResolver) Osz(setId int, noVideo bool) (io.ReadCloser, error) {
	resolver.logger.Debug(
		"Downloading osz...",
		"set_id", setId,
	)

	resourceType := constants.BeatmapResourceTypeOsz
	if noVideo {
		resourceType = constants.BeatmapResourceTypeOszNoVideo
	}

	return resolver.FetchStream(resourceType, setId)
}

func (resolver *MirrorResolver) Osu(beatmapId int) (io.ReadCloser, error) {
	resolver.logger.Debug(
		"Downloading beatmap...",
		"beatmap_id", beatmapId,
	)

	mirrors, err := resolver.mirrors.FetchByTypeAll(constants.BeatmapResourceTypeBeatmap)
	if err != nil {
		return nil, err
	}

	// Special case for beatmaps: we want to check titanic first for beatmaps
	// TODO: Remove this special case & make the download_server consistent
	return resolver.FetchStreamFromMirrors(beatmapId, mirrors)
}

func (resolver *MirrorResolver) Preview(setId int) (io.ReadCloser, error) {
	resolver.logger.Debug("Downloading preview...", "set_id", setId)
	return resolver.FetchStream(constants.BeatmapResourceTypeAudio, setId)
}

func (resolver *MirrorResolver) Background(setId int, large bool) (io.ReadCloser, error) {
	resolver.logger.Debug("Downloading background...", "set_id", setId)

	resourceType := constants.BeatmapResourceTypeThumbnail
	if large {
		resourceType = constants.BeatmapResourceTypeBackground
	}

	return resolver.FetchStream(resourceType, setId)
}

// FetchStream resolves the mirrors for the given resource type and returns a
// stream to the first mirror that responds successfully.
func (resolver *MirrorResolver) FetchStream(resourceType constants.BeatmapResourceType, setId int) (io.ReadCloser, error) {
	mirrors := resolver.ResolveMirrors(resourceType, resolver.server)
	return resolver.FetchStreamFromMirrors(setId, mirrors)
}

// FetchStreamFromMirrors iterates through the provided mirrors, returning a stream
// from the first mirror that responds successfully.
func (resolver *MirrorResolver) FetchStreamFromMirrors(setId int, mirrors []*schemas.BeatmapMirror) (io.ReadCloser, error) {
	if len(mirrors) == 0 {
		return nil, ErrNoMirrorsAvailable
	}

	for _, mirror := range mirrors {
		response := resolver.PerformMirrorRequest(
			resolveMirrorUrl(mirror.Url, setId),
			mirror,
		)
		if response == nil {
			continue
		}

		if response.ContentLength == 0 {
			response.Body.Close()
			continue
		}
		return response.Body, nil
	}

	return nil, ErrResourceNotFound
}

// PerformMirrorRequest sends a request to a single mirror, returning its
// response or nil if the mirror is unavailable, rate limited or errored out.
func (resolver *MirrorResolver) PerformMirrorRequest(url string, mirror *schemas.BeatmapMirror) *http.Response {
	if resolver.CheckRatelimit(mirror.Url) {
		return nil
	}

	response, err := resolver.session.Get(url)
	if err != nil {
		resolver.logger.Error(
			"Failed to send request",
			"url", url, "error", err.Error(),
		)
		return nil
	}

	ratelimitRemaining, hasRemaining := resolveHeaderInt(response, "X-Ratelimit-Remaining")
	ratelimitReset, _ := resolveHeaderInt(response, "X-Ratelimit-Reset")

	if hasRemaining && ratelimitRemaining <= 1 {
		resolver.logger.Warn(
			"Remaining units low, blocking mirror",
			"mirror", mirror.Url, "seconds", ratelimitReset,
		)
		resolver.SetRatelimit(mirror.Url, ratelimitReset)
	}

	dailyRemaining, hasDaily := resolveHeaderInt(response, "X-Daily-Remaining")
	dailyReset, _ := resolveHeaderInt(response, "X-Ratelimit-Daily-Reset")

	if hasDaily && dailyRemaining <= 1 {
		resolver.logger.Warn(
			"Daily limit reached on mirror",
			"mirror", mirror.Url, "seconds", dailyReset,
		)
		resolver.SetRatelimit(mirror.Url, dailyReset)
	}

	if response.StatusCode == http.StatusTooManyRequests {
		retry, ok := resolveHeaderInt(response, "Retry-After")
		if !ok {
			retry = 120
		}
		resolver.logger.Warn(
			"Rate limited on mirror",
			"mirror", mirror.Url, "seconds", retry,
		)
		resolver.SetRatelimit(mirror.Url, retry)
		response.Body.Close()
		return nil
	}

	if response.StatusCode < 200 || response.StatusCode >= 400 {
		resolver.LogError(response.Request.URL.String(), response.StatusCode)
		response.Body.Close()
		return nil
	}

	return response
}

// CheckRatelimit reports whether the given mirror is currently rate limited.
func (resolver *MirrorResolver) CheckRatelimit(rawUrl string) bool {
	domain := resolveMirrorDomain(rawUrl)
	key := fmt.Sprintf("ratelimit:%s", domain)

	exists, err := resolver.cache.Exists(context.Background(), key).Result()
	return err == nil && exists > 0
}

// SetRatelimit blocks the given mirror for the provided number of seconds.
func (resolver *MirrorResolver) SetRatelimit(rawUrl string, seconds int) {
	if seconds <= 0 {
		seconds = 120
	}
	domain := resolveMirrorDomain(rawUrl)

	resolver.cache.Set(
		context.Background(),
		fmt.Sprintf("ratelimit:%s", domain), 1,
		time.Duration(seconds)*time.Second,
	)
	resolver.logger.Warn(
		"Rate limited mirror",
		"mirror", domain, "seconds", seconds,
	)
}

// ResolveMirrors returns the mirrors for the given type & server, rotated by a
// round-robin index and filtered to exclude rate limited mirrors.
func (resolver *MirrorResolver) ResolveMirrors(resourceType constants.BeatmapResourceType, server constants.BeatmapServer) []*schemas.BeatmapMirror {
	roundRobinKey := fmt.Sprintf("roundrobin:%d:%d", resourceType, server)
	index := resolver.RoundRobinIndex(roundRobinKey)

	available, err := resolver.mirrors.FetchByType(resourceType, server)
	if err != nil {
		resolver.logger.Error(
			"Failed to fetch mirrors",
			"error", err.Error(),
		)
		return nil
	}

	mirrors := make([]*schemas.BeatmapMirror, 0, len(available))
	for _, mirror := range available {
		if resolver.CheckRatelimit(mirror.Url) {
			continue
		}
		mirrors = append(mirrors, mirror)
	}
	if len(mirrors) == 0 {
		return nil
	}

	index = index % len(mirrors)
	nextIndex := (index + 1) % len(mirrors)

	resolver.cache.Set(
		context.Background(),
		roundRobinKey,
		nextIndex, 60*time.Second,
	)
	return append(mirrors[index:], mirrors[:index]...)
}

func (resolver *MirrorResolver) RoundRobinIndex(key string) int {
	value, err := resolver.cache.Get(context.Background(), key).Result()
	if err != nil {
		return 0
	}

	index, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return index
}

func (resolver *MirrorResolver) LogError(url string, statusCode int) {
	if url == "" || statusCode == 0 {
		return
	}

	if statusCode == http.StatusNotFound {
		resolver.logger.Debug(
			"Failed to find resource",
			"url", url, "status", statusCode,
		)
		return
	}

	resolver.logger.Error(
		"Error while sending request",
		"url", url, "status", statusCode,
	)
}

func resolveMirrorUrl(rawUrl string, id int) string {
	return strings.ReplaceAll(rawUrl, "{}", strconv.Itoa(id))
}

func resolveMirrorDomain(rawUrl string) string {
	parsed, err := url.Parse(resolveMirrorUrl(rawUrl, 0))
	if err != nil {
		return rawUrl
	}
	return parsed.Host
}

func resolveHeaderInt(response *http.Response, header string) (int, bool) {
	value := response.Header.Get(header)
	if value == "" {
		return 0, false
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}
	return parsed, true
}
