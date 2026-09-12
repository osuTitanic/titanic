package replays

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/ulikunitz/xz/lzma"
)

const maxReplaySize = 50 << 20 // TODO: this is pretty large, maybe lower this

type Frame struct {
	Delta   int
	Time    int
	X       float64
	Y       float64
	Buttons ButtonState
}

type ButtonState int

const (
	NoButton ButtonState = 0
	Left1    ButtonState = 1 << iota
	Right1
	Left2
	Right2
	Smoke
)

const allButtons = Left1 | Right1 | Left2 | Right2 | Smoke

func (buttons ButtonState) Valid() bool {
	return buttons&^allButtons == 0
}

func (buttons ButtonState) Has(button ButtonState) bool {
	return buttons&button != 0
}

var lzmaConfig = lzma.ReaderConfig{DictCap: maxReplaySize}

func DeserializeFrames(data []byte) (frames []Frame, seed int64, err error) {
	reader, err := lzmaConfig.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, 0, fmt.Errorf("open replay stream: %w", err)
	}

	decoded, err := io.ReadAll(io.LimitReader(reader, maxReplaySize+1))
	if err != nil {
		return nil, 0, fmt.Errorf("decompress replay: %w", err)
	}
	if len(decoded) > maxReplaySize {
		return nil, 0, fmt.Errorf("decompress replay: data exceeds %d bytes", maxReplaySize)
	}

	rawFrames := strings.Split(string(decoded), ",")
	frames = make([]Frame, 0, len(rawFrames))
	currentTime := 0

	for index, rawFrame := range rawFrames {
		if rawFrame == "" {
			continue
		}

		fields := strings.Split(rawFrame, "|")
		if len(fields) != 4 {
			return nil, 0, fmt.Errorf("invalid replay frame %d: got %d fields, want 4", index, len(fields))
		}
		if fields[0] == "-12345" {
			seed, err = strconv.ParseInt(fields[3], 10, 64)
			if err != nil {
				return nil, 0, fmt.Errorf("invalid replay seed %q: %w", fields[3], err)
			}
			continue
		}

		delta, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, 0, fmt.Errorf("invalid replay frame %d delta %q: %w", index, fields[0], err)
		}
		// TODO: delta should never be negative, right?

		x, err := strconv.ParseFloat(fields[1], 64)
		if err != nil || math.IsNaN(x) || math.IsInf(x, 0) {
			return nil, 0, fmt.Errorf("invalid replay frame %d x-coordinate %q", index, fields[1])
		}
		y, err := strconv.ParseFloat(fields[2], 64)
		if err != nil || math.IsNaN(y) || math.IsInf(y, 0) {
			return nil, 0, fmt.Errorf("invalid replay frame %d y-coordinate %q", index, fields[2])
		}

		buttons, err := strconv.Atoi(fields[3])
		if err != nil {
			return nil, 0, fmt.Errorf("invalid buttons %q: %w", fields[3], err)
		}
		buttonState := ButtonState(buttons)
		if !buttonState.Valid() {
			return nil, 0, fmt.Errorf("invalid button state: %d", buttons)
		}

		// Convert delta time into absolute replay time
		currentTime = currentTime + delta

		frames = append(frames, Frame{
			Time:    currentTime,
			Delta:   delta,
			X:       x,
			Y:       y,
			Buttons: buttonState,
		})
	}
	return frames, seed, nil
}
