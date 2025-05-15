package testutil

import (
	"bytes"
	"testing"
)

func TestNewTestLogger(t *testing.T) {
	// Test with a buffer
	buf := &bytes.Buffer{}
	logger := NewTestLogger(buf)
	if logger == nil {
		t.Error("NewTestLogger returned nil")
	}

	// Test logging
	logger.Info("test message", "key", "value")
	if buf.Len() == 0 {
		t.Error("Logger did not write to buffer")
	}

	// Test with nil writer (should use io.Discard)
	logger = NewTestLogger(nil)
	if logger == nil {
		t.Error("NewTestLogger returned nil with nil writer")
	}
}

func TestDiscardLogger(t *testing.T) {
	logger := DiscardLogger()
	if logger == nil {
		t.Error("DiscardLogger returned nil")
	}

	// Test that it doesn't panic
	logger.Info("test message", "key", "value")
	logger.Debug("debug message", "key", "value")
	logger.Warn("warning message", "key", "value")
	logger.Error("error message", "key", "value")
}
