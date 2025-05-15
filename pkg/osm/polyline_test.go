// Package osm provides utilities for working with OpenStreetMap data.
package osm

import (
	"testing"

	"github.com/NERVsystems/osmmcp/pkg/geo"
)

// TestDecodePolyline tests the decoding of polyline strings using the Polyline5 format.
// All test cases use 5 decimal places of precision (1e-5) for coordinates.
func TestDecodePolyline(t *testing.T) {
	testCases := []struct {
		name     string
		encoded  string
		expected []geo.Location
	}{
		{
			name:     "Empty string",
			encoded:  "",
			expected: []geo.Location{},
		},
		{
			name:    "Single point",
			encoded: "_p~iF~ps|U",
			expected: []geo.Location{
				{Latitude: 38.5, Longitude: -120.2},
			},
		},
		{
			name:    "Multiple points",
			encoded: "_p~iF~ps|U_ulLnnqC_mqNvxq`@",
			expected: []geo.Location{
				{Latitude: 38.5, Longitude: -120.2},
				{Latitude: 40.7, Longitude: -120.95},
				{Latitude: 43.252, Longitude: -126.453},
			},
		},
		{
			name:    "Negative coordinates",
			encoded: "f{xyCwuy~W",
			expected: []geo.Location{
				{Latitude: -25.363882, Longitude: 131.044922},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := DecodePolyline(tc.encoded)

			// Check length
			if len(result) != len(tc.expected) {
				t.Errorf("Expected %d points, got %d", len(tc.expected), len(result))
				return
			}

			// Check each point
			for i, expected := range tc.expected {
				if !almostEqual(result[i].Latitude, expected.Latitude, 0.00001) ||
					!almostEqual(result[i].Longitude, expected.Longitude, 0.00001) {
					t.Errorf("Point %d: expected %v, got %v", i, expected, result[i])
				}
			}
		})
	}
}

// TestEncodePolyline tests the encoding of location points to polyline strings using the Polyline5 format.
// All test cases use 5 decimal places of precision (1e-5) for coordinates.
func TestEncodePolyline(t *testing.T) {
	testCases := []struct {
		name     string
		points   []geo.Location
		expected string
	}{
		{
			name:     "Empty slice",
			points:   []geo.Location{},
			expected: "",
		},
		{
			name: "Single point",
			points: []geo.Location{
				{Latitude: 38.5, Longitude: -120.2},
			},
			expected: "_p~iF~ps|U",
		},
		{
			name: "Multiple points",
			points: []geo.Location{
				{Latitude: 38.5, Longitude: -120.2},
				{Latitude: 40.7, Longitude: -120.95},
				{Latitude: 43.252, Longitude: -126.453},
			},
			expected: "_p~iF~ps|U_ulLnnqC_mqNvxq`@",
		},
		{
			name: "Negative coordinates",
			points: []geo.Location{
				{Latitude: -25.363882, Longitude: 131.044922},
			},
			expected: "f{xyCwuy~W",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := EncodePolyline(tc.points)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// TestPolylineRoundTrip tests that encoding and decoding a set of points
// results in the same coordinates, within a small tolerance.
func TestPolylineRoundTrip(t *testing.T) {
	testCases := []struct {
		name   string
		points []geo.Location
	}{
		{
			name:   "Empty slice",
			points: []geo.Location{},
		},
		{
			name: "Single point",
			points: []geo.Location{
				{Latitude: 38.5, Longitude: -120.2},
			},
		},
		{
			name: "Multiple points",
			points: []geo.Location{
				{Latitude: 38.5, Longitude: -120.2},
				{Latitude: 40.7, Longitude: -120.95},
				{Latitude: 43.252, Longitude: -126.453},
			},
		},
		{
			name: "SF to Oakland",
			points: []geo.Location{
				{Latitude: 37.7749, Longitude: -122.4194}, // SF
				{Latitude: 37.8044, Longitude: -122.2711}, // Oakland
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encode
			encoded := EncodePolyline(tc.points)

			// Decode
			decoded := DecodePolyline(encoded)

			// Compare
			if len(decoded) != len(tc.points) {
				t.Errorf("Round trip length mismatch: original %d, result %d", len(tc.points), len(decoded))
				return
			}

			for i, original := range tc.points {
				if !almostEqual(decoded[i].Latitude, original.Latitude, 0.00001) ||
					!almostEqual(decoded[i].Longitude, original.Longitude, 0.00001) {
					t.Errorf("Point %d mismatch after round trip: original %v, result %v",
						i, original, decoded[i])
				}
			}
		})
	}
}

// almostEqual checks if two float64 values are equal within a tolerance.
// This is used for comparing floating-point coordinates.
func almostEqual(a, b, tolerance float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}
