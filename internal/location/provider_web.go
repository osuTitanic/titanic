package location

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// WebProvider resolves geolocation data using ip-api.com.
type WebProvider struct {
	client *http.Client
	logger *slog.Logger
}

type locationResponse struct {
	Status      string  `json:"status"`
	Message     string  `json:"message"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	City        string  `json:"city"`
	Latitude    float64 `json:"lat"`
	Longitude   float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Query       string  `json:"query"`
}

func NewWebProvider() *WebProvider {
	return &WebProvider{
		client: &http.Client{Timeout: 5 * time.Second},
		logger: slog.Default().With("component", "location"),
	}
}

func (p *WebProvider) Setup() error {
	// bleh
	return nil
}

func (p *WebProvider) Resolve(ip string) (*Location, error) {
	p.logger.Debug("Resolving location", "ip", ip)

	// When requesting /json/ (without an IP set), the API provides our own location.
	// We can take advantage of that for local IPs, since they won't be resolvable otherwise.
	if IsLocalIP(ip) {
		ip = ""
	}

	location, err := p.Fetch(ip)
	if err != nil {
		return DefaultLocation(), err
	}

	return location, nil
}

func (p *WebProvider) Fetch(ip string) (*Location, error) {
	response, err := p.client.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		return nil, fmt.Errorf("location: request failed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("location: unexpected status code: %d", response.StatusCode)
	}

	var payload locationResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("location: failed to decode response: %w", err)
	}

	if payload.Status != "success" {
		return nil, fmt.Errorf("location: lookup failed: %s", payload.Message)
	}

	location := &Location{
		IP:        payload.Query,
		Latitude:  payload.Latitude,
		Longitude: payload.Longitude,
		Timezone:  payload.Timezone,
		City:      payload.City,
	}
	location.SetCountryCode(payload.CountryCode)
	return location, nil
}
