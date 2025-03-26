package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/imkira/go-observer/v2"

	"github.com/bars43ru/bus2map/internal/model"
	"github.com/bars43ru/bus2map/internal/repository"
)

type BusTracking struct {
	location  observer.Property[*model.BusTrackingInfo]
	route     *repository.Route
	transport *repository.Transport
	schedule  *repository.Schedule
}

func New(
	route *repository.Route,
	transport *repository.Transport,
	schedule *repository.Schedule,
) *BusTracking {
	return &BusTracking{
		location:  observer.NewProperty[*model.BusTrackingInfo](nil),
		route:     route,
		transport: transport,
		schedule:  schedule,
	}
}

func (s *BusTracking) SubscribeLocation() observer.Stream[*model.BusTrackingInfo] {
	return s.location.Observe()
}

func (s *BusTracking) ProcessGPSData(ctx context.Context, gpsData model.GPS) {
	transport, err := s.transport.Get(gpsData.UID)
	if err != nil {
		l := slog.With(slog.String("uid", gpsData.UID))
		if errors.Is(err, repository.ErrNotFound) {
			l.WarnContext(ctx, "not found UID in transport")
			// TODO тут записываем что есть какой то UID сигнала по которому нет данных о привязке к транспорту
			return
		}
		l.ErrorContext(ctx, "get transport from UID", slog.Any("error", err))
		return
	}
	schedule, err := s.schedule.GetCurrent(transport.StateNumber, gpsData.Time)
	if err != nil {
		l := slog.With(
			slog.String("state_number", transport.StateNumber.String()),
			slog.Time("gps_time", gpsData.Time),
		)
		if errors.Is(err, repository.ErrNotFound) {
			l.WarnContext(ctx, "not found schedule for transport")
			// TODO тут записываем что есть какой то UID сигнала по которому нет данных о привязке к транспорту
			return
		}
		l.ErrorContext(ctx, "get schedule for transport", slog.Any("error", err))
		return
	}

	route, err := s.route.GetRoute(schedule.Number)
	if err != nil {
		l := slog.With(
			slog.String("route_number", schedule.Number.String()),
		)
		if errors.Is(err, repository.ErrNotFound) {
			l.WarnContext(ctx, "not found route")
			// TODO тут записываем что есть какой то UID сигнала по которому нет данных о привязке к транспорту
			return
		}
		l.ErrorContext(ctx, "get route", slog.Any("error", err))
		return
	}

	s.location.Update(&model.BusTrackingInfo{
		Route:     route,
		Transport: transport,
		Location:  gpsData,
		Schedule:  schedule,
	})
}
