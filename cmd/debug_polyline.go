package main

import (
	"fmt"
	"os"

	"github.com/NERVsystems/osmmcp/pkg/geo"
	"github.com/NERVsystems/osmmcp/pkg/osm"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: debug_polyline <encoded_polyline>")
		os.Exit(1)
	}
	encoded := os.Args[1]
	points := osm.DecodePolyline(encoded)
	for i, pt := range points {
		fmt.Printf("Decoded Point %d: Latitude: %.8f, Longitude: %.8f\n", i, pt.Latitude, pt.Longitude)
	}

	// Now encode the test point and print the encoded string
	testPt := geo.Location{Latitude: -25.363882, Longitude: 131.044922}
	latE5 := int(testPt.Latitude * 1e5)
	lngE5 := int(testPt.Longitude * 1e5)
	latE6 := int(testPt.Latitude * 1e6)
	lngE6 := int(testPt.Longitude * 1e6)
	fmt.Printf("\nTest Point: {Latitude: %.8f, Longitude: %.8f}\n", testPt.Latitude, testPt.Longitude)
	fmt.Printf("As integers (1e5): latE5 = %d, lngE5 = %d\n", latE5, lngE5)
	fmt.Printf("As integers (1e6): latE6 = %d, lngE6 = %d\n", latE6, lngE6)
	encodedTestE5 := osm.EncodePolyline([]geo.Location{testPt})
	fmt.Printf("Encoded string (1e5): %s\n", encodedTestE5)
	// Try encoding with 1e6 scaling
	fakePt := geo.Location{Latitude: float64(latE6) / 1e5, Longitude: float64(lngE6) / 1e5}
	encodedTestE6 := osm.EncodePolyline([]geo.Location{fakePt})
	fmt.Printf("Encoded string (1e6 as 1e5): %s\n", encodedTestE6)
}
