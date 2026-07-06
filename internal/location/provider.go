package location

import "sync"

// Provider defines how to interface with a geolocation backend.
type Provider interface {
	// Setup prepares the provider for use.
	Setup() error

	// Resolve returns a Location for the given IP address.
	Resolve(ip string) (*Location, error)
}

// provider combines all available providers, using local sources first
// and falling back to web sources if needed. Successful lookups are
// cached in memory and keyed by IP address.
type provider struct {
	web *WebProvider

	mutex sync.RWMutex
	cache map[string]*Location
}

func NewProvider() Provider {
	return &provider{
		cache: make(map[string]*Location),
	}
}

func (p *provider) Setup() error {
	p.web = NewWebProvider()
	return p.web.Setup()
}

func (p *provider) Resolve(ip string) (*Location, error) {
	if location, ok := p.lookup(ip); ok {
		return location, nil
	}

	// TODO: Add other proivders such as mmdb, etc. and combine them here
	location, err := p.web.Resolve(ip)
	if err != nil {
		// Don't cache failures, so they can be retried later on
		return location, err
	}

	p.store(ip, location)
	return location, nil
}

func (p *provider) lookup(ip string) (*Location, bool) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	location, ok := p.cache[ip]
	return location, ok
}

func (p *provider) store(ip string, location *Location) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.cache[ip] = location
}
