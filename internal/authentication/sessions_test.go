//go:build integration

package authentication

import (
	"context"
	"testing"
	"time"

	"github.com/osuTitanic/titanic-go/internal/testkit"
)

func TestWebsiteSessionStoreWithRedisContainer(t *testing.T) {
	// Spin up a temporary redis container to complete this test
	ctx := context.Background()
	client := testkit.RedisClient(t)
	store := NewWebsiteSessionStore(client)

	// Create a new website session and upsert it into the store
	now := time.Unix(1_700_000_000, 0).UTC()
	session := &WebsiteSession{
		Id:        "session",
		UserId:    69,
		CreatedAt: now.Unix(),
		ExpiresAt: now.Add(time.Hour).Unix(),
	}

	if err := store.Upsert(ctx, session, now); err != nil {
		t.Fatalf("failed to upsert website session: %v", err)
	}

	got, err := store.Validate(ctx, session.Id, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("failed to validate website session: %v", err)
	}
	if got == nil {
		t.Fatal("expected website session, got nil")
	}
	if got.Id != session.Id {
		t.Fatalf("session id = %q, want %q", got.Id, session.Id)
	}
	if got.UserId != session.UserId {
		t.Fatalf("user id = %d, want %d", got.UserId, session.UserId)
	}

	if err := store.Delete(ctx, session.Id); err != nil {
		t.Fatalf("failed to delete website session: %v", err)
	}
	got, err = store.Get(ctx, session.Id)
	if err != nil {
		t.Fatalf("failed to get deleted website session: %v", err)
	}
	if got != nil {
		t.Fatalf("expected deleted website session to be missing, got %+v", got)
	}
}
