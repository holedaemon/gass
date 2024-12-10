// Package logger provides a shorthand for creating a new *[slog.Logger].
package logger

import (
	"log/slog"
	"os"
)

// New creates a new *[slog.Logger].
func New(debug bool) *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
	}

	if debug {
		opts.Level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}
