package performance

import (
	"bytes"
	"errors"
	"io"

	_ "embed"
)

//go:embed beatmap_disco_prince.osu
var beatmapDataDiscoPrince []byte

//go:embed beatmap_freedom_dive.osu
var beatmapDataFreedomDive []byte

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
