// Package server provides the MCP server implementation for the OpenStreetMap integration.
package server

import (
	"log/slog"

	"github.com/NERVsystems/osmmcp/pkg/tools"
	"github.com/mark3labs/mcp-go/server"
)

const (
	// ServerName is the name of the MCP server
	ServerName = "osm-mcp-server"

	// ServerVersion is the version of the MCP server
	ServerVersion = "0.1.0"
)

// Server encapsulates the MCP server with OpenStreetMap tools.
type Server struct {
	srv *server.MCPServer
}

// NewServer creates a new OpenStreetMap MCP server with all tools registered.
func NewServer() (*Server, error) {
	logger := slog.Default()
	logger.Info("initializing OpenStreetMap MCP server",
		"name", ServerName,
		"version", ServerVersion)

	// Create MCP server with options
	srv := server.NewMCPServer(
		ServerName,
		ServerVersion,
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	// Create tool registry and register all tools
	registry := tools.NewRegistry(logger)
	registry.RegisterTools(srv)

	return &Server{srv: srv}, nil
}

// Run starts the MCP server using stdin/stdout for communication.
func (s *Server) Run() error {
	return server.ServeStdio(s.srv)
}
