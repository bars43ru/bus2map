package wialonips

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoordinate_ToWgs84(t *testing.T) {
	tests := []struct {
		name  string
		value Coordinate
		wgs84 float64
	}{
		{
			name:  "Latitude success",
			value: 5844.6826,
			wgs84: 58.74471,
		},
		{
			name:  "Longitude success",
			value: 05010.7126,
			wgs84: 50.178543,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.value.ToWgs84(), tt.wgs84)
		})
	}
}
