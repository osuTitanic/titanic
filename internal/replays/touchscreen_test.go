package replays

import (
	"slices"
	"testing"

	"github.com/osuTitanic/titanic/internal/constants"
)

var touchscreenUsers = []string{
	"skita",
	"wyoming",
}

func TestTouchscreenUsage(t *testing.T) {
	dir, err := testdata.ReadDir("testdata")
	if err != nil {
		t.Skipf("can't open replays: %v", err)
	}

	for _, entry := range dir {
		t.Run(entry.Name(), func(t *testing.T) {
			s, frames, _ := replayFixtureDeserialized(t, entry.Name())
			if s.Mode != constants.ModeOsu {
				t.Skip("not an osu standard replay")
			}

			detected, score := DetectTouchscreenUsage(frames, 0.45)
			t.Logf("detected: %t, score: %f", detected, score)

			isTouchscreenUser := slices.Contains(touchscreenUsers, s.User.Name)
			if detected != isTouchscreenUser {
				t.Errorf("detected: %t, expected: %t", detected, isTouchscreenUser)
			}
		})
	}
}
