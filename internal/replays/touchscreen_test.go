package replays

import (
	"testing"

	"github.com/osuTitanic/titanic/internal/constants"
)

func TestTouchscreenUsage(t *testing.T) {
	dir, err := testdata.ReadDir("testdata")
	if err != nil {
		t.Skipf("can't open replays: %v", err)
	}

	for _, entry := range dir {
		t.Run(entry.Name(), func(t *testing.T) {
			score, frames, _ := replayFixtureDeserialized(t, entry.Name())
			if score.Mode != constants.ModeOsu {
				t.Skip("not an osu standard replay")
			}

			detected, confidence := DetectTouchscreenUsage(frames, 0.5)
			t.Logf("detected: %t, score: %f", detected, confidence)

			// TODO: Add assertions
		})
	}
}
