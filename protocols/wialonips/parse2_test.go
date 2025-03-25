package wialonips

import (
	"context"
	"fmt"

	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestExample(t *testing.T) {
	file, err := os.Open("./egts/0010-11:25:07.egts")
	require.NoError(t, err)
	parse, err := NewParse(file)
	require.NoError(t, err)
	for _, point := range parse.Points(context.Background()) {
		fmt.Println(point, point.Latitude.ToWgs84(), point.Longitude.ToWgs84())
	}
}
