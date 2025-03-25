package wialonips

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNew_MessageD(t *testing.T) {
	tests := []struct {
		name    string
		uid     string
		arg     string
		wantMsg *messageD
		wantErr bool
	}{
		{
			name: "success",
			uid:  "353173067939817",
			arg:  "#L#353173067939817;NA",
		},
		{
			name:    "error in messageL",
			uid:     "",
			arg:     "#L353173067939817;NA",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parse, err := NewParse(strings.NewReader(tt.arg))
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, parse)
			} else {
				require.NoError(t, err)
				require.NotNil(t, parse)
				require.Equal(t, tt.uid, parse.uid())
			}
		})
	}
}

func TestNew_Points(t *testing.T) {
	const source = `#L#353173067939817;NA
#D#060521;081606;5844.6826;N;05010.7126;E;8;131;113.000000;15;7.000000;3;NA;NA;;SOS:1:1,avl_driver:3:,Odom:1:851171,Speed:1:9
#D#06rrr0521;081606;5844.6826;N;05010.7126;E;8;131;113.000000;15;7.000000;3;NA;NA;;SOS:1:1,avl_driver:3:,Odom:1:851171,Speed:1:9
#D#060521;081606;90.0;N;0.0;E;8;131;113.000000;15;7.000000;3;NA;NA;;SOS:1:1,avl_driver:3:,Odom:1:851171,Speed:1:9

#L#0eee60521;081606;5844.6826;N;05010.7126;E;8;131;113.000000;15;7.000000;3;NA;NA;;SOS:1:1,avl_driver:3:,Odom:1:851171,Speed:1:9

#D#060521;081606;5844.6826;N;05010.7126;E;24;131;113.000000;15;7.000000;3;NA;NA;;SOS:1:1,avl_driver:3:,Odom:1:851171,Speed:1:9
`
	want := []Point{
		{
			messageL: messageL{
				UID: "353173067939817",
			},
			messageD: messageD{
				Time:      time.Date(2021, 5, 6, 8, 16, 6, 0, time.UTC),
				Latitude:  5844.6826,
				Longitude: 05010.7126,
				Speed:     8,
				Course:    131,
			},
		},
		{
			messageL: messageL{
				UID: "353173067939817",
			},
			messageD: messageD{
				Time:      time.Date(2021, 5, 6, 8, 16, 6, 0, time.UTC),
				Latitude:  5844.6826,
				Longitude: 05010.7126,
				Speed:     24,
				Course:    131,
			},
		},
	}
	parse, err := NewParse(strings.NewReader(source))
	require.NoError(t, err)
	require.NotNil(t, parse)

	for index, point := range parse.Points(context.Background()) {
		if index > len(want) {
			require.ErrorIs(t, err, io.EOF)
			break
		} else {
			require.NoError(t, err)
		}
		require.Equal(t, want[index], point)
	}
}
