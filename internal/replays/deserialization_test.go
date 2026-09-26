package replays

import (
	"embed"
	"testing"

	"github.com/osuTitanic/titanic/internal/schemas"
)

//go:embed testdata/*.osr
var testdata embed.FS

func TestReplayDeserialization(t *testing.T) {
	dir, err := testdata.ReadDir("testdata")
	if err != nil {
		t.Skipf("can't open replays: %v", err)
	}

	for _, entry := range dir {
		t.Run(entry.Name(), func(t *testing.T) {
			_, frames, seed := replayFixtureDeserialized(t, entry.Name())
			t.Logf("frames: %d, seed: %d", len(frames), seed)
		})
	}
}

func replayFixture(t *testing.T, name string) []byte {
	t.Helper()

	data, err := testdata.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatalf("read replay fixture %q: %v", name, err)
	}
	return data
}

func replayFixtureDeserialized(t *testing.T, name string) (*schemas.Score, []Frame, int64) {
	t.Helper()

	score, frames, seed, err := Deserialize(replayFixture(t, name))
	if err != nil {
		t.Fatalf("deserialize replay fixture %q: %v", name, err)
	}
	return score, frames, seed
}
