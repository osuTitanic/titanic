package server

import (
	"log/slog"
)

type BanchoTcpServer[C any] struct {
	logger  slog.Logger
	timeout int
}

// TODO: ...
