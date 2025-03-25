package sender

import (
	"context"
	"github.com/bars43ru/bus2map/internal/model"
	"github.com/bars43ru/bus2map/internal/model/transport_type"
	"github.com/bars43ru/bus2map/protocols/yandex"
	"github.com/imkira/go-observer/v2"
	"log/slog"
	"time"
)

var _TransportTypeToVehicleType = map[transport_type.Type]yandex.VehicleType{
	transport_type.TypeBUS:        yandex.BusVehicleType,
	transport_type.TypeTROLLEYBUS: yandex.TrolleybusVehicleType,
	transport_type.TypeTRAMWAY:    yandex.TramwayVehicleType,
	transport_type.TypeMINIBUS:    yandex.MinibusVehicleType,
}

func BridgeYandex(
	cliYandex yandex.Client,
	observer observer.Stream[*model.BusTrackingInfo],
) func(ctx context.Context) error {
	const sizeChunk = 50
	makeChunk := func(ctx context.Context) []model.BusTrackingInfo {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		chunk := make([]model.BusTrackingInfo, 0, sizeChunk)
	L:
		for {
			select {
			case <-observer.Changes():
				busTrackingInfo := observer.Next()
				if busTrackingInfo == nil {
					continue
				}
				chunk = append(chunk, *busTrackingInfo)
				if len(chunk) >= 50 {
					break L
				}
			case <-ctx.Done():
				break L
			}
		}
		return chunk
	}

	send := func(ctx context.Context, tracks []yandex.Track) error {
		ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
		defer cancel()
		return cliYandex.Send(ctx, tracks)
	}

	return func(ctx context.Context) error {
		for ctx.Err() == nil {
			busTrackingInfoItems := makeChunk(ctx)
			if len(busTrackingInfoItems) == 0 {
				continue
			}
			tracks := make([]yandex.Track, 0, len(busTrackingInfoItems))
			for _, busTrackingInfo := range busTrackingInfoItems {
				track := yandex.Track{
					UUID:        busTrackingInfo.Transport.StateNumber.String(),
					Category:    yandex.NormalGpsSignal,
					Route:       busTrackingInfo.Route.YandexNumber,
					VehicleType: _TransportTypeToVehicleType[busTrackingInfo.Transport.Type],
					Point: yandex.Point{
						Latitude:  busTrackingInfo.Location.Latitude,
						Longitude: busTrackingInfo.Location.Longitude,
						AvgSpeed:  uint(busTrackingInfo.Location.Speed),
						Direction: uint(busTrackingInfo.Location.Course),
						Time:      yandex.CustomTime(busTrackingInfo.Location.Time),
					},
				}
				tracks = append(tracks, track)
			}

			if err := send(ctx, tracks); err != nil {
				slog.ErrorContext(ctx, "send tracks to yandex", slog.Any("error", err))
			}
		}
		return nil
	}
}
