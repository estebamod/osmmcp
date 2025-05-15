package server

import (
	"context"
	"testing"
)

func TestNewServer(t *testing.T) {
	// TODO: Implement server creation tests
	s, err := NewServer()
	if err != nil {
		t.Errorf("NewServer() error = %v", err)
	}
	if s == nil {
		t.Error("NewServer() returned nil server")
	}
}

func TestServer_Run(t *testing.T) {
	// TODO: Implement server run tests
	s, err := NewServer()
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the server in a goroutine
	go func() {
		if err := s.RunWithContext(ctx); err != nil {
			t.Errorf("RunWithContext() error = %v", err)
		}
	}()

	// Shutdown the server
	s.Shutdown()
	s.WaitForShutdown()
}
