// Package osm provides utilities for working with OpenStreetMap data.
package osm

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

const (
	// API endpoints
	NominatimBaseURL = "https://nominatim.openstreetmap.org"
	OverpassBaseURL  = "https://overpass-api.de/api/interpreter"
	OSRMBaseURL      = "https://router.project-osrm.org"

	// User agent for API requests (required by Nominatim's usage policy)
	UserAgent = "osm-mcp-server/0.1.0"

	// Earth radius in meters (approximate)
	EarthRadius = 6371000.0
)

// NewClient returns an HTTP client configured for OSM API requests
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

// HaversineDistance calculates the great-circle distance between two points on a sphere
// using the Haversine formula. The result is in meters.
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Haversine formula
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := EarthRadius * c

	return distance
}

// BoundingBox represents a geographic bounding box with southwest and northeast corners
type BoundingBox struct {
	MinLat float64
	MinLon float64
	MaxLat float64
	MaxLon float64
}

// NewBoundingBox creates a new empty bounding box
func NewBoundingBox() *BoundingBox {
	return &BoundingBox{
		MinLat: 90.0,
		MinLon: 180.0,
		MaxLat: -90.0,
		MaxLon: -180.0,
	}
}

// ExtendWithPoint extends the bounding box to include the specified point
func (bb *BoundingBox) ExtendWithPoint(lat, lon float64) {
	if lat < bb.MinLat {
		bb.MinLat = lat
	}
	if lat > bb.MaxLat {
		bb.MaxLat = lat
	}
	if lon < bb.MinLon {
		bb.MinLon = lon
	}
	if lon > bb.MaxLon {
		bb.MaxLon = lon
	}
}

// Buffer adds a buffer around the bounding box in degrees
func (bb *BoundingBox) Buffer(bufferMeters float64) {
	// Convert meters to approximate degrees (crude approximation)
	// 0.01 degrees ≈ 1.11 km at the equator
	bufferDegrees := bufferMeters / 111000
	bb.MinLat -= bufferDegrees
	bb.MaxLat += bufferDegrees
	bb.MinLon -= bufferDegrees
	bb.MaxLon += bufferDegrees

	// Ensure coordinates are within valid ranges
	if bb.MinLat < -90 {
		bb.MinLat = -90
	}
	if bb.MaxLat > 90 {
		bb.MaxLat = 90
	}
	if bb.MinLon < -180 {
		bb.MinLon = -180
	}
	if bb.MaxLon > 180 {
		bb.MaxLon = 180
	}
}

// String returns a string representation of the bounding box
func (bb *BoundingBox) String() string {
	return fmt.Sprintf("(%f,%f,%f,%f)", bb.MinLat, bb.MinLon, bb.MaxLat, bb.MaxLon)
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
