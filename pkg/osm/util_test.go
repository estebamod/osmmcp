package osm

import (
	"testing"
)

func TestValidateCoords(t *testing.T) {
	tests := []struct {
		name    string
		lat     float64
		lon     float64
		wantErr bool
	}{
		{
			name:    "valid coordinates",
			lat:     40.7128,
			lon:     -74.0060,
			wantErr: false,
		},
		{
			name:    "valid coordinates at boundaries",
			lat:     90.0,
			lon:     180.0,
			wantErr: false,
		},
		{
			name:    "valid coordinates at negative boundaries",
			lat:     -90.0,
			lon:     -180.0,
			wantErr: false,
		},
		{
			name:    "invalid latitude too high",
			lat:     91.0,
			lon:     -74.0060,
			wantErr: true,
		},
		{
			name:    "invalid latitude too low",
			lat:     -91.0,
			lon:     -74.0060,
			wantErr: true,
		},
		{
			name:    "invalid longitude too high",
			lat:     40.7128,
			lon:     181.0,
			wantErr: true,
		},
		{
			name:    "invalid longitude too low",
			lat:     40.7128,
			lon:     -181.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCoords(tt.lat, tt.lon)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCoords() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
