// Package osm provides utilities for working with OpenStreetMap data.
package osm

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NERVsystems/osmmcp/pkg/geo"
)

const (
	// API endpoints
	NominatimBaseURL = "https://nominatim.openstreetmap.org"
	OverpassBaseURL  = "https://overpass-api.de/api/interpreter"
	OSRMBaseURL      = "https://router.project-osrm.org"

	// User agent for API requests (required by Nominatim's usage policy)
	UserAgent = "osm-mcp-server/0.1.0"

	// Earth radius in meters (approximate) - re-exported from geo package
	EarthRadius = geo.EarthRadius
)

// NewClient returns an HTTP client configured for OSM API requests
// Deprecated: Use GetClient(ctx) instead for connection pooling
func NewClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
			IdleConnTimeout:     30 * time.Second,
		},
	}
}

// HaversineDistance calculates the great-circle distance between two points
// Deprecated: Use geo.HaversineDistance instead
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	return geo.HaversineDistance(lat1, lon1, lat2, lon2)
}

// NewBoundingBox creates a new empty bounding box
// Deprecated: Use geo.NewBoundingBox instead
func NewBoundingBox() *geo.BoundingBox {
	return geo.NewBoundingBox()
}

// CategoryMap maps common category names to OSM tags
var CategoryMap = map[string]map[string][]string{
	"restaurant": {
		"amenity": {"restaurant", "fast_food", "cafe", "bar", "pub"},
	},
	"cafe": {
		"amenity": {"cafe"},
	},
	"bar": {
		"amenity": {"bar", "pub"},
	},
	"hotel": {
		"tourism": {"hotel", "motel", "hostel", "guest_house"},
	},
	"park": {
		"leisure": {"park", "garden", "nature_reserve"},
	},
	"shop": {
		"shop": {"supermarket", "convenience", "mall", "department_store"},
	},
	"supermarket": {
		"shop": {"supermarket"},
	},
	"hospital": {
		"amenity": {"hospital", "clinic"},
	},
	"pharmacy": {
		"amenity": {"pharmacy"},
	},
	"bank": {
		"amenity": {"bank", "atm"},
	},
	"school": {
		"amenity": {"school", "university", "college"},
	},
	"gas_station": {
		"amenity": {"fuel"},
	},
	"parking": {
		"amenity": {"parking"},
	},
	"museum": {
		"tourism": {"museum", "gallery"},
	},
	"cinema": {
		"amenity": {"cinema"},
	},
	"gym": {
		"leisure": {"fitness_centre", "sports_centre"},
	},
	"library": {
		"amenity": {"library"},
	},
	"bus_station": {
		"highway": {"bus_stop"},
		"amenity": {"bus_station"},
	},
	"train_station": {
		"railway": {"station", "halt", "tram_stop"},
	},
	"airport": {
		"aeroway": {"aerodrome", "terminal"},
	},
	// EV specific categories
	"charging_station": {
		"amenity": {"charging_station"},
	},
}

// ValidateCoords validates latitude and longitude values
// Returns an error if the coordinates are invalid
func ValidateCoords(lat, lon float64) error {
	if lat < -90 || lat > 90 {
		return fmt.Errorf("invalid latitude: %f (must be between -90 and 90)", lat)
	}
	if lon < -180 || lon > 180 {
		return fmt.Errorf("invalid longitude: %f (must be between -180 and 180)", lon)
	}
	return nil
}
