package bss

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
 * This file contains staging logic used only by pre-osz2 beatmap submission.
 * This is not relevant for the modern osz2 submission system.
 *
 * Pre-osz2 refers to the system used before the introduction of osz2 in 2011,
 * which came with a massive overhaul of the beatmap submission system itself.
 * Before this overhaul, the system had a stateful design:
 *
 * getid[2-5] receives one .osu file per request, which we stage in redis to be used later.
 * Once all .osu files are staged, the last request uses the staged submission to update (or create)
 * the beatmapset & all of its difficulties in a single transaction.
 */

// StagedSubmission contains the state that is shared
// between the individual pre-osz2 getid requests.
type StagedSubmission struct {
	UserId       int `json:"user_id"`
	BeatmapsetId int `json:"beatmapset_id"`

	// OszFilename is returned by getid, which the client will use
	// as a filename for the osz file in upload.php
	OszFilename string `json:"osz_filename,omitempty"`

	NewBeatmapset bool `json:"new_beatmapset"`
	HasVideo      bool `json:"has_video"`
	HasStoryboard bool `json:"has_storyboard"`

	Beatmaps []StagedSubmissionBeatmap `json:"beatmaps"`
}

// "Binary" encoding used for redis .Scan(...) & .Set(...)

func (submission StagedSubmission) MarshalBinary() ([]byte, error) {
	return json.Marshal(submission)
}

func (submission *StagedSubmission) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, submission)
}

// StagedSubmissionBeatmap contains an assigned beatmap ID and
// the .osu file staged by a pre-osz2 getid request.
type StagedSubmissionBeatmap struct {
	Id       int    `json:"id"`
	Filename string `json:"filename,omitempty"`
	Checksum string `json:"checksum,omitempty"`
	Contents []byte `json:"contents,omitempty"`
}

// StagedSubmissionStore manages the current staged submission for each user.
type StagedSubmissionStore struct {
	Redis  redis.Cmdable
	Prefix string        // redis key prefix
	TTL    time.Duration // how long the staged files will be kept
}

func NewStagedSubmissionStore(client redis.Cmdable) *StagedSubmissionStore {
	return &StagedSubmissionStore{
		Redis:  client,
		Prefix: "bss:staging:",
		TTL:    time.Hour,
	}
}

func (store *StagedSubmissionStore) SubmissionKey(userId int) string {
	return store.Prefix + strconv.Itoa(userId)
}

func (store *StagedSubmissionStore) Save(ctx context.Context, submission *StagedSubmission) error {
	if submission == nil {
		return errors.New("bss: nil staged submission")
	}
	if submission.UserId <= 0 || submission.BeatmapsetId <= 0 {
		return errors.New("bss: invalid staged submission")
	}

	ttl := store.TTL
	if ttl <= 0 {
		ttl = time.Hour
	}

	if err := store.Redis.Set(
		ctx,
		store.SubmissionKey(submission.UserId),
		submission,
		ttl,
	).Err(); err != nil {
		return fmt.Errorf("bss: save staged submission: %w", err)
	}
	return nil
}

func (store *StagedSubmissionStore) Get(ctx context.Context, userId int) (*StagedSubmission, error) {
	if userId <= 0 {
		return nil, errors.New("bss: invalid staged submission")
	}
	var submission StagedSubmission

	err := store.Redis.Get(ctx, store.SubmissionKey(userId)).Scan(&submission)
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("bss: load staged submission: %w", err)
	}

	if submission.UserId != userId || submission.BeatmapsetId <= 0 {
		return nil, errors.New("bss: invalid staged submission")
	}
	return &submission, nil
}

func (store *StagedSubmissionStore) Delete(ctx context.Context, userId int) error {
	if userId <= 0 {
		return errors.New("bss: invalid staged submission")
	}

	if err := store.Redis.Del(ctx, store.SubmissionKey(userId)).Err(); err != nil {
		return fmt.Errorf("bss: delete staged submission: %w", err)
	}
	return nil
}
