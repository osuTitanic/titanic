package state

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/osuTitanic/titanic/internal/authentication"
	"github.com/osuTitanic/titanic/internal/config"
	"github.com/osuTitanic/titanic/internal/database"
	"github.com/osuTitanic/titanic/internal/discord"
	"github.com/osuTitanic/titanic/internal/email"
	"github.com/osuTitanic/titanic/internal/location"
	"github.com/osuTitanic/titanic/internal/logging"
	"github.com/osuTitanic/titanic/internal/performance"
	"github.com/osuTitanic/titanic/internal/permissions"
	"github.com/osuTitanic/titanic/internal/rankings"
	"github.com/osuTitanic/titanic/internal/resources"
	"github.com/osuTitanic/titanic/internal/storage"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// State holds the application state, including database
// repositories, configuration, logger, storage & email services.
type State struct {
	*Repositories

	// Core components
	Config     *config.Config
	Logger     *slog.Logger
	Database   *gorm.DB
	Redis      *redis.Client
	Storage    storage.Storage
	Email      email.Email
	Officer    *discord.Officer
	Location   location.Provider
	Extensions map[string]any

	// Services
	Permissions permissions.Resolver
	Resources   resources.BeatmapResourceProvider
	Rankings    *rankings.RankingsService
	PPv1        *performance.PPv1Service
	PPv2        performance.IPPv2Service

	// Authentication
	SessionStore    *authentication.WebsiteSessionStore
	SessionStoreApi *authentication.SessionStore
	CSRFStore       *authentication.CSRFStore
}

func NewState(environmentFiles ...string) (*State, error) {
	cfg, err := config.LoadConfig(environmentFiles...)
	if err != nil {
		return nil, fmt.Errorf("state: failed to load config: %w", err)
	}

	logLevel := slog.LevelInfo
	if cfg.Debug {
		logLevel = slog.LevelDebug
	}
	logging.SetDefault("titanic", logLevel)
	logger := slog.Default()

	var storageProvider storage.Storage = storage.NewFileStorage(cfg.DataPath)
	var s3Config = cfg.S3Config()

	if cfg.S3Enabled && s3Config == nil {
		return nil, fmt.Errorf("state: S3 is enabled but S3 config is missing")
	}
	if cfg.S3Enabled {
		storageProvider, err = storage.NewS3Storage(*s3Config)
		if err != nil {
			return nil, fmt.Errorf("state: failed to create S3 storage: %w", err)
		}
	}

	if err := storageProvider.Setup(); err != nil {
		return nil, fmt.Errorf("state: failed to setup storage: %w", err)
	}

	db, err := database.CreateSession(cfg)
	if err != nil {
		return nil, fmt.Errorf("state: failed to create database session: %w", err)
	}

	mailer, err := email.NewEmailFromConfig(cfg)
	if err != nil {
		database.CloseSession(db)
		return nil, fmt.Errorf("state: failed to create email service: %w", err)
	}

	if err := mailer.Setup(); err != nil {
		database.CloseSession(db)
		return nil, fmt.Errorf("state: failed to setup email service: %w", err)
	}

	geolocation := location.NewProvider()
	if err := geolocation.Setup(); err != nil {
		database.CloseSession(db)
		return nil, fmt.Errorf("state: failed to setup location service: %w", err)
	}

	redisPassword := ""
	if cfg.RedisPass != nil {
		redisPassword = *cfg.RedisPass
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: redisPassword,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		database.CloseSession(db)
		return nil, fmt.Errorf("state: failed to connect to Redis: %w", err)
	}

	repos := NewRepositories(db)

	beatmapResources := resources.NewBeatmapProvider(
		cfg, redisClient, storageProvider,
		repos.ResourceMirrors,
		repos.Beatmapsets,
	)
	if err := beatmapResources.Setup(); err != nil {
		database.CloseSession(db)
		return nil, fmt.Errorf("state: failed to setup beatmap resources: %w", err)
	}

	return &State{
		Config:          cfg,
		Database:        db,
		Storage:         storageProvider,
		Logger:          logger,
		Email:           mailer,
		Officer:         discord.NewOfficerFromConfig(cfg),
		Location:        geolocation,
		Redis:           redisClient,
		Repositories:    repos,
		Resources:       beatmapResources,
		Extensions:      make(map[string]any),
		Rankings:        rankings.NewRankingsService(redisClient),
		PPv1:            performance.NewPPv1Service(repos.Scores, repos.Beatmaps),
		PPv2:            performance.NewPPv2Service(beatmapResources), // TODO: set caching layer
		Permissions:     permissions.New(repos.Permissions, repos.Groups),
		CSRFStore:       authentication.NewCSRFStore(redisClient),
		SessionStore:    authentication.NewWebsiteSessionStore(redisClient),
		SessionStoreApi: authentication.NewSessionStore(redisClient),
	}, nil
}

// DatabaseTransaction executes the given function within a database transaction.
// Example usage:
//
//	err := state.DatabaseTransaction(func(repos *Repositories) error {
//	    // Perform database operations using repos
//	 	repos.User.Create(...)
//
//		// If an error is returned, the transaction will be rolled back
//	    // If nil is returned, the transaction will be committed
//	    return nil
//	})
func (state *State) DatabaseTransaction(fn func(repositories *Repositories) error) error {
	if state == nil || state.Database == nil {
		return fmt.Errorf("state: database is not initialized")
	}
	return state.Database.Transaction(func(tx *gorm.DB) error {
		return fn(NewRepositories(tx))
	})
}

// Close gracefully closes the state
func (state *State) Close() error {
	if state == nil || state.Database == nil {
		return nil
	}
	return database.CloseSession(state.Database)
}
