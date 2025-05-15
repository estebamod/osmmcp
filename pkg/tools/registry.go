// Package tools provides the OpenStreetMap MCP tools implementations.
package tools

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Registry holds all MCP tool registrations for the OpenStreetMap service.
type Registry struct {
	logger *slog.Logger
}

// NewRegistry creates a new MCP tool registry.
func NewRegistry(logger *slog.Logger) *Registry {
	return &Registry{
		logger: logger,
	}
}

// ToolDefinition represents an OpenStreetMap MCP tool definition.
type ToolDefinition struct {
	Name        string
	Description string
	Tool        mcp.Tool
	Handler     func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

// GetToolDefinitions returns all OpenStreetMap MCP tool definitions.
func (r *Registry) GetToolDefinitions() []ToolDefinition {
	return []ToolDefinition{
		// Geocoding Tools
		{
			Name:        "geocode_address",
			Description: "Convert an address or place name to geographic coordinates",
			Tool:        GeocodeAddressTool(),
			Handler:     HandleGeocodeAddress,
		},
		{
			Name:        "reverse_geocode",
			Description: "Convert geographic coordinates to a human-readable address",
			Tool:        ReverseGeocodeTool(),
			Handler:     HandleReverseGeocode,
		},

		// Place Search Tools
		{
			Name:        "find_nearby_places",
			Description: "Find points of interest near a specific location",
			Tool:        FindNearbyPlacesTool(),
			Handler:     HandleFindNearbyPlaces,
		},
		{
			Name:        "search_category",
			Description: "Search for places by category within a rectangular area",
			Tool:        SearchCategoryTool(),
			Handler:     HandleSearchCategory,
		},

		// Routing Tools
		{
			Name:        "get_route_directions",
			Description: "Get directions for a route between two locations",
			Tool:        GetRouteDirectionsTool(),
			Handler:     HandleGetRouteDirections,
		},
		{
			Name:        "suggest_meeting_point",
			Description: "Suggest an optimal meeting point for multiple people",
			Tool:        SuggestMeetingPointTool(),
			Handler:     HandleSuggestMeetingPoint,
		},

		// Exploration Tools
		{
			Name:        "explore_area",
			Description: "Explore an area and get comprehensive information about it",
			Tool:        ExploreAreaTool(),
			Handler:     HandleExploreArea,
		},

		// EV Tools
		{
			Name:        "find_charging_stations",
			Description: "Find electric vehicle charging stations near a location",
			Tool:        FindChargingStationsTool(),
			Handler:     HandleFindChargingStations,
		},
		{
			Name:        "find_route_charging_stations",
			Description: "Find electric vehicle charging stations along a route",
			Tool:        FindRouteChargingStationsTool(),
			Handler:     HandleFindRouteChargingStations,
		},

		// Education Tools
		{
			Name:        "find_schools_nearby",
			Description: "Find educational institutions near a specific location",
			Tool:        FindSchoolsNearbyTool(),
			Handler:     HandleFindSchoolsNearby,
		},

		// Commute Tools
		{
			Name:        "analyze_commute",
			Description: "Analyze transportation options between home and work locations",
			Tool:        AnalyzeCommuteTool(),
			Handler:     HandleAnalyzeCommute,
		},

		// Neighborhood Analysis Tools
		{
			Name:        "analyze_neighborhood",
			Description: "Evaluate neighborhood livability for real estate and relocation decisions",
			Tool:        AnalyzeNeighborhoodTool(),
			Handler:     HandleAnalyzeNeighborhood,
		},

		// Parking Tools
		{
			Name:        "find_parking_facilities",
			Description: "Find parking facilities near a specific location",
			Tool:        FindParkingAreasTool(),
			Handler:     HandleFindParkingFacilities,
		},
	}
}

// RegisterTools registers all tools with the MCP server.
func (r *Registry) RegisterTools(mcpServer *server.MCPServer) {
	for _, def := range r.GetToolDefinitions() {
		r.logger.Info("registering tool", "name", def.Name)
		mcpServer.AddTool(def.Tool, def.Handler)
	}
}
