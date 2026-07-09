//go:build integration

package rankings_test

import (
	"testing"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/rankings"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/testkit"
)

func TestRankingsServiceUpdate(t *testing.T) {
	service := rankings.NewRankingsService(testkit.RedisClient(t))
	mode := constants.ModeOsu

	entries := []struct {
		stats   *schemas.Stats
		country string
	}{
		{
			stats: &schemas.Stats{
				UserId: 1,
				Mode:   mode,
				PP:     100,
				Rscore: 1000,
				Tscore: 10000,
				Acc:    97.5,
				CountA: 2,
			},
			country: "US",
		},
		{
			stats: &schemas.Stats{
				UserId: 2,
				Mode:   mode,
				PP:     250,
				Rscore: 3000,
				Tscore: 30000,
				Acc:    99.1,
				CountX: 3,
			},
			country: "US",
		},
		{
			stats: &schemas.Stats{
				UserId: 3,
				Mode:   mode,
				PP:     175,
				Rscore: 2000,
				Tscore: 20000,
				Acc:    98.2,
				CountS: 1,
			},
			country: "CA",
		},
	}

	for _, entry := range entries {
		if err := service.Update(entry.stats, entry.country); err != nil {
			t.Fatalf("Update(user %d) error = %v", entry.stats.UserId, err)
		}
	}

	assertRank(t, service, 2, mode, 1)
	assertRank(t, service, 3, mode, 2)
	assertRank(t, service, 1, mode, 3)
	assertCountryRank(t, service, 2, mode, "US", 1)
	assertCountryRank(t, service, 1, mode, "US", 2)
	assertCountryRank(t, service, 3, mode, "US", 0)

	score, err := service.Score(2, mode)
	if err != nil {
		t.Fatalf("Score() error = %v", err)
	}
	if score != 3000 {
		t.Fatalf("Score() = %d, want 3000", score)
	}

	count, err := service.PlayerCount(mode, "performance", nil)
	if err != nil {
		t.Fatalf("PlayerCount() error = %v", err)
	}
	if count != 3 {
		t.Fatalf("PlayerCount() = %d, want 3", count)
	}

	top, err := service.TopPlayers(mode, 0, 2, "performance", nil)
	if err != nil {
		t.Fatalf("TopPlayers() error = %v", err)
	}
	if len(top) != 2 {
		t.Fatalf("TopPlayers() returned %d players, want 2", len(top))
	}
	if top[0].UserId != 2 || top[1].UserId != 3 {
		t.Fatalf("TopPlayers() ids = [%d, %d], want [2, 3]", top[0].UserId, top[1].UserId)
	}

	if err := service.Remove(2, "US"); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}
	assertRank(t, service, 2, mode, 0)
	assertRank(t, service, 3, mode, 1)
	assertCountryRank(t, service, 1, mode, "US", 1)
}

func assertRank(t *testing.T, service *rankings.RankingsService, userId int, mode constants.Mode, want int) {
	t.Helper()

	got, err := service.GlobalRank(userId, mode)
	if err != nil {
		t.Fatalf("GlobalRank(%d) error = %v", userId, err)
	}
	if got != want {
		t.Fatalf("GlobalRank(%d) = %d, want %d", userId, got, want)
	}
}

func assertCountryRank(t *testing.T, service *rankings.RankingsService, userId int, mode constants.Mode, country string, want int) {
	t.Helper()

	got, err := service.CountryRank(userId, mode, country)
	if err != nil {
		t.Fatalf("CountryRank(%d, %q) error = %v", userId, country, err)
	}
	if got != want {
		t.Fatalf("CountryRank(%d, %q) = %d, want %d", userId, country, got, want)
	}
}
