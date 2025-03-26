package controller

import (
	"io"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/bars43ru/bus2map/api/bustracking"
	"github.com/bars43ru/bus2map/internal/model"
	"github.com/bars43ru/bus2map/internal/model/transport_type"
	"github.com/bars43ru/bus2map/internal/service"
)

var (
	_TransportTypeToPbTransportType = map[transport_type.Type]pb.Transport_Type{
		transport_type.TypeBUS:        pb.Transport_BUS,
		transport_type.TypeMINIBUS:    pb.Transport_MINIBUS,
		transport_type.TypeTRAMWAY:    pb.Transport_TRAMWAY,
		transport_type.TypeTROLLEYBUS: pb.Transport_TROLLEYBUS,
	}
	_PbTransportTypeToTransportType = map[pb.Transport_Type]transport_type.Type{
		pb.Transport_BUS:        transport_type.TypeBUS,
		pb.Transport_MINIBUS:    transport_type.TypeMINIBUS,
		pb.Transport_TRAMWAY:    transport_type.TypeTRAMWAY,
		pb.Transport_TROLLEYBUS: transport_type.TypeTROLLEYBUS,
	}
)

type BusTracking struct {
	pb.UnsafeBusTrackingServiceServer
	service *service.BusTracking
}

func NewBusTrackingService(service *service.BusTracking) *BusTracking {
	return &BusTracking{
		service: service,
	}
}

func (s *BusTracking) StreamGPSData(stream grpc.ClientStreamingServer[pb.GPSData, pb.StreamGPSDataResponse]) error {
	ctx := stream.Context()
	for {
		// Получаем данные от клиента
		pbGPSData, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamGPSDataResponse{})
		}
		if err != nil {
			slog.ErrorContext(ctx, "receiving GPS data in StreamRawGPSData", slog.Any("error", err))
			return err
		}
		gpsData := model.GPS{
			UID:       pbGPSData.GetUid(),
			Time:      pbGPSData.Time.AsTime(),
			Latitude:  pbGPSData.Latitude,
			Longitude: pbGPSData.Longitude,
			Speed:     pbGPSData.Speed,
			Course:    pbGPSData.Course,
		}
		s.service.ProcessGPSData(ctx, gpsData)
	}
}

func (s *BusTracking) StreamBusTrackingInfo(
	req *pb.StreamBusDataRequest,
	stream grpc.ServerStreamingServer[pb.BusTrackingInfo],
) error {
	ctx := stream.Context()
	watcher := s.service.SubscribeLocation()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-watcher.Changes():
			busTrackingInfo := watcher.Next()
			pbBusTrackingInfo := &pb.BusTrackingInfo{
				GpsData:   s.gpsDataToPbGPSData(busTrackingInfo.Location),
				Route:     s.routeToPbRoute(busTrackingInfo.Route),
				Transport: s.transportToPbTransport(busTrackingInfo.Transport),
				Schedule:  s.scheduleToPbSchedule(busTrackingInfo.Schedule),
			}
			err := stream.Send(pbBusTrackingInfo)
			if err != nil {
				slog.ErrorContext(ctx, "sending BusTrackingInfo to subscribe client", slog.Any("error", err))
				return err
			}
		}
	}
}

func (s *BusTracking) gpsDataToPbGPSData(gps model.GPS) *pb.GPSData {
	return &pb.GPSData{
		Uid:       gps.UID,
		Latitude:  gps.Latitude,
		Longitude: gps.Longitude,
		Speed:     gps.Speed,
		Course:    gps.Course,
		Time:      timestamppb.New(gps.Time),
	}
}

func (s *BusTracking) routeToPbRoute(route model.Route) *pb.Route {
	return &pb.Route{
		Number: route.Number.String(),
		Yandex: route.YandexNumber,
		TwoGis: route.TwoGISNumber,
	}
}

func (s *BusTracking) transportToPbTransport(transport model.Transport) *pb.Transport {
	return &pb.Transport{
		Uuid:        transport.GUID,
		StateNumber: transport.StateNumber.String(),
		Type:        _TransportTypeToPbTransportType[transport.Type],
	}
}

func (s *BusTracking) scheduleToPbSchedule(schedule model.Schedule) *pb.Schedule {
	return &pb.Schedule{
		Number:      schedule.Number.String(),
		StateNumber: schedule.StateNumber.String(),
		From:        timestamppb.New(schedule.From),
		To:          timestamppb.New(schedule.To),
	}
}
