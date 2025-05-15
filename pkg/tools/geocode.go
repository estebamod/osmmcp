package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	// Nominatim is OSM's geocoding service
	nominatimBaseURL = "https://nominatim.openstreetmap.org"
)

// GeocodeAddressInput defines the input parameters for geocoding an address
type GeocodeAddressInput struct {
	Address string `json:"address"`
}

// GeocodeAddressOutput defines the output format for geocoded addresses
type GeocodeAddressOutput struct {
	Place Place `json:"place"`
}

// GeocodeAddressTool returns a tool definition for geocoding addresses
func GeocodeAddressTool() mcp.Tool {
	return mcp.NewTool("geocode_address",
		mcp.WithDescription("Convert an address or place name to geographic coordinates"),
		mcp.WithString("address",
			mcp.Required(),
			mcp.Description("The address or place name to geocode"),
		),
	)
}

// HandleGeocodeAddress implements the geocoding functionality
func HandleGeocodeAddress(ctx context.Context, rawInput mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger := slog.Default().With("tool", "geocode_address")

	// Parse input
	var input GeocodeAddressInput
	address := mcp.ParseString(rawInput, "address", "")
	input.Address = address

	if input.Address == "" {
		return ErrorResponse("Address must not be empty"), nil
	}

	// Build request URL
	reqURL, err := url.Parse(fmt.Sprintf("%s/search", nominatimBaseURL))
	if err != nil {
		logger.Error("failed to parse URL", "error", err)
		return ErrorResponse("Internal server error"), nil
	}

	// Add query parameters
	q := reqURL.Query()
	q.Add("q", input.Address)
	q.Add("format", "json")
	q.Add("limit", "1")
	reqURL.RawQuery = q.Encode()

	// Make HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		logger.Error("failed to create request", "error", err)
		return ErrorResponse("Failed to create request"), nil
	}

	// Set user agent (required by Nominatim's usage policy)
	req.Header.Set("User-Agent", "osm-mcp-server/0.1.0")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("failed to execute request", "error", err)
		return ErrorResponse("Failed to communicate with geocoding service"), nil
	}
	defer resp.Body.Close()

	// Process response
	if resp.StatusCode != http.StatusOK {
		logger.Error("geocoding service returned error", "status", resp.StatusCode)
		return ErrorResponse(fmt.Sprintf("Geocoding service error: %d", resp.StatusCode)), nil
	}

	// Parse response
	var results []struct {
		PlaceID     string  `json:"place_id"`
		DisplayName string  `json:"display_name"`
		Lat         string  `json:"lat"`
		Lon         string  `json:"lon"`
		Type        string  `json:"type"`
		Importance  float64 `json:"importance"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		logger.Error("failed to decode response", "error", err)
		return ErrorResponse("Failed to parse geocoding response"), nil
	}

	// Handle no results
	if len(results) == 0 {
		return ErrorResponse("No results found for the address"), nil
	}

	// Get first result
	result := results[0]

	// Convert lat/lon to float64
	var lat, lon float64
	fmt.Sscanf(result.Lat, "%f", &lat)
	fmt.Sscanf(result.Lon, "%f", &lon)

	// Create output
	output := GeocodeAddressOutput{
		Place: Place{
			ID:   result.PlaceID,
			Name: result.DisplayName,
			Location: Location{
				Latitude:  lat,
				Longitude: lon,
			},
			Address: Address{
				Formatted: result.DisplayName,
			},
		},
	}

	// Return result
	resultBytes, err := json.Marshal(output)
	if err != nil {
		logger.Error("failed to marshal result", "error", err)
		return ErrorResponse("Failed to generate result"), nil
	}

	return mcp.NewToolResultText(string(resultBytes)), nil
}

// ReverseGeocodeInput defines the input parameters for reverse geocoding
type ReverseGeocodeInput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// ReverseGeocodeOutput defines the output format for reverse geocoded coordinates
type ReverseGeocodeOutput struct {
	Place Place `json:"place"`
}

// ReverseGeocodeTool returns a tool definition for reverse geocoding
func ReverseGeocodeTool() mcp.Tool {
	return mcp.NewTool("reverse_geocode",
		mcp.WithDescription("Convert geographic coordinates to a human-readable address"),
		mcp.WithNumber("latitude",
			mcp.Required(),
			mcp.Description("The latitude coordinate"),
		),
		mcp.WithNumber("longitude",
			mcp.Required(),
			mcp.Description("The longitude coordinate"),
		),
	)
}

// HandleReverseGeocode implements the reverse geocoding functionality
func HandleReverseGeocode(ctx context.Context, rawInput mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger := slog.Default().With("tool", "reverse_geocode")

	// Parse input
	latitude := mcp.ParseFloat64(rawInput, "latitude", 0)
	longitude := mcp.ParseFloat64(rawInput, "longitude", 0)

	// Basic validation
	if latitude < -90 || latitude > 90 {
		return ErrorResponse("Latitude must be between -90 and 90"), nil
	}
	if longitude < -180 || longitude > 180 {
		return ErrorResponse("Longitude must be between -180 and 180"), nil
	}

	// Build request URL
	reqURL, err := url.Parse(fmt.Sprintf("%s/reverse", nominatimBaseURL))
	if err != nil {
		logger.Error("failed to parse URL", "error", err)
		return ErrorResponse("Internal server error"), nil
	}

	// Add query parameters
	q := reqURL.Query()
	q.Add("lat", fmt.Sprintf("%f", latitude))
	q.Add("lon", fmt.Sprintf("%f", longitude))
	q.Add("format", "json")
	q.Add("addressdetails", "1")
	reqURL.RawQuery = q.Encode()

	// Make HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		logger.Error("failed to create request", "error", err)
		return ErrorResponse("Failed to create request"), nil
	}

	// Set user agent (required by Nominatim's usage policy)
	req.Header.Set("User-Agent", "osm-mcp-server/0.1.0")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("failed to execute request", "error", err)
		return ErrorResponse("Failed to communicate with geocoding service"), nil
	}
	defer resp.Body.Close()

	// Process response
	if resp.StatusCode != http.StatusOK {
		logger.Error("geocoding service returned error", "status", resp.StatusCode)
		return ErrorResponse(fmt.Sprintf("Geocoding service error: %d", resp.StatusCode)), nil
	}

	// Parse response
	var result struct {
		PlaceID     string `json:"place_id"`
		DisplayName string `json:"display_name"`
		Lat         string `json:"lat"`
		Lon         string `json:"lon"`
		Address     struct {
			Road        string `json:"road"`
			HouseNumber string `json:"house_number"`
			City        string `json:"city"`
			Town        string `json:"town"`
			State       string `json:"state"`
			Country     string `json:"country"`
			PostCode    string `json:"postcode"`
		} `json:"address"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Error("failed to decode response", "error", err)
		return ErrorResponse("Failed to parse geocoding response"), nil
	}

	// Get city (could be in city or town field)
	city := result.Address.City
	if city == "" {
		city = result.Address.Town
	}

	// Convert lat/lon to float64
	var lat, lon float64
	fmt.Sscanf(result.Lat, "%f", &lat)
	fmt.Sscanf(result.Lon, "%f", &lon)

	// Create output
	output := ReverseGeocodeOutput{
		Place: Place{
			ID:   result.PlaceID,
			Name: result.DisplayName,
			Location: Location{
				Latitude:  lat,
				Longitude: lon,
			},
			Address: Address{
				Street:      result.Address.Road,
				HouseNumber: result.Address.HouseNumber,
				City:        city,
				State:       result.Address.State,
				Country:     result.Address.Country,
				PostalCode:  result.Address.PostCode,
				Formatted:   result.DisplayName,
			},
		},
	}

	// Return result
	resultBytes, err := json.Marshal(output)
	if err != nil {
		logger.Error("failed to marshal result", "error", err)
		return ErrorResponse("Failed to generate result"), nil
	}

	return mcp.NewToolResultText(string(resultBytes)), nil
}
