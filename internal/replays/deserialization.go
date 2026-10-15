package replays

import (
	"bufio"
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

const (
	maxReplaySize      = 50 << 20 // TODO: this is pretty large, maybe lower this
	maxReplayFrames    = 2_500_000
	maxReplayFrameSize = 1 << 10 // 1 KiB
)

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
	lzmaReader, err := lzmaConfig.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, 0, fmt.Errorf("open replay stream: %w", err)
	}

	reader := &io.LimitedReader{
		R: lzmaReader,
		N: maxReplaySize + 1,
	}
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, maxReplayFrameSize+1), maxReplayFrameSize+1) // allocate ~1 KB buffer for frame

	// Use a scanner to split the replay into frames
	// We could strings.Split(...) the frames in one go, but
	// a scanner is more suitable for really large replay files
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if index := bytes.IndexByte(data, ','); index >= 0 {
			// Advance to the next byte & return the frame
			return index + 1, data[:index], nil
		}
		if atEOF && len(data) > 0 {
			// Return the remaining data as the last frame
			return len(data), data, nil
		}
		return 0, nil, nil
	})

	frames = make([]Frame, 0, 4096)
	currentTime := 0
	index := 0

	for scanner.Scan() {
		if index >= maxReplayFrames+1 {
			return nil, 0, fmt.Errorf("replay contains too many frames: maximum is %d", maxReplayFrames)
		}
		frameIndex := index
		index++

		rawFrame := scanner.Text()
		if rawFrame == "" {
			continue
		}

		fields := strings.Split(rawFrame, "|")
		if len(fields) != 4 {
			return nil, 0, fmt.Errorf("invalid replay frame %d: got %d fields, want 4", frameIndex, len(fields))
		}
		if fields[0] == "-12345" {
			seed, err = strconv.ParseInt(fields[3], 10, 64)
			if err != nil {
				return nil, 0, fmt.Errorf("invalid replay seed %q: %w", fields[3], err)
			}
			continue
		}
		if len(frames) >= maxReplayFrames {
			return nil, 0, fmt.Errorf("replay contains too many frames: maximum is %d", maxReplayFrames)
		}

		delta, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, 0, fmt.Errorf("invalid replay frame %d delta %q: %w", frameIndex, fields[0], err)
		}
		// TODO: delta should never be negative, right?

		x, err := strconv.ParseFloat(fields[1], 64)
		if err != nil || math.IsNaN(x) || math.IsInf(x, 0) {
			return nil, 0, fmt.Errorf("invalid replay frame %d x-coordinate %q", frameIndex, fields[1])
		}
		y, err := strconv.ParseFloat(fields[2], 64)
		if err != nil || math.IsNaN(y) || math.IsInf(y, 0) {
			return nil, 0, fmt.Errorf("invalid replay frame %d y-coordinate %q", frameIndex, fields[2])
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
	if err := scanner.Err(); err != nil {
		return nil, 0, fmt.Errorf("scan replay frames: %w", err)
	}
	if reader.N == 0 {
		// N = 0 indicates the decompressed data exceeds maxReplaySize
		return nil, 0, fmt.Errorf("decompress replay: data exceeds %d bytes", maxReplaySize)
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
