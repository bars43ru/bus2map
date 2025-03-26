package receiver

import (
	"context"
	"fmt"
	"io"

	"github.com/bars43ru/bus2map/internal/model"
	"github.com/bars43ru/bus2map/protocols/tcp"
	"github.com/bars43ru/bus2map/protocols/wialonips"
)

func BridgeWialonIPS(gpsLocator GPSLocator) tcp.ConnectionHandlerFunc {
	return func(ctx context.Context, r io.Reader) error {
		datasource, err := wialonips.NewParse(r)
		if err != nil {
			return fmt.Errorf("new parse WialonIPS: %w", err)
		}
		for _, point := range datasource.Points(ctx) {
			rawGPS := model.GPS{
				UID:       point.UID,
				Time:      point.Time,
				Latitude:  point.Latitude.ToWgs84(),
				Longitude: point.Longitude.ToWgs84(),
				Speed:     uint32(point.Speed),
				Course:    uint32(point.Course),
			}
			gpsLocator.ProcessGPSData(ctx, rawGPS)
		}
		return nil
	}
}
