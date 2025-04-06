package sender

import (
	"context"
	"log/slog"
	"time"

	"github.com/imkira/go-observer/v2"

	"github.com/bars43ru/bus2map/internal/model"
	"github.com/bars43ru/bus2map/internal/model/transport_type"
	"github.com/bars43ru/bus2map/internal/protocols/yandex"
	"github.com/bars43ru/bus2map/pkg/xslog"
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
					slog.InfoContext(ctx, "data packet has been formed for sending")
					break L
				}
			case <-ctx.Done():
				slog.InfoContext(ctx, "the data packet accumulation time has expired")
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
				slog.InfoContext(ctx, "no data to send")
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

			slog.InfoContext(ctx, "data sending")
			if err := send(ctx, tracks); err != nil {
				slog.ErrorContext(ctx, "send tracks to yandex", xslog.Error(err))
			}
		}
		return nil
	}
}
