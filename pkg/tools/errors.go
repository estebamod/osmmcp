// Package tools provides the OpenStreetMap MCP tools implementations.
package tools

import (
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

// APIError represents an error that occurred while communicating with
// an external API service, with information to help users recover.
type APIError struct {
	Service     string // The API service name (e.g., "Nominatim", "Overpass")
	StatusCode  int    // HTTP status code
	Message     string // Error message
	Recoverable bool   // Whether the error can be recovered from
	Guidance    string // Guidance for users on how to recover
}

// Error implements the error interface and provides a formatted error message.
func (e *APIError) Error() string {
	if e.Guidance != "" {
		return fmt.Sprintf("%s API error (%d): %s. %s", e.Service, e.StatusCode, e.Message, e.Guidance)
	}
	return fmt.Sprintf("%s API error (%d): %s", e.Service, e.StatusCode, e.Message)
}

// Common error guidance messages
const (
	// Nominatim guidance
	GuidanceNominatimAddressFormat = "Try using a more standard address format or provide city and country."
	GuidanceNominatimRateLimit     = "Please try again in a few seconds."
	GuidanceNominatimTimeout       = "Check your internet connection and try again, or use different geocoding parameters."
	GuidanceNominatimGeneral       = "Check your address formatting and try again."

	// Overpass guidance
	GuidanceOverpassTimeout   = "Consider simplifying your query by reducing the search radius or adding more specific filters."
	GuidanceOverpassRateLimit = "The Overpass API is currently experiencing high load. Please try again in a minute."
	GuidanceOverpassSyntax    = "There's an issue with the query format. Try simplifying your search."
	GuidanceOverpassMemory    = "The query requires too much memory. Try reducing the search area or adding more specific filters."
	GuidanceOverpassGeneral   = "Try a smaller search radius or fewer search criteria."

	// OSRM guidance
	GuidanceOSRMRouteNotFound = "No route could be found between the specified points. Try locations with accessible roads."
	GuidanceOSRMRateLimit     = "The routing service is experiencing high load. Please try again in a few seconds."
	GuidanceOSRMTimeout       = "The routing request timed out. Try a shorter route or check your internet connection."
	GuidanceOSRMGeneral       = "Check that your coordinates are accessible by the specified transport mode."

	// Generic guidance
	GuidanceGeneral      = "Please try again later or modify your request parameters."
	GuidanceNetworkError = "Check your internet connection and try again."
	GuidanceDataError    = "The data received was incomplete or malformed. Try different search parameters."
)

// NewAPIError creates a new APIError with appropriate guidance based on status code.
func NewAPIError(service string, statusCode int, message, guidance string) *APIError {
	// Use provided guidance if available, otherwise infer based on status code
	if guidance == "" {
		switch statusCode {
		case http.StatusTooManyRequests:
			guidance = "Rate limit exceeded. Please try again in a few moments."
		case http.StatusRequestTimeout, http.StatusGatewayTimeout:
			guidance = "The request timed out. Try reducing the search area or simplifying the query."
		case http.StatusBadRequest:
			guidance = "The request was invalid. Check your parameters and try again."
		case http.StatusInternalServerError:
			guidance = "The server encountered an error. This is likely temporary, please try again later."
		case http.StatusServiceUnavailable:
			guidance = "The service is temporarily unavailable. Please try again later."
		default:
			guidance = GuidanceGeneral
		}
	}

	return &APIError{
		Service:     service,
		StatusCode:  statusCode,
		Message:     message,
		Recoverable: statusCode != http.StatusBadRequest, // Most errors except bad requests are recoverable
		Guidance:    guidance,
	}
}

// ErrorWithGuidance returns a properly formatted error response with user guidance.
func ErrorWithGuidance(err *APIError) *mcp.CallToolResult {
	errorText := fmt.Sprintf("Error: %s\n\nGuidance: %s", err.Message, err.Guidance)
	return mcp.NewToolResultError(errorText)
}

// ValidationError creates an error for invalid coordinate or radius parameters.
func ValidationError(lat, lon, radius float64, maxRadius float64) *APIError {
	var message string

	if lat < -90 || lat > 90 {
		message = fmt.Sprintf("Invalid latitude value: %f (must be between -90 and 90)", lat)
	} else if lon < -180 || lon > 180 {
		message = fmt.Sprintf("Invalid longitude value: %f (must be between -180 and 180)", lon)
	} else if radius <= 0 {
		message = "Radius must be greater than 0"
	} else if radius > maxRadius {
		message = fmt.Sprintf("Radius too large: %f (maximum allowed is %f meters)", radius, maxRadius)
	} else {
		message = "Invalid parameters"
	}

	guidance := "Please correct the parameters and try again."

	return &APIError{
		Service:     "Validation",
		StatusCode:  http.StatusBadRequest,
		Message:     message,
		Recoverable: true,
		Guidance:    guidance,
	}
}
