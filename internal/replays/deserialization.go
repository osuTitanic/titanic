package replays

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
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
	Left1    ButtonState = 1
	Right1   ButtonState = 2
	Left2    ButtonState = 4
	Right2   ButtonState = 8
	Smoke    ButtonState = 16
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

func Deserialize(data []byte) (*schemas.Score, []Frame, int64, error) {
	reader := bytes.NewReader(data)

	mode, err := readU8(reader)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read mode: %w", err)
	}

	clientVersion, err := readS32(reader)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read client version: %w", err)
	}

	beatmapChecksum, err := readString(reader)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read beatmap checksum: %w", err)
	}

	username, err := readString(reader)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read username: %w", err)
	}

	scoreChecksum, err := readString(reader)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read score checksum: %w", err)
	}

	count300, err := readU16(reader)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read count300: %w", err)
	}

	count100, err := readU16(reader)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read count100: %w", err)
	}

	count50, err := readU16(reader)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read count50: %w", err)
	}

	countGeki, err := readU16(reader)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read countGeki: %w", err)
	}

	countKatu, err := readU16(reader)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read countKatu: %w", err)
	}

	countMiss, err := readU16(reader)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read countMiss: %w", err)
	}

	totalScore, err := readU32(reader)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read total score: %w", err)
	}

	maxCombo, err := readU16(reader)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read max combo: %w", err)
	}

	perfect, err := readBool(reader)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read perfect: %w", err)
	}

	mods, err := readU32(reader)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read mods: %w", err)
	}

	if _, err := readString(reader); err != nil {
		return nil, nil, 0, fmt.Errorf("read hp graph: %w", err)
	}

	timestamp, err := readS64(reader)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read timestamp: %w", err)
	}

	replayLength, err := readU32(reader)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read replay length: %w", err)
	}

	if int(replayLength) > reader.Len() {
		return nil, nil, 0, fmt.Errorf(
			"invalid replay length: %d bytes requested, %d remaining",
			replayLength, reader.Len(),
		)
	}

	replayData := make([]byte, replayLength)
	if _, err := io.ReadFull(reader, replayData); err != nil {
		return nil, nil, 0, fmt.Errorf("read replay data: %w", err)
	}

	var scoreId uint64

	if clientVersion >= 20140721 {
		scoreId, err = readU64(reader)
	} else {
		var id uint32
		id, err = readU32(reader)
		scoreId = uint64(id)
	}
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read score id: %w", err)
	}

	frames, seed, err := DeserializeFrames(replayData)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("deserialize replay frames: %w", err)
	}

	score := &schemas.Score{
		Id:            int64(scoreId),
		ClientVersion: int(clientVersion),
		Checksum:      scoreChecksum,
		Mode:          constants.Mode(mode),
		Count300:      int(count300),
		Count100:      int(count100),
		Count50:       int(count50),
		CountGeki:     int(countGeki),
		CountKatu:     int(countKatu),
		CountMiss:     int(countMiss),
		TotalScore:    int64(totalScore),
		MaxCombo:      int(maxCombo),
		Perfect:       perfect,
		Mods:          constants.Mods(mods),
		SubmittedAt:   timeFromTicks(timestamp),
		Beatmap: &schemas.Beatmap{
			Checksum: beatmapChecksum,
		},
		User: &schemas.User{
			Name: username,
		},
	}
	return score, frames, seed, nil
}
