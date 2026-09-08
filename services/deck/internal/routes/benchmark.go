package routes

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

const (
	benchmarkMaxFramerate = 1_000_000
	benchmarkMaxScore     = 1_000_000_000
)

type BenchmarkRequest struct {
	Smoothness float64
	Framerate  int
	RawScore   int64
	Client     string
	Hardware   BenchmarkHardware
}

func (r *BenchmarkRequest) Validate() error {
	if math.IsNaN(r.Smoothness) || math.IsInf(r.Smoothness, 0) {
		return fmt.Errorf("smoothness must be a valid number: %v", r.Smoothness)
	}
	if r.Smoothness < 0 || r.Smoothness > 100 {
		return fmt.Errorf("smoothness must be between 0 and 100: %v", r.Smoothness)
	}
	if r.Framerate <= 0 || r.Framerate > benchmarkMaxFramerate {
		return fmt.Errorf("framerate must be between 1 and %d: %d", benchmarkMaxFramerate, r.Framerate)
	}
	if r.RawScore <= 0 || r.RawScore > benchmarkMaxScore {
		return fmt.Errorf("raw score must be between 1 and %d: %d", benchmarkMaxScore, r.RawScore)
	}
	if err := r.Hardware.Validate(); err != nil {
		return fmt.Errorf("validate hardware: %w", err)
	}
	return nil
}

func (r BenchmarkRequest) Grade() string {
	switch {
	case r.Smoothness == 100:
		return "SS"
	case r.Smoothness > 95:
		return "S"
	case r.Smoothness > 90:
		return "A"
	case r.Smoothness > 80:
		return "B"
	case r.Smoothness > 70:
		return "C"
	default:
		return "D"
	}
}

func NewBenchmarkRequest(ctx *server.Context) (request BenchmarkRequest, err error) {
	smoothness, err := ctx.FormValueFloat("s")
	if err != nil {
		return request, fmt.Errorf("parse smoothness: %w", err)
	}
	framerate, err := ctx.FormValueInt("f")
	if err != nil {
		return request, fmt.Errorf("parse framerate: %w", err)
	}
	rawScore, err := ctx.FormValueInt64("r")
	if err != nil {
		return request, fmt.Errorf("parse raw score: %w", err)
	}

	var hardware BenchmarkHardware
	var hardwareRaw = ctx.FormValue("h")

	decoder := json.NewDecoder(strings.NewReader(hardwareRaw))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&hardware); err != nil {
		return request, fmt.Errorf("decode hardware: %w", err)
	}

	return BenchmarkRequest{
		Smoothness: smoothness,
		Framerate:  framerate,
		RawScore:   rawScore,
		Client:     ctx.FormValue("c"),
		Hardware:   hardware,
	}, nil
}

type BenchmarkHardware struct {
	Renderer           string  `json:"renderer"`
	Resolution         *string `json:"resolution,omitempty"`
	Fullscreen         *bool   `json:"fullscreen,omitempty"`
	Letterboxing       *bool   `json:"letterboxing,omitempty"`
	DotNetVersion      *string `json:"dotnet_version,omitempty"`
	ClientArchitecture *string `json:"client_architecture,omitempty"`

	CPU                     *string `json:"cpu,omitempty"`
	Cores                   *int64  `json:"cores,omitempty"`
	Threads                 *int64  `json:"threads,omitempty"`
	GPU                     *string `json:"gpu,omitempty"`
	RAM                     *int64  `json:"ram,omitempty"`
	OS                      *string `json:"os,omitempty"`
	MotherboardManufacturer *string `json:"motherboard_manufacturer,omitempty"`
	Motherboard             *string `json:"motherboard,omitempty"`
}

type Validator[T any] func(T) error

var benchmarkHardwareValidators = []Validator[*BenchmarkHardware]{
	(*BenchmarkHardware).ValidateRenderer,
	(*BenchmarkHardware).ValidateResolution,
	(*BenchmarkHardware).ValidateArchitecture,
	(*BenchmarkHardware).ValidateDotNetVersion,
	(*BenchmarkHardware).ValidateFullHardware,
}

func (h *BenchmarkHardware) Validate() error {
	for _, validator := range benchmarkHardwareValidators {
		if err := validator(h); err != nil {
			return err
		}
	}
	return nil
}

// /web/osu-benchmark.php -> Submit an osu! client benchmark
//
// This custom endpoint is not used by the official osu! client.
// It was added as an easter egg in modified clients to revive the old benchmark feature.
func Benchmark(ctx *server.Context) {
	request, err := NewBenchmarkRequest(ctx)
	if err != nil {
		ctx.Logger.Warn("Failed to decode benchmark request", "error", err)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := request.Validate(); err != nil {
		ctx.Logger.Warn("Invalid benchmark request", "error", err)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	user, ok := ctx.HandleUserAuthenticationFormSimple("u", "p", true)
	if !ok {
		return
	}
	if !user.Activated || user.Restricted {
		ctx.Response.WriteHeader(http.StatusUnauthorized)
		return
	}

	hardware, err := json.Marshal(request.Hardware)
	if err != nil {
		ctx.Logger.Error("Failed to encode benchmark hardware", "user_id", user.Id, "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	benchmark := &schemas.Benchmark{
		UserId:     user.Id,
		Smoothness: request.Smoothness,
		Framerate:  request.Framerate,
		Score:      request.RawScore,
		Client:     request.Client,
		Grade:      request.Grade(),
		Hardware:   hardware,
	}
	if err := ctx.State.Benchmarks.Create(benchmark); err != nil {
		ctx.Logger.Error("Failed to create benchmark", "user_id", user.Id, "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.LatestActivity = time.Now()
	ctx.State.Users.Update(user, "latest_activity")

	ctx.Logger.Info("Submitted benchmark", "user_id", user.Id, "id", benchmark.Id)
	ctx.RenderText(http.StatusOK, strconv.Itoa(benchmark.Id))
}

func (h *BenchmarkHardware) ValidateRenderer() error {
	h.Renderer = strings.TrimSpace(h.Renderer)
	if h.Renderer == "" || len(h.Renderer) > 12 {
		return fmt.Errorf("renderer must contain between 1 and 12 characters")
	}
	return nil
}

func (h *BenchmarkHardware) ValidateResolution() error {
	if h.Resolution == nil {
		return nil
	}

	resolution := strings.ToLower(strings.TrimSpace(*h.Resolution))
	if resolution == "" || len(resolution) > 32 {
		return fmt.Errorf("resolution must contain between 1 and 32 characters")
	}

	width, height, ok := strings.Cut(resolution, "x")
	if !ok {
		return fmt.Errorf("resolution must use <width>x<height> format")
	}
	if _, err := strconv.ParseUint(width, 10, 64); err != nil {
		return fmt.Errorf("parse resolution width: %w", err)
	}
	if _, err := strconv.ParseUint(height, 10, 64); err != nil {
		return fmt.Errorf("parse resolution height: %w", err)
	}
	h.Resolution = new(resolution)
	return nil
}

func (h *BenchmarkHardware) ValidateDotNetVersion() error {
	if h.DotNetVersion == nil {
		return nil
	}
	dotNetVersion := strings.TrimSpace(*h.DotNetVersion)
	if dotNetVersion == "" || len(dotNetVersion) > 64 {
		return fmt.Errorf("dotnet_version must contain between 1 and 64 characters")
	}
	h.DotNetVersion = new(dotNetVersion)
	return nil
}

func (h *BenchmarkHardware) ValidateArchitecture() error {
	if h.ClientArchitecture == nil {
		return nil
	}
	architecture := strings.ToLower(strings.TrimSpace(*h.ClientArchitecture))
	if architecture != "32 bit" && architecture != "64 bit" {
		return fmt.Errorf("client_architecture must be 32 bit or 64 bit")
	}
	h.ClientArchitecture = new(architecture)
	return nil
}

func (h *BenchmarkHardware) ValidateFullHardware() error {
	if h.Cores != nil && *h.Cores <= 0 {
		return fmt.Errorf("cores must be a positive integer")
	}
	if h.Threads != nil && *h.Threads <= 0 {
		return fmt.Errorf("threads must be a positive integer")
	}
	if h.RAM != nil && *h.RAM <= 0 {
		return fmt.Errorf("ram must be a positive integer")
	}
	return nil
}
