package location

import "testing"

func TestProvider(t *testing.T) {
	provider := NewProvider()
	if err := provider.Setup(); err != nil {
		t.Fatalf("failed to setup provider: %v", err)
	}

	result, err := provider.Resolve("1.1.1.1")
	if err != nil {
		t.Fatalf("failed to resolve IP: %v", err)
	}
	if result == nil {
		t.Fatal("expected a location result, got nil")
	}

	t.Logf("Resolved location: %+v", result)
}
