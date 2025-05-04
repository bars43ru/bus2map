package service

import (
	"context"
	"fmt"
	"github.com/bars43ru/bus2map/internal/ctype"
	"github.com/bars43ru/bus2map/internal/protocols/egts"
	"github.com/bars43ru/bus2map/pkg/tcp"
	"github.com/imkira/go-observer/v2"
	"io"
)

func NewReciver(tracker observer.Property[ctype.RawGPSData]) tcp.ConnectionHandlerFunc {
	return func(ctx context.Context, r io.Reader) error {
		datasource := egts.NewParse(r)
		for _, point := range datasource.Points(ctx) {
			rawGPS := ctype.RawGPSData{
				UID:       fmt.Sprintf("%d", point.Client),
				Time:      point.Time,
				Latitude:  point.Latitude,
				Longitude: point.Longitude,
				Speed:     uint32(point.Speed),
				Course:    uint32(point.Course),
			}
			tracker.Update(rawGPS)
		}
		return nil
	}
}
