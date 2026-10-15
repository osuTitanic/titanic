package bss

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrInvalidSubmissionTicket = errors.New("bss: invalid submission ticket")

type SubmissionTicketKind string

const (
	SubmissionTicketOsz     SubmissionTicketKind = "osz"
	SubmissionTicketBeatmap SubmissionTicketKind = "beatmap"
)

func (kind SubmissionTicketKind) Valid() bool {
	return kind == SubmissionTicketOsz || kind == SubmissionTicketBeatmap
}

// SubmissionTicket authorizes a resource upload
// It gets sent with a getid response and is used
// in the upload endpoint as the "c" parameter
// to identify the user and resource to upload.
type SubmissionTicket struct {
	UserId       int                  `json:"u"`
	BeatmapsetId int                  `json:"s"`
	Kind         SubmissionTicketKind `json:"k"`
	Filename     string               `json:"f"`
	ExpiresAt    int64                `json:"e"`
}

func (ticket SubmissionTicket) Validate() error {
	if ticket.UserId <= 0 || ticket.BeatmapsetId <= 0 {
		return ErrInvalidSubmissionTicket
	}
	if !ticket.Kind.Valid() {
		return ErrInvalidSubmissionTicket
	}
	if ticket.Filename == "" {
		return ErrInvalidSubmissionTicket
	}
	if ticket.ExpiresAt <= time.Now().Unix() {
		return ErrInvalidSubmissionTicket
	}
	return nil
}

func (ticket SubmissionTicket) Encode(secret string) (string, error) {
	payload, err := json.Marshal(ticket)
	if err != nil {
		return "", fmt.Errorf("bss: encode submission ticket: %w", err)
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payload)
	return payloadEncoded + "." + SignSubmissionTicket(secret, payloadEncoded), nil
}

func NewSubmissionTicket(secret string, userId, beatmapsetId int, kind SubmissionTicketKind, filename string) (string, error) {
	ticket := SubmissionTicket{
		UserId:       userId,
		BeatmapsetId: beatmapsetId,
		Kind:         kind,
		Filename:     filename,
		ExpiresAt:    time.Now().Add(time.Hour).Unix(),
	}
	if err := ticket.Validate(); err != nil {
		return "", err
	}
	return ticket.Encode(secret)
}

func ParseSubmissionTicket(secret, value string) (*SubmissionTicket, error) {
	if secret == "" {
		return nil, errors.New("bss: submission ticket secret is empty")
	}

	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return nil, ErrInvalidSubmissionTicket
	}

	payloadEncoded := parts[0]
	signature := []byte(parts[1])
	expected := SignSubmissionTicket(secret, payloadEncoded)

	if !hmac.Equal([]byte(expected), signature) {
		return nil, ErrInvalidSubmissionTicket
	}

	payload, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return nil, ErrInvalidSubmissionTicket
	}

	var ticket SubmissionTicket
	if err := json.Unmarshal(payload, &ticket); err != nil {
		return nil, ErrInvalidSubmissionTicket
	}
	if err := ticket.Validate(); err != nil {
		return nil, err
	}
	return &ticket, nil
}

func SignSubmissionTicket(secret, payload string) string {
	signature := hmac.New(sha256.New, []byte("bss:"+secret))
	signature.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(signature.Sum(nil))
}
