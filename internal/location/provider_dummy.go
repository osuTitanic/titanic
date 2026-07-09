package location

// DummyProvider is a geolocation provider for tests.
type DummyProvider struct {
	Result *Location
}

func NewDummyProvider() Provider {
	return &DummyProvider{Result: DefaultLocation()}
}

func (p *DummyProvider) Setup() error {
	return nil
}

func (p *DummyProvider) Resolve(ip string) (*Location, error) {
	copy := *p.Result
	location := &copy
	location.IP = ip
	return location, nil
}
