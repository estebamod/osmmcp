package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/NERVsystems/osmmcp/pkg/server"
)

// Version information
var (
	version        bool
	debug          bool
	generateConfig string

	// Build information
	buildVersion = "0.1.0"
	buildCommit  = "unknown"
	buildDate    = "unknown"
)

func init() {
	flag.BoolVar(&version, "version", false, "Display version information")
	flag.BoolVar(&debug, "debug", false, "Enable debug logging")
	flag.StringVar(&generateConfig, "generate-config", "", "Generate a Claude Desktop Client config file at the specified path")
}

func main() {
	flag.Parse()

	// Configure logging
	var logLevel slog.Level
	if debug {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	// Show version and exit if requested
	if version {
		showVersion()
		return
	}

	// Generate Claude Desktop config if requested
	if generateConfig != "" {
		if err := generateClientConfig(generateConfig); err != nil {
			logger.Error("failed to generate config", "error", err)
			os.Exit(1)
		}
		logger.Info("successfully generated Claude Desktop Client config", "path", generateConfig)
		return
	}

	logger.Info("starting OpenStreetMap MCP server",
		"version", version,
		"log_level", logLevel.String())

	// Create and run the MCP server
	srv, err := server.NewServer()
	if err != nil {
		logger.Error("failed to create server", "error", err)
		os.Exit(1)
	}

	logger.Info("server initialized, waiting for requests")
	if err := srv.Run(); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}

// generateClientConfig creates or updates a Claude Desktop Client config file
func generateClientConfig(outputPath string) error {
	logger := slog.Default()

	// Get absolute path to executable
	execPath, err := os.Executable()
	if err != nil {
		execPath = os.Args[0] // Fallback to args if cannot get executable path
	}
	absExecPath, err := filepath.Abs(execPath)
	if err != nil {
		absExecPath = execPath // Use as is if cannot resolve absolute path
	}

	// Prepare our server config
	osmConfig := map[string]interface{}{
		"command": absExecPath,
		"args":    []string{},
	}

	// Define the config structure
	var config map[string]interface{}

	// Check if file exists already
	if _, err := os.Stat(outputPath); err == nil {
		// File exists, read it
		data, err := os.ReadFile(outputPath)
		if err != nil {
			return fmt.Errorf("failed to read existing config: %w", err)
		}

		// Parse existing JSON
		if err := json.Unmarshal(data, &config); err != nil {
			logger.Warn("existing config is not valid JSON, will create new", "error", err)
			config = make(map[string]interface{})
		}
	} else {
		// File doesn't exist, create new config
		config = make(map[string]interface{})
	}

	// Check if mcpServers exists, create it if not
	mcpServers, ok := config["mcpServers"].(map[string]interface{})
	if !ok {
		mcpServers = make(map[string]interface{})
		config["mcpServers"] = mcpServers
	}

	// Add or update our server
	mcpServers["OSM"] = osmConfig

	// Marshal to JSON with pretty printing
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Validate the JSON by attempting to unmarshal it
	var validation interface{}
	if err := json.Unmarshal(data, &validation); err != nil {
		return fmt.Errorf("generated invalid JSON: %w", err)
	}

	// Add a newline at the end for better formatting
	data = append(data, '\n')

	// Make sure parent directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write to the output file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// showVersion displays version information and exits
func showVersion() {
	fmt.Printf("osmmcp version %s (%s) built on %s\n", buildVersion, buildCommit, buildDate)
}
