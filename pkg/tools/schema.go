// Package tools provides the OpenStreetMap MCP tools implementations.
package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
)

// Schema helper functions have been removed as they're not needed
// with the new mcp-go API v0.27.1.
// The new API uses a more fluent builder pattern directly.

// ErrorResponse is used for consistent error reporting
func ErrorResponse(message string) *mcp.CallToolResult {
	return mcp.NewToolResultError(message)
}
