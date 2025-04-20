package service

import (
	"context"
	"fmt"
	"github.com/bars43ru/bus2map/internal/ctype"
	"github.com/bars43ru/bus2map/internal/protocols/wialonips"
	"github.com/bars43ru/bus2map/pkg/tcp"
	"github.com/imkira/go-observer/v2"
	"io"
)

func NewReciver(tracker observer.Property[ctype.RawGPSData]) tcp.ConnectionHandlerFunc {
	return func(ctx context.Context, r io.Reader) error {
		datasource, err := wialonips.NewParse(r)
		if err != nil {
			return fmt.Errorf("new parse WialonIPS: %w", err)
		}
		for _, point := range datasource.Points(ctx) {
			rawGPS := ctype.RawGPSData{
				UID:       point.UID,
				Time:      point.Time,
				Latitude:  point.Latitude.ToWgs84(),
				Longitude: point.Longitude.ToWgs84(),
				Speed:     uint32(point.Speed),
				Course:    uint32(point.Course),
			}
			tracker.Update(rawGPS)
		}
		return nil
	}
}
