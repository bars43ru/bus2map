package receiver

import (
	"context"
	"fmt"
	"github.com/bars43ru/bus2map/internal/model"
	"github.com/bars43ru/bus2map/protocols/egts"
	"github.com/bars43ru/bus2map/protocols/tcp"
	"io"
)

func BridgeEGTS(gpsLocator GPSLocator) tcp.ConnectionHandlerFunc {
	return func(ctx context.Context, r io.Reader) error {
		datasource := egts.NewParse(r)
		for _, point := range datasource.Points(ctx) {
			rawGPS := model.GPS{
				UID:       fmt.Sprintf("%d", point.Client),
				Time:      point.Time,
				Latitude:  point.Latitude,
				Longitude: point.Longitude,
				Speed:     uint32(point.Speed),
				Course:    uint32(point.Course),
			}
			gpsLocator.ProcessGPSData(ctx, rawGPS)
		}
		return nil
	}
}
