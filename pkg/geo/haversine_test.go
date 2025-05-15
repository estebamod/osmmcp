package geo

import (
	"math"
	"testing"
)

func TestHaversineDistance(t *testing.T) {
	// Test cases with known distances
	tests := []struct {
		lat1      float64
		lon1      float64
		lat2      float64
		lon2      float64
		expected  float64
		name      string
		tolerance float64 // Now represents relative tolerance (e.g., 0.001 for 0.1%)
	}{
		{
			name:      "Same point",
			lat1:      37.7749,
			lon1:      -122.4194,
			lat2:      37.7749,
			lon2:      -122.4194,
			expected:  0,
			tolerance: 0.0001, // 0.01% for zero case
		},
		{
			name:      "Short distance - SF downtown to Market St",
			lat1:      37.7749,
			lon1:      -122.4194,
			lat2:      37.7734,
			lon2:      -122.4167,
			expected:  290.06, // Updated from GeographicLib
			tolerance: 0.001,  // 0.1% relative tolerance
		},
		{
			name:      "Medium distance - SF to Oakland",
			lat1:      37.7749,
			lon1:      -122.4194,
			lat2:      37.8044,
			lon2:      -122.2712,
			expected:  13429.63, // Updated from GeographicLib
			tolerance: 0.001,    // 0.1% relative tolerance
		},
		{
			name:      "Long distance - SF to NYC",
			lat1:      37.7749,
			lon1:      -122.4194,
			lat2:      40.7128,
			lon2:      -74.0060,
			expected:  4129936.81, // ~4130 km
			tolerance: 0.001,      // 0.1% relative tolerance
		},
		{
			name:      "Antipodal points",
			lat1:      37.7749,
			lon1:      -122.4194,
			lat2:      -37.7749,
			lon2:      57.5806,    // Opposite side of Earth
			expected:  20015086.8, // ~20,015 km (approx Earth diameter * π/2)
			tolerance: 0.001,      // 0.1% relative tolerance
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := HaversineDistance(tc.lat1, tc.lon1, tc.lat2, tc.lon2)

			// Use relative tolerance for non-zero distances
			var difference float64
			if tc.expected == 0 {
				difference = math.Abs(result)
			} else {
				difference = math.Abs(result-tc.expected) / tc.expected
			}

			if difference > tc.tolerance {
				t.Errorf("HaversineDistance(%f, %f, %f, %f) = %f, expected %f ± %.1f%%",
					tc.lat1, tc.lon1, tc.lat2, tc.lon2, result, tc.expected, tc.tolerance*100)
			}
		})
	}
}

func TestBoundingBox(t *testing.T) {
	t.Run("Creation and extension", func(t *testing.T) {
		bbox := NewBoundingBox()

		// Check initial state
		if bbox.MinLat != 90.0 || bbox.MinLon != 180.0 || bbox.MaxLat != -90.0 || bbox.MaxLon != -180.0 {
			t.Errorf("NewBoundingBox() incorrect initial state: %+v", bbox)
		}

		// Extend with a single point
		bbox.ExtendWithPoint(37.7749, -122.4194) // San Francisco

		if bbox.MinLat != 37.7749 || bbox.MaxLat != 37.7749 ||
			bbox.MinLon != -122.4194 || bbox.MaxLon != -122.4194 {
			t.Errorf("ExtendWithPoint didn't set values correctly with single point: %+v", bbox)
		}

		// Extend with another point
		bbox.ExtendWithPoint(40.7128, -74.0060) // New York

		if bbox.MinLat != 37.7749 || bbox.MaxLat != 40.7128 ||
			bbox.MinLon != -122.4194 || bbox.MaxLon != -74.0060 {
			t.Errorf("ExtendWithPoint didn't extend correctly with second point: %+v", bbox)
		}

		// Add point that should be ignored (already contained)
		bbox.ExtendWithPoint(39.0, -100.0) // Somewhere in the middle

		if bbox.MinLat != 37.7749 || bbox.MaxLat != 40.7128 ||
			bbox.MinLon != -122.4194 || bbox.MaxLon != -74.0060 {
			t.Errorf("ExtendWithPoint changed bounding box when it shouldn't have: %+v", bbox)
		}
	})

	t.Run("Buffer", func(t *testing.T) {
		bbox := NewBoundingBox()
		bbox.ExtendWithPoint(37.7749, -122.4194) // San Francisco

		// Add a 10km buffer
		original := *bbox
		bbox.Buffer(10000)

		// Check that the buffer was added (approximately 0.09 degrees ≈ 10km)
		bufferDegrees := 10000.0 / 111000.0 // Approximate conversion

		if math.Abs((bbox.MinLat-original.MinLat)+bufferDegrees) > 0.001 ||
			math.Abs((original.MaxLat-bbox.MaxLat)+bufferDegrees) > 0.001 ||
			math.Abs((bbox.MinLon-original.MinLon)+bufferDegrees) > 0.001 ||
			math.Abs((original.MaxLon-bbox.MaxLon)+bufferDegrees) > 0.001 {
			t.Errorf("Buffer didn't add correctly. Original: %+v, Buffered: %+v", original, bbox)
		}
	})

	t.Run("Boundary clipping", func(t *testing.T) {
		bbox := NewBoundingBox()
		bbox.ExtendWithPoint(89.0, 179.0)

		// Add a large buffer that should be clipped
		bbox.Buffer(1000000) // 1000km, should hit boundaries

		if bbox.MinLat < -90.0 || bbox.MaxLat > 90.0 ||
			bbox.MinLon < -180.0 || bbox.MaxLon > 180.0 {
			t.Errorf("Buffer didn't clip to valid boundaries: %+v", bbox)
		}
	})

	t.Run("String format", func(t *testing.T) {
		bbox := NewBoundingBox()
		bbox.ExtendWithPoint(37.7749, -122.4194)
		bbox.ExtendWithPoint(40.7128, -74.0060)

		expected := "(37.774900,-122.419400,40.712800,-74.006000)"
		if bbox.String() != expected {
			t.Errorf("String() = %s, expected %s", bbox.String(), expected)
		}
	})
}
