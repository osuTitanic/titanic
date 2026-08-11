package location

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/osuTitanic/titanic/internal/testkit"
)

func TestProviderGeoLite(t *testing.T) {
	databasePath := filepath.Join("..", "..", ".data", "GeoLite2-City.mmdb")
	if _, err := os.Stat(databasePath); err != nil {
		t.Skipf("GeoLite database is unavailable at %q: %v", databasePath, err)
	}

	fakeWeb := &testingProvider{resolveErr: errors.New("web provider should not be used")}
	geoLite := NewGeoLiteProvider(databasePath, "")
	resolver := newProviderFromInterfaces(geoLite, fakeWeb)

	t.Cleanup(func() {
		if err := resolver.Close(); err != nil {
			t.Errorf("failed to close provider: %v", err)
		}
	})
	if err := resolver.Setup(); err != nil {
		t.Fatalf("failed to setup provider: %v", err)
	}
	if resolver.geoLite == nil {
		t.Fatal("expected GeoLite database to be loaded")
	}

	result, err := resolver.Resolve("1.1.1.1")
	if err != nil {
		t.Fatalf("failed to resolve IP with GeoLite: %v", err)
	}
	if result.IP != "1.1.1.1" {
		t.Fatalf("IP = %q, want %q", result.IP, "1.1.1.1")
	}
	if result.CountryCode == "" || result.CountryCode == "XX" {
		t.Fatalf("expected resolved country code, got %q", result.CountryCode)
	}
	if fakeWeb.resolveCalls != 0 {
		t.Fatalf("web lookups = %d, want 0", fakeWeb.resolveCalls)
	}
}

func TestProviderWebFallback(t *testing.T) {
	if !testkit.IsInternetAvailable() {
		t.Skip("internet connection is unavailable")
	}

	geoLite := &testingProvider{resolveErr: errors.New("address not found")}
	web := NewWebProvider()
	resolver := newProviderFromInterfaces(geoLite, web)

	if err := resolver.Setup(); err != nil {
		t.Fatalf("failed to setup provider: %v", err)
	}

	result, err := resolver.Resolve("1.1.1.1")
	if err != nil {
		t.Fatalf("failed to resolve IP with web fallback: %v", err)
	}
	if result.IP != "1.1.1.1" {
		t.Fatalf("IP = %q, want %q", result.IP, "1.1.1.1")
	}
	if result.CountryCode != "AU" {
		t.Fatalf("CountryCode = %q, want %q", result.CountryCode, "AU")
	}
	if result.CountryName != "Australia" {
		t.Fatalf("CountryName = %q, want %q", result.CountryName, "Australia")
	}

	cached, err := resolver.Resolve("1.1.1.1")
	if err != nil {
		t.Fatalf("failed to resolve cached IP: %v", err)
	}
	if cached != result {
		t.Fatal("expected cached web result")
	}
	if geoLite.resolveCalls != 1 {
		t.Fatalf("GeoLite lookups = %d, want 1", geoLite.resolveCalls)
	}
}

func TestProviderCachesResult(t *testing.T) {
	want := &Location{
		IP:          "1.1.1.1",
		CountryCode: "AU",
		CountryName: "Australia",
		City:        "Brisbane",
	}
	fakeGeoLite := &testingProvider{result: want}
	fakeWeb := &testingProvider{resolveErr: errors.New("web provider should not be used")}
	resolver := newProviderFromInterfaces(fakeGeoLite, fakeWeb)

	if err := resolver.Setup(); err != nil {
		t.Fatalf("failed to setup provider: %v", err)
	}

	result, err := resolver.Resolve("1.1.1.1")
	if err != nil {
		t.Fatalf("failed to resolve IP with GeoLite: %v", err)
	}
	if result != want {
		t.Fatal("expected result from GeoLite provider")
	}

	cached, err := resolver.Resolve("1.1.1.1")
	if err != nil {
		t.Fatalf("failed to resolve cached IP: %v", err)
	}
	if cached != result {
		t.Fatal("expected cached GeoLite result")
	}
	if fakeGeoLite.resolveCalls != 1 {
		t.Fatalf("GeoLite lookups = %d, want 1", fakeGeoLite.resolveCalls)
	}
	if fakeWeb.resolveCalls != 0 {
		t.Fatalf("web lookups = %d, want 0", fakeWeb.resolveCalls)
	}
}

type testingProvider struct {
	result       *Location
	setupErr     error
	resolveErr   error
	closeErr     error
	setupCalls   int
	resolveCalls int
	closeCalls   int
}

func (provider *testingProvider) Setup() error {
	provider.setupCalls++
	return provider.setupErr
}

func (provider *testingProvider) Resolve(string) (*Location, error) {
	provider.resolveCalls++
	return provider.result, provider.resolveErr
}

func (provider *testingProvider) Close() error {
	provider.closeCalls++
	return provider.closeErr
}
