package location

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/osuTitanic/titanic/internal/caching"
)

// Provider defines how to interface with a geolocation backend.
type Provider interface {
	// Setup prepares the provider for use.
	Setup() error

	// Resolve returns a Location for the given IP address.
	Resolve(ip string) (*Location, error)

	// Close releases resources held by the provider.
	Close() error
}

// provider combines all available providers, using local sources first
// and falling back to web sources if needed. Successful lookups are
// cached in memory and keyed by IP address.
type provider struct {
	geoLite Provider
	web     Provider
	logger  *slog.Logger

	mutex sync.RWMutex
	cache *caching.Cache[string, *Location]
}

func NewProvider(databasePath, downloadUrl string) Provider {
	return newProviderFromInterfaces(
		NewGeoLiteProvider(databasePath, downloadUrl),
		NewWebProvider(),
	)
}

func newProviderFromInterfaces(geoLite, web Provider) *provider {
	return &provider{
		geoLite: geoLite,
		web:     web,
		logger:  slog.Default().With("component", "location"),
		cache:   caching.New[string, *Location](time.Hour),
	}
}

func (p *provider) Setup() error {
	if p.web == nil {
		return errors.New("location: web provider is not configured")
	}
	if err := p.web.Setup(); err != nil {
		return fmt.Errorf("location: failed to setup web provider: %w", err)
	}

	if p.geoLite == nil {
		p.logger.Warn("GeoLite provider is not configured, falling back to web provider")
		return nil
	}
	if err := p.geoLite.Setup(); err != nil {
		p.logger.Warn(
			"Failed to setup GeoLite provider, falling back to web provider",
			"error", err,
		)
		p.geoLite.Close()
		p.geoLite = nil
	}

	return nil
}

func (p *provider) Resolve(ip string) (*Location, error) {
	// Will either load from cache, or populate it using the resolver function
	return p.cache.GetOrLoad(ip, func() (*Location, error) {
		var geoliteErr error
		if p.geoLite != nil {
			location, err := p.geoLite.Resolve(ip)
			if err == nil {
				// Successfully resolved through geolite -> cached result
				return location, nil
			}
			geoliteErr = fmt.Errorf("location: GeoLite lookup failed: %w", err)
		}

		location, err := p.web.Resolve(ip)
		if err != nil {
			// Don't cache failures, so they can be retried later on
			return location, errors.Join(
				geoliteErr,
				fmt.Errorf("location: web lookup failed: %w", err),
			)
		}

		// Resolved through web -> cache result
		return location, nil
	})
}

func (p *provider) Close() error {
	var geoLiteErr error
	var webErr error

	if p.geoLite != nil {
		if err := p.geoLite.Close(); err != nil {
			geoLiteErr = fmt.Errorf("location: failed to close GeoLite provider: %w", err)
		}
	}
	if p.web != nil {
		if err := p.web.Close(); err != nil {
			webErr = fmt.Errorf("location: failed to close web provider: %w", err)
		}
	}
	p.cache.Clear()

	return errors.Join(geoLiteErr, webErr)
}
