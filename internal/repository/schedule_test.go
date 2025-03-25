package repository_test

import (
	"github.com/bars43ru/bus2map/internal/repository"
	"testing"
	"time"
)

func TestSchedule_ParseDateTime(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Time
		wantErr  bool
	}{
		{
			input:    "02/06/2020T12:55:00Z+03:00",
			expected: time.Date(2020, 6, 2, 12, 55, 0, 0, time.FixedZone("UTC+3", 3*60*60)),
			wantErr:  false,
		},
		{
			input:   "invalid-date",
			wantErr: true,
		},
		{
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		got, err := (*repository.Schedule)(nil).ParseDateTime(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("parseDateTime(%q): expected error, got nil", tt.input)
			}
			continue
		}

		if err != nil {
			t.Errorf("parseDateTime(%q): unexpected error: %v", tt.input, err)
			continue
		}

		if !got.Equal(tt.expected) {
			t.Errorf("parseDateTime(%q): expected %v, got %v", tt.input, tt.expected, got)
		}
	}
}
