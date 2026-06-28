package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"strings"
)

const (
	darkGrey = "\x1b[90m"
	grey     = "\x1b[38;20m"
	yellow   = "\x1b[33;20m"
	red      = "\x1b[31;20m"
	boldRed  = "\x1b[31;1m"
	cyan     = "\x1b[96m"
	reset    = "\x1b[0m"
)

// ColorLogger is a custom logger that is used by
// default on all titanic related projects
type ColorLogger struct {
	out   io.Writer
	level slog.Level
	attrs []slog.Attr
}

func NewColorLogger(out io.Writer, level slog.Level) *ColorLogger {
	return &ColorLogger{
		out:   out,
		level: level,
	}
}

func NewLogger(out io.Writer, level slog.Level) *slog.Logger {
	handler := NewColorLogger(out, level)
	return slog.New(handler)
}

func NewMultiLogger(level slog.Level, writers ...io.Writer) *slog.Logger {
	multi := io.MultiWriter(writers...)
	handler := NewColorLogger(multi, level)
	return slog.New(handler)
}

func NewComponentLogger(component string, level slog.Level, writers ...io.Writer) *slog.Logger {
	multi := io.MultiWriter(writers...)
	handler := NewColorLogger(multi, level)
	return slog.New(handler).With("component", component)
}

func (h *ColorLogger) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *ColorLogger) Handle(_ context.Context, record slog.Record) error {
	var attrs []string
	var attrString string
	var component string

	processAttr := func(a slog.Attr) {
		if a.Key == "component" {
			component = a.Value.String()
		} else {
			attrs = append(attrs, formatAttr(a))
		}
	}
	for _, a := range h.attrs {
		processAttr(a)
	}

	// Extract attributes from the record and format them
	record.Attrs(func(a slog.Attr) bool {
		processAttr(a)
		return true
	})

	// If no "component" attribute is found, we use `titanic` by default
	if component == "" {
		component = "titanic"
	}

	timestamp := record.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(record.Level.String())
	color := colorForLevel(record.Level)

	if len(attrs) > 0 {
		attrString += " ->"
		attrString += " " + strings.Join(attrs, " ")
	}

	fmt.Fprintf(
		h.out,
		"%s[%s] - <%s> %s: %s%s%s\n",
		grey,
		timestamp,
		component,
		color+level,
		record.Message,
		reset+darkGrey+attrString,
		reset,
	)
	return nil
}

func (h *ColorLogger) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := slices.Concat(h.attrs, attrs)
	return &ColorLogger{
		out:   h.out,
		level: h.level,
		attrs: newAttrs,
	}
}

func (h *ColorLogger) WithGroup(name string) slog.Handler {
	return h
}

func colorForLevel(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return red
	case level >= slog.LevelWarn:
		return yellow
	case level >= slog.LevelInfo:
		return cyan
	default:
		return grey
	}
}

func formatAttr(a slog.Attr) string {
	switch a.Value.Kind() {
	case slog.KindString:
		return fmt.Sprintf("%s=%q", a.Key, a.Value.String())
	case slog.KindGroup:
		var parts []string
		for _, ga := range a.Value.Group() {
			parts = append(parts, formatAttr(ga))
		}
		return fmt.Sprintf("%s={%s}", a.Key, strings.Join(parts, " "))
	default:
		return fmt.Sprintf("%s=%v", a.Key, a.Value.Any())
	}
}
