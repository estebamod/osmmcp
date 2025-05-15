// Package tools provides the OpenStreetMap MCP tools implementations.
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/NERVsystems/osmmcp/pkg/osm"
	"github.com/mark3labs/mcp-go/mcp"
)

// ParkingArea represents a parking facility
type ParkingArea struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Location     Location `json:"location"`
	Distance     float64  `json:"distance,omitempty"`     // in meters
	Type         string   `json:"type,omitempty"`         // e.g., surface, underground, multi-storey
	Access       string   `json:"access,omitempty"`       // e.g., public, private, customers
	Capacity     int      `json:"capacity,omitempty"`     // number of parking spaces if available
	Fee          bool     `json:"fee,omitempty"`          // whether there's a parking fee
	MaxStay      string   `json:"max_stay,omitempty"`     // maximum parking duration if available
	Availability string   `json:"availability,omitempty"` // if real-time availability is known
	Wheelchair   bool     `json:"wheelchair,omitempty"`   // wheelchair accessibility
	Operator     string   `json:"operator,omitempty"`     // who operates the facility
}

// FindParkingAreasTool returns a tool definition for finding parking facilities
func FindParkingAreasTool() mcp.Tool {
	return mcp.NewTool("find_parking_facilities",
		mcp.WithDescription("Find parking facilities near a specific location"),
		mcp.WithNumber("latitude",
			mcp.Required(),
			mcp.Description("The latitude coordinate of the center point"),
		),
		mcp.WithNumber("longitude",
			mcp.Required(),
			mcp.Description("The longitude coordinate of the center point"),
		),
		mcp.WithNumber("radius",
			mcp.Description("Search radius in meters (max 5000)"),
			mcp.DefaultNumber(1000),
		),
		mcp.WithString("type",
			mcp.Description("Optional type filter (e.g., surface, underground, multi-storey)"),
			mcp.DefaultString(""),
		),
		mcp.WithBoolean("include_private",
			mcp.Description("Whether to include private parking facilities"),
			mcp.DefaultBool(false),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of results to return"),
			mcp.DefaultNumber(10),
		),
	)
}

// HandleFindParkingFacilities implements finding parking facilities functionality
func HandleFindParkingFacilities(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger := slog.Default().With("tool", "find_parking_facilities")

	// Parse input parameters
	latitude := mcp.ParseFloat64(req, "latitude", 0)
	longitude := mcp.ParseFloat64(req, "longitude", 0)
	radius := mcp.ParseFloat64(req, "radius", 1000)
	facilityType := mcp.ParseString(req, "type", "")
	includePrivate := mcp.ParseBoolean(req, "include_private", false)
	limit := int(mcp.ParseFloat64(req, "limit", 10))

	// Basic validation
	if latitude < -90 || latitude > 90 {
		return ErrorResponse("Latitude must be between -90 and 90"), nil
	}
	if longitude < -180 || longitude > 180 {
		return ErrorResponse("Longitude must be between -180 and 180"), nil
	}
	if radius <= 0 || radius > 5000 {
		return ErrorResponse("Radius must be between 1 and 5000 meters"), nil
	}
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 50 {
		limit = 50 // Max limit
	}

	// Build Overpass query for parking facilities
	var queryBuilder strings.Builder
	queryBuilder.WriteString("[out:json];")

	// Search for nodes with amenity=parking
	queryBuilder.WriteString(fmt.Sprintf("(node(around:%f,%f,%f)[amenity=parking];", radius, latitude, longitude))

	// Search for ways (areas) with amenity=parking
	queryBuilder.WriteString(fmt.Sprintf("way(around:%f,%f,%f)[amenity=parking];", radius, latitude, longitude))

	// Search for relations with amenity=parking (for complex parking structures)
	queryBuilder.WriteString(fmt.Sprintf("relation(around:%f,%f,%f)[amenity=parking];", radius, latitude, longitude))

	// Complete the query
	queryBuilder.WriteString(");out center;")

	// Build request
	reqURL, err := url.Parse(osm.OverpassBaseURL)
	if err != nil {
		logger.Error("failed to parse URL", "error", err)
		return ErrorResponse("Internal server error"), nil
	}

	// Make HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL.String(), strings.NewReader("data="+url.QueryEscape(queryBuilder.String())))
	if err != nil {
		logger.Error("failed to create request", "error", err)
		return ErrorResponse("Failed to create request"), nil
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("User-Agent", osm.UserAgent)

	// Execute request
	client := osm.NewClient()
	resp, err := client.Do(httpReq)
	if err != nil {
		logger.Error("failed to execute request", "error", err)
		return ErrorResponse("Failed to communicate with OSM service"), nil
	}
	defer resp.Body.Close()

	// Process response
	if resp.StatusCode != http.StatusOK {
		logger.Error("OSM service returned error", "status", resp.StatusCode)
		return ErrorResponse(fmt.Sprintf("OSM service error: %d", resp.StatusCode)), nil
	}

	// Parse response
	var overpassResp struct {
		Elements []struct {
			ID     int     `json:"id"`
			Type   string  `json:"type"`
			Lat    float64 `json:"lat,omitempty"`
			Lon    float64 `json:"lon,omitempty"`
			Center *struct {
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			} `json:"center,omitempty"`
			Tags map[string]string `json:"tags"`
		} `json:"elements"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&overpassResp); err != nil {
		logger.Error("failed to decode response", "error", err)
		return ErrorResponse("Failed to parse parking facilities data"), nil
	}

	// Convert to ParkingArea objects and calculate distances
	facilities := make([]ParkingArea, 0)
	for _, element := range overpassResp.Elements {
		// Get coordinates (handling both nodes and ways/relations)
		var lat, lon float64
		if element.Type == "node" {
			lat = element.Lat
			lon = element.Lon
		} else if (element.Type == "way" || element.Type == "relation") && element.Center != nil {
			lat = element.Center.Lat
			lon = element.Center.Lon
		} else {
			continue // Skip elements without coordinates
		}

		// Skip private facilities if not requested
		if !includePrivate {
			access := strings.ToLower(element.Tags["access"])
			if access == "private" || access == "customers" || access == "permit" {
				continue
			}
		}

		// Apply facility type filter if specified
		if facilityType != "" {
			parkingType := strings.ToLower(element.Tags["parking"])
			if parkingType != "" && !strings.Contains(parkingType, strings.ToLower(facilityType)) {
				continue
			}
		}

		// Calculate distance
		distance := osm.HaversineDistance(
			latitude, longitude,
			lat, lon,
		)

		// Parse capacity if available
		capacity := 0
		if capacityStr := element.Tags["capacity"]; capacityStr != "" {
			_, _ = fmt.Sscanf(capacityStr, "%d", &capacity)
		} else if capacityStr := element.Tags["capacity:disabled"]; capacityStr != "" {
			_, _ = fmt.Sscanf(capacityStr, "%d", &capacity)
		}

		// Determine if there's a fee
		hasFee := false
		if feeStr := element.Tags["fee"]; feeStr == "yes" || feeStr == "true" {
			hasFee = true
		}

		// Determine wheelchair accessibility
		hasWheelchair := false
		if wheelchairStr := element.Tags["wheelchair"]; wheelchairStr == "yes" || wheelchairStr == "designated" {
			hasWheelchair = true
		}

		// Create facility object
		name := element.Tags["name"]
		if name == "" {
			// Generate a generic name if none exists
			parkingType := element.Tags["parking"]
			if parkingType == "" {
				parkingType = "parking"
			}
			name = fmt.Sprintf("%s parking", strings.Title(parkingType))
		}

		facility := ParkingArea{
			ID:   fmt.Sprintf("%d", element.ID),
			Name: name,
			Location: Location{
				Latitude:  lat,
				Longitude: lon,
			},
			Distance:   distance,
			Type:       element.Tags["parking"],
			Access:     element.Tags["access"],
			Capacity:   capacity,
			Fee:        hasFee,
			MaxStay:    element.Tags["maxstay"],
			Wheelchair: hasWheelchair,
			Operator:   element.Tags["operator"],
		}

		facilities = append(facilities, facility)
	}

	// Sort facilities by distance (closest first)
	for i := 0; i < len(facilities); i++ {
		for j := i + 1; j < len(facilities); j++ {
			if facilities[i].Distance > facilities[j].Distance {
				facilities[i], facilities[j] = facilities[j], facilities[i]
			}
		}
	}

	// Limit results
	if len(facilities) > limit {
		facilities = facilities[:limit]
	}

	// Create output
	output := struct {
		Facilities []ParkingArea `json:"facilities"`
	}{
		Facilities: facilities,
	}

	// Return result
	resultBytes, err := json.Marshal(output)
	if err != nil {
		logger.Error("failed to marshal result", "error", err)
		return ErrorResponse("Failed to generate result"), nil
	}

	return mcp.NewToolResultText(string(resultBytes)), nil
}
