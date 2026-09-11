package scoring

import (
	"cmp"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

const (
	MultipartMemory       = 1 << 20
	MaxReplaySize         = 10 << 20
	MaxSubmissionBodySize = MaxReplaySize + 1<<20
)

func ResolveSubmissionContext(ctx *server.Context, endpoint Endpoint) (*SubmissionContext, error) {
	userAgent := cmp.Or(
		ctx.Request.Header.Get("user-agent"),
		"osu!", // TODO: forgot if it makes sense to actually use a fallback here
	)
	if !constants.OsuUserAgent.MatchString(userAgent) {
		return nil, fmt.Errorf("score submission: invalid user agent %q", userAgent)
	}

	err := ctx.EnsureMultipartForm(
		MaxSubmissionBodySize,
		MultipartMemory,
	)
	if err != nil && !errors.Is(err, http.ErrNotMultipart) {
		return nil, fmt.Errorf("parse submission form: %w", err)
	}

	scoreData := endpoint.RequestValue(ctx.Request, "score")
	if scoreData == "" {
		return nil, fmt.Errorf("parse submission: missing score data")
	}

	replay, err := readReplay(ctx.Request)
	if err != nil {
		return nil, err
	}

	ivEncoded := ctx.PostFormValue("iv")
	clientHash := ctx.PostFormValue("s")
	processes := ctx.PostFormValue("pl")
	funSpoiler := ctx.PostFormValue("fs")
	osuVersion := ctx.PostFormValue("osuver")

	if endpoint.UsesEncryptedFormScore() && ivEncoded != "" {
		iv, decodeErr := base64.StdEncoding.DecodeString(ivEncoded)
		if decodeErr != nil {
			return nil, fmt.Errorf("decode submission IV: %w", decodeErr)
		}

		key := encryptionKeyV1
		if osuVersion != "" {
			// yummy osu score burger
			key = encryptionKeyV2 + osuVersion
		}

		scoreData, err = DecryptString(scoreData, iv, key)
		if err != nil {
			return nil, fmt.Errorf("decrypt score: %w", err)
		}
		if clientHash != "" {
			clientHash, err = DecryptString(clientHash, iv, key)
			if err != nil {
				return nil, fmt.Errorf("decrypt client hash: %w", err)
			}
		}
		if funSpoiler != "" {
			funSpoiler, err = DecryptString(funSpoiler, iv, key)
			if err != nil {
				return nil, fmt.Errorf("decrypt fun spoiler: %w", err)
			}
		}
		if processes != "" {
			processes, err = DecryptString(processes, iv, key)
			if err != nil {
				return nil, fmt.Errorf("decrypt process list: %w", err)
			}
		}
	}

	submission, err := parseScoreData(scoreData, endpoint)
	if err != nil {
		return nil, err
	}

	submission.Replay = replay
	submission.FunSpoiler = funSpoiler
	submission.ClientHash = clientHash
	submission.Processes = processes

	if replay != nil {
		replayMd5 := fmt.Sprintf("%x", md5.Sum(replay))
		submission.ReplayMd5 = new(replayMd5)
	}

	failtimeValue := endpoint.RequestValue(ctx.Request, "ft")
	exitedValue := endpoint.RequestValue(ctx.Request, "x")
	submission.Exited = exitedValue != ""

	if failtimeValue != "" {
		failtime, parseErr := strconv.Atoi(failtimeValue)
		if parseErr != nil {
			return nil, fmt.Errorf("parse fail time %q: %w", failtimeValue, parseErr)
		}
		submission.Failtime = new(failtime)
	}
	if submission.Failtime == nil {
		submission.Failtime = new(0)
	}
	if submission.Passed {
		submission.Failtime = nil
		submission.Exited = false
	}

	submission.Acc = submission.Accuracy()
	// TODO: Grade calc
	return submission, nil
}

func readReplay(request *http.Request) ([]byte, error) {
	if request.MultipartForm == nil {
		return nil, nil
	}

	files := request.MultipartForm.File["score"]
	if len(files) == 0 {
		return nil, nil
	}

	// NOTE: The form data can contain two "score" sections, where one
	// 		 of them is the score data, and the other is the replay
	header := files[len(files)-1]

	if header.Filename != "replay" && header.Filename != "score" {
		return nil, fmt.Errorf("read replay: invalid replay filename %q", header.Filename)
	}
	if header.Size > MaxReplaySize {
		return nil, fmt.Errorf("read replay: replay is too large: %d bytes", header.Size)
	}

	file, err := header.Open()
	if err != nil {
		return nil, fmt.Errorf("open replay: %w", err)
	}
	defer file.Close()

	replay, err := io.ReadAll(io.LimitReader(file, MaxReplaySize+1))
	if err != nil {
		return nil, fmt.Errorf("read replay: %w", err)
	}
	if len(replay) > MaxReplaySize {
		return nil, fmt.Errorf("read replay: replay is too large: %d bytes", len(replay))
	}
	return replay, nil
}

func parseScoreData(scoreData string, endpoint Endpoint) (*SubmissionContext, error) {
	fields := strings.Split(scoreData, ":")
	if len(fields) < 15 {
		return nil, fmt.Errorf("parse score: got %d fields, want at least 15", len(fields))
	}

	count300, err := resolveInt(fields, 3, "300 count")
	if err != nil {
		return nil, err
	}
	count100, err := resolveInt(fields, 4, "100 count")
	if err != nil {
		return nil, err
	}
	count50, err := resolveInt(fields, 5, "50 count")
	if err != nil {
		return nil, err
	}
	countGeki, err := resolveInt(fields, 6, "geki count")
	if err != nil {
		return nil, err
	}
	countKatu, err := resolveInt(fields, 7, "katu count")
	if err != nil {
		return nil, err
	}
	countMiss, err := resolveInt(fields, 8, "miss count")
	if err != nil {
		return nil, err
	}
	totalScore, err := resolveInt64(fields, 9, "total score")
	if err != nil {
		return nil, err
	}
	maxCombo, err := resolveInt(fields, 10, "max combo")
	if err != nil {
		return nil, err
	}
	perfect, err := resolveBool(fields, 11, "perfect")
	if err != nil {
		return nil, err
	}
	grade, err := resolveEnum[constants.Grade](fields, 12, "grade")
	if err != nil {
		return nil, err
	}
	modsValue, err := resolveUint(fields, 13, "mods")
	if err != nil {
		return nil, err
	}
	passed, err := resolveBool(fields, 14, "passed")
	if err != nil {
		return nil, err
	}

	mode, err := resolveEnum[constants.Mode](fields, 15, "mode")
	if err != nil {
		return nil, err
	}

	versionField := fieldAt(fields, 17)
	clientVersion := 0

	flagsRaw := strings.Count(versionField, " ")
	flags := constants.IntegrityFlags(flagsRaw)

	if versionField != "" {
		clientVersion, err = strconv.Atoi(strings.TrimSpace(versionField))
		if err != nil {
			return nil, fmt.Errorf("parse score client version %q: %w", versionField, err)
		}
	}

	score := &schemas.Score{
		ClientVersion: clientVersion,
		Checksum:      fields[2],
		Mode:          mode,
		TotalScore:    totalScore,
		MaxCombo:      maxCombo,
		Mods:          constants.Mods(modsValue),
		Perfect:       perfect,
		Count300:      count300,
		Count100:      count100,
		Count50:       count50,
		CountMiss:     countMiss,
		CountGeki:     countGeki,
		CountKatu:     countKatu,
		Grade:         grade,
	}
	submission := NewSubmissionContext(endpoint, score)
	submission.Username = strings.TrimSpace(fields[1])
	submission.BeatmapChecksum = fields[0]
	submission.Flags = flags
	submission.Passed = passed
	return submission, nil
}

func fieldAt(fields []string, index int) string {
	if index >= len(fields) {
		return ""
	}
	return fields[index]
}

func resolveInt(fields []string, index int, name string) (int, error) {
	value, err := strconv.Atoi(fields[index])
	if err != nil {
		return 0, fmt.Errorf("parse score %s %q: %w", name, fields[index], err)
	}
	return value, nil
}

func resolveUint(fields []string, index int, name string) (uint, error) {
	value, err := strconv.ParseUint(fields[index], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("parse score %s %q: %w", name, fields[index], err)
	}
	return uint(value), nil
}

func resolveInt64(fields []string, index int, name string) (int64, error) {
	value, err := strconv.ParseInt(fields[index], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse score %s %q: %w", name, fields[index], err)
	}
	return value, nil
}

func resolveBool(fields []string, index int, name string) (bool, error) {
	switch {
	case strings.EqualFold(fields[index], "true"):
		return true, nil
	case strings.EqualFold(fields[index], "false"):
		return false, nil
	case strings.EqualFold(fields[index], "True"):
		return true, nil
	case strings.EqualFold(fields[index], "False"):
		return false, nil
	case strings.EqualFold(fields[index], "1"):
		return true, nil
	case strings.EqualFold(fields[index], "0"):
		return false, nil
	default:
		return false, fmt.Errorf("parse score: invalid %s value %q", name, fields[index])
	}
}

func resolveEnum[T constants.HasValidator](fields []string, index int, name string) (T, error) {
	field := fieldAt(fields, index)
	if field == "" {
		var zero T
		return zero, nil
	}

	value, err := constants.ResolveEnum[T](field)
	if err != nil {
		return value, fmt.Errorf("parse score %s %q: %w", name, field, err)
	}
	if !value.Valid() {
		return value, fmt.Errorf("parse score: invalid %s value %q", name, field)
	}
	return value, nil
}
