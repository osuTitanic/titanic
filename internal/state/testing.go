//go:build integration

package state

import (
	"io"
	"log/slog"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/osuTitanic/titanic-go/internal/authentication"
	"github.com/osuTitanic/titanic-go/internal/config"
	"github.com/osuTitanic/titanic-go/internal/database"
	"github.com/osuTitanic/titanic-go/internal/email"
	"github.com/osuTitanic/titanic-go/internal/location"
	"github.com/osuTitanic/titanic-go/internal/performance"
	"github.com/osuTitanic/titanic-go/internal/permissions"
	"github.com/osuTitanic/titanic-go/internal/rankings"
	"github.com/osuTitanic/titanic-go/internal/resources"
	"github.com/osuTitanic/titanic-go/internal/storage"
	"github.com/osuTitanic/titanic-go/internal/testkit"
	"github.com/redis/go-redis/v9"
)

// NOTE: Usage example for the test state: services/stern/cmd/web/main_test.go
// 		 Too lazy to document this right now in the readme.

type TestStateOptions struct {
	configure  []func(*config.Config)
	migrations []any
	logger     *slog.Logger
}
type TestStateOption func(*TestStateOptions)

// NewTestState creates an application state backed by temporary PostgreSQL & Redis containers
func NewTestState(t testing.TB, opts ...TestStateOption) *State {
	t.Helper()
	t.Log("creating new test state")

	options := &TestStateOptions{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	for _, opt := range opts {
		opt(options)
	}
	slog.SetDefault(options.logger)

	t.Log("creating new postgres session")
	cfg := testkit.PostgresConfig(t)
	cfg.DataPath = t.TempDir()
	applyTestConfigDefaults(cfg)

	t.Log("creating new redis session")
	redisClient := testkit.RedisClient(t)
	applyRedisConfig(t, cfg, redisClient)

	for _, configure := range options.configure {
		configure(cfg)
	}

	db, err := database.CreateSession(cfg)
	if err != nil {
		t.Fatalf("failed to create postgres session: %v", err)
	}
	t.Cleanup(func() {
		if err := database.CloseSession(db); err != nil {
			t.Fatalf("failed to close postgres session: %v", err)
		}
	})

	if len(options.migrations) > 0 {
		if err := db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto").Error; err != nil {
			t.Fatalf("failed to create pgcrypto extension: %v", err)
		}
		if err := db.AutoMigrate(options.migrations...); err != nil {
			t.Fatalf("failed to migrate test state schemas: %v", err)
		}
		// TODO: We would eventually want to run our own migrations, unless that takes too long
		// 		 Another option would be to add default test data to the database
	}

	geolocation := location.NewDummyProvider()
	if err := geolocation.Setup(); err != nil {
		t.Fatalf("failed to setup location service: %v", err)
	}

	// TODO: Add a more suitable provider for storage.
	// 		 Right now it depends on an internet connection to github.
	storage := storage.NewFileStorage(cfg.DataPath)
	if err := storage.Setup(); err != nil {
		t.Fatalf("failed to setup storage: %v", err)
	}
	repositories := NewRepositories(db)

	beatmapResources := resources.NewBeatmapProvider(
		cfg, redisClient, storage,
		repositories.ResourceMirrors,
		repositories.Beatmapsets,
	)
	if err := beatmapResources.Setup(); err != nil {
		t.Fatalf("failed to setup beatmap resources: %v", err)
	}

	return &State{
		Config:          cfg,
		Logger:          options.logger,
		Database:        db,
		Redis:           redisClient,
		Storage:         storage,
		Email:           email.NewNoopEmail(cfg.EmailSender),
		Location:        geolocation,
		Resources:       beatmapResources,
		Repositories:    repositories,
		Extensions:      map[string]any{},
		Rankings:        rankings.NewRankingsService(redisClient),
		PPv1:            performance.NewPPv1Service(repositories.Scores, repositories.Beatmaps),
		Permissions:     permissions.New(repositories.Permissions, repositories.Groups),
		CSRFStore:       authentication.NewCSRFStore(redisClient),
		SessionStore:    authentication.NewWebsiteSessionStore(redisClient),
		SessionStoreApi: authentication.NewSessionStore(redisClient),
	}
}

// WithTestConfig applies given configuration before the database session is opened
func WithTestConfig(configure func(*config.Config)) TestStateOption {
	return func(options *TestStateOptions) {
		if configure != nil {
			options.configure = append(options.configure, configure)
		}
	}
}

// WithTestMigrations runs AutoMigrate for the given schemas after postgres is available
func WithTestMigrations(models ...any) TestStateOption {
	return func(options *TestStateOptions) {
		options.migrations = append(options.migrations, models...)
	}
}

// WithTestLogger overrides the default logger
func WithTestLogger(logger *slog.Logger) TestStateOption {
	return func(options *TestStateOptions) {
		if logger != nil {
			options.logger = logger
		}
	}
}

func applyRedisConfig(t testing.TB, cfg *config.Config, client *redis.Client) {
	t.Helper()

	options := client.Options()
	host, portText, err := net.SplitHostPort(options.Addr)
	if err != nil {
		t.Fatalf("failed to parse redis address %q: %v", options.Addr, err)
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		t.Fatalf("failed to parse redis port %q: %v", portText, err)
	}

	cfg.RedisHost = host
	cfg.RedisPort = port
	if options.Password != "" {
		cfg.RedisPass = &options.Password
	}
}

func applyTestConfigDefaults(cfg *config.Config) {
	cfg.DomainName = "localhost"
	cfg.EmailProvider = "noop"
	cfg.EmailSender = "support@titanic.sh"
	cfg.FrontendSecretKey = "verysecret"
	cfg.ScoreResponseLimit = 50
	cfg.BeatmapFavoritesLimit = 100
	cfg.SuperFriendlyUsers = config.IntSlice{1}
	cfg.BeginningEndedAt = config.DynamicTime(time.Date(2023, 12, 31, 6, 0, 0, 0, time.UTC))
	cfg.WikiDefaultLanguage = "en"
	cfg.PostgresPoolEnabled = false
	cfg.ValidImageServicesOverride = config.StringSlice{}
}
