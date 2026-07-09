package logging

import (
	"io"
	"log/slog"
)

func SetDefault(component string, level slog.Level, additionalWriters ...io.Writer) {
	writers := append([]io.Writer{GetConsoleWriter()}, additionalWriters...)
	slog.SetDefault(NewComponentLogger(component, level, writers...))
}

func init() {
	// Set default logger to "titanic" with Info level
	// This may be overridden by the application
	SetDefault("titanic", slog.LevelInfo)
}
