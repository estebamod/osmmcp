// Package testutil provides utilities for testing.
package testutil

import (
	"io"
	"log/slog"
)

// NewTestLogger creates a new logger for testing
// If writer is nil, it will use io.Discard
func NewTestLogger(w io.Writer) *slog.Logger {
	if w == nil {
		w = io.Discard
	}
	return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

// DiscardLogger returns a logger that discards all output
func DiscardLogger() *slog.Logger {
	return NewTestLogger(nil)
}
