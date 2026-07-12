//go:build rosu

package performance

import (
	"bytes"
	"errors"
	"io"
	"math"
	"testing"

	_ "embed"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
)

//go:embed beatmap_disco_prince.osu
var beatmapDataDiscoPrince []byte

//go:embed beatmap_freedom_dive.osu
var beatmapDataFreedomDive []byte

func TestPPv2ServiceRosu(t *testing.T) {
	provider := &testProvider{}
	service, err := NewPPv2Service(provider)
	if err != nil {
		t.Fatalf("initializing ppv2: %v", err)
	}

	difficulty, err := service.CalculateDifficulty(75, constants.ModeOsu, constants.NoMod)
	if err != nil {
		t.Fatalf("calculate difficulty: %v", err)
	}
	if !provider.closed {
		t.Fatal("beatmap stream was not closed")
	}
	if difficulty.StarRating <= 0 || math.IsNaN(difficulty.StarRating) {
		t.Fatalf("unexpected star rating: %f", difficulty.StarRating)
	}
	provider.closed = false

	t.Logf("calculated stars: %f", difficulty.StarRating)

	pp, err := service.CalculatePerformance(&schemas.Score{
		BeatmapId:  75,
		Mode:       constants.ModeOsu,
		MaxCombo:   314,
		Count300:   194,
		TotalScore: 1_491_676,
	})
	if err != nil {
		t.Fatalf("calculate performance: %v", err)
	}
	if !provider.closed {
		t.Fatal("beatmap stream was not closed")
	}
	if pp <= 0 || math.IsNaN(pp) {
		t.Fatalf("unexpected performance: %f", pp)
	}

	t.Logf("calculated pp: %f", pp)
}

type testProvider struct {
	closed bool
}

func (*testProvider) Setup() error {
	return nil
}

func (*testProvider) Osz(int, bool) (io.ReadCloser, int64, error) {
	return nil, 0, errors.New("not implemented")
}

func (provider *testProvider) Osu(int) (io.ReadCloser, error) {
	return &trackedReadCloser{
		// TODO: Test both disco prince & freedom dive
		Reader: bytes.NewReader(beatmapDataDiscoPrince),
		close:  func() { provider.closed = true },
	}, nil
}

func (*testProvider) Preview(int) (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}

func (*testProvider) Background(int, bool) (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}

type trackedReadCloser struct {
	io.Reader
	close func()
}

func (stream *trackedReadCloser) Close() error {
	stream.close()
	return nil
}
