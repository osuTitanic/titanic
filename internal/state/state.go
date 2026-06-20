package state

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/osuTitanic/titanic-go/internal/authentication"
	"github.com/osuTitanic/titanic-go/internal/config"
	"github.com/osuTitanic/titanic-go/internal/database"
	"github.com/osuTitanic/titanic-go/internal/email"
	"github.com/osuTitanic/titanic-go/internal/logging"
	"github.com/osuTitanic/titanic-go/internal/performance"
	"github.com/osuTitanic/titanic-go/internal/permissions"
	"github.com/osuTitanic/titanic-go/internal/rankings"
	"github.com/osuTitanic/titanic-go/internal/resources"
	"github.com/osuTitanic/titanic-go/internal/storage"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// State holds the application state, including database
// repositories, configuration, logger, storage & email services.
type State struct {
	*Repositories

	// Core components
	Config   *config.Config
	Logger   *slog.Logger
	Database *gorm.DB
	Redis    *redis.Client
	Storage  storage.Storage
	Email    email.Email

	// Services
	Permissions permissions.Resolver
	Resources   resources.BeatmapResourceProvider
	Rankings    *rankings.RankingsService
	PPv1        *performance.PPv1Service

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

	fs := storage.NewFileStorage(cfg.DataPath)
	if err := fs.Setup(); err != nil {
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
		cfg, redisClient, fs,
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
		Storage:         fs,
		Logger:          logger,
		Email:           mailer,
		Redis:           redisClient,
		Repositories:    repos,
		Resources:       beatmapResources,
		Rankings:        rankings.NewRankingsService(redisClient),
		PPv1:            performance.NewPPv1Service(repos.Scores, repos.Beatmaps),
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
