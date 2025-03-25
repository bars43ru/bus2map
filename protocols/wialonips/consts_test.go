package wialonips

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_layoutTime(t *testing.T) {
	expect := time.Date(2021, 5, 6, 8, 16, 06, 0, time.UTC)
	actual, err := time.Parse(layoutTime, "060521081606")
	require.NoError(t, err)
	require.Equal(t, expect, actual)

	_, err = time.Parse(layoutTime, "062421081606")
	require.Error(t, err)
}
