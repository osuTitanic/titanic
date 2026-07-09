package location

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

// roundTripFunc implements the http.RoundTripper interface for testing
type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	// Call the function to handle the request & return the response
	return fn(req)
}

func TestProvider(t *testing.T) {
	resolver := NewProvider()
	if err := resolver.Setup(); err != nil {
		t.Fatalf("failed to setup provider: %v", err)
	}
	requests := 0

	// Create a simulated HTTP client that returns a predefined response for the location lookup.
	// We don't want to make calls to the actual ip-api.com for unit tests.
	simulatedRequest := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		requests++
		if req.URL.String() != "http://ip-api.com/json/1.1.1.1" {
			t.Fatalf("unexpected location lookup url: %s", req.URL.String())
		}

		body := `{
			"status": "success",
			"country": "Australia",
			"countryCode": "AU",
			"city": "Brisbane",
			"lat": -27.4698,
			"lon": 153.0251,
			"timezone": "Australia/Brisbane",
			"query": "1.1.1.1"
		}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})

	// Override the HTTP client in the provider with our simulated client
	resolver.(*provider).web.client = &http.Client{
		Transport: simulatedRequest,
	}

	result, err := resolver.Resolve("1.1.1.1")
	if err != nil {
		t.Fatalf("failed to resolve IP: %v", err)
	}
	if result == nil {
		t.Fatal("expected a location result, got nil")
	}
	if result.IP != "1.1.1.1" {
		t.Fatalf("IP = %q, want %q", result.IP, "1.1.1.1")
	}
	if result.CountryCode != "AU" {
		t.Fatalf("CountryCode = %q, want %q", result.CountryCode, "AU")
	}
	if result.City != "Brisbane" {
		t.Fatalf("City = %q, want %q", result.City, "Brisbane")
	}

	cached, err := resolver.Resolve("1.1.1.1")
	if err != nil {
		t.Fatalf("failed to resolve cached IP: %v", err)
	}
	if cached != result {
		t.Fatal("expected cached location result")
	}
	if requests != 1 {
		t.Fatalf("location lookup requests = %d, want 1", requests)
	}
}
