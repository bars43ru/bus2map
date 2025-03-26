package repository_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bars43ru/bus2map/internal/model"
	"github.com/bars43ru/bus2map/internal/repository"
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

func TestSchedule_GetCurrent(t *testing.T) {
	schedule := repository.NewSchedule("")
	schedule.Replace([]model.Schedule{
		{
			"4",
			"3",
			time.Date(2025, 0o1, 0o1, 0o1, 0o1, 0o1, 0o0, time.UTC),
			time.Date(2025, 0o1, 0o1, 0o1, 0o1, 0o3, 0o0, time.UTC),
		},
		{
			"1",
			"2",
			time.Date(2025, 0o1, 0o1, 0o1, 0o1, 0o1, 0o0, time.UTC),
			time.Date(2025, 0o1, 0o1, 0o1, 0o1, 0o3, 0o0, time.UTC),
		},
		{
			"3",
			"2",
			time.Date(2025, 0o1, 0o1, 0o1, 0o1, 0o4, 0o0, time.UTC),
			time.Date(2025, 0o1, 0o1, 0o1, 0o1, 10, 0o0, time.UTC),
		},
	})
	t.Run("before period",
		func(t *testing.T) {
			_, err := schedule.GetCurrent("2", time.Date(2025, 0o1, 0o1, 0o1, 0o1, 0o0, 0o0, time.UTC))
			require.ErrorIs(t, err, repository.ErrNotFound)
		},
	)
	t.Run("in period",
		func(t *testing.T) {
			v, err := schedule.GetCurrent("2", time.Date(2025, 0o1, 0o1, 0o1, 0o1, 0o2, 0o0, time.UTC))
			require.NoError(t, err)
			require.Equal(t, v.Number.String(), "1")

			v, err = schedule.GetCurrent("2", time.Date(2025, 0o1, 0o1, 0o1, 0o1, 0o6, 0o0, time.UTC))
			require.NoError(t, err)
			require.Equal(t, v.Number.String(), "3")
		},
	)
	t.Run("after period",
		func(t *testing.T) {
			_, err := schedule.GetCurrent("2", time.Date(2025, 0o1, 0o1, 0o1, 0o1, 11, 0o0, time.UTC))
			require.ErrorIs(t, err, repository.ErrNotFound)
		},
	)
	t.Run("boundary period",
		func(t *testing.T) {
			v, err := schedule.GetCurrent("3", time.Date(2025, 0o1, 0o1, 0o1, 0o1, 0o1, 0o0, time.UTC))
			require.NoError(t, err)
			require.Equal(t, v.Number.String(), "4")

			v, err = schedule.GetCurrent("3", time.Date(2025, 0o1, 0o1, 0o1, 0o1, 0o3, 0o0, time.UTC))
			require.NoError(t, err)
			require.Equal(t, v.Number.String(), "4")
		},
	)
	t.Run("not found state number",
		func(t *testing.T) {
			_, err := schedule.GetCurrent("25", time.Date(2025, 0o1, 0o1, 0o1, 0o1, 0o2, 0o0, time.UTC))
			require.ErrorIs(t, err, repository.ErrNotFound)
		},
	)
}
