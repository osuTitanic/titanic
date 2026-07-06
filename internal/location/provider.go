package location

// Provider defines how to interface with a geolocation backend.
type Provider interface {
	// Setup prepares the provider for use.
	Setup() error

	// Resolve returns a Location for the given IP address.
	Resolve(ip string) (*Location, error)
}

// provider combines all available providers, using local sources first
// and falling back to web sources if needed.
type provider struct {
	web *WebProvider
}

func NewProvider() Provider {
	return &provider{}
}

func (p *provider) Setup() error {
	p.web = NewWebProvider()
	return p.web.Setup()
}

func (p *provider) Resolve(ip string) (*Location, error) {
	// TODO: Add other proivders such as mmdb, etc. and combine them here
	return p.web.Resolve(ip)
}
