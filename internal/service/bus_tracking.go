// Package service содержит основные сервисы приложения.
// Реализует бизнес-логику обработки данных и управления состоянием системы.
package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/imkira/go-observer/v2"

	"github.com/bars43ru/bus2map/internal/model"
	"github.com/bars43ru/bus2map/internal/repository"
	"github.com/bars43ru/bus2map/pkg/xslog"
)

// BusTracking представляет сервис для отслеживания местоположения транспортных средств.
// Обрабатывает GPS-данные, сопоставляет их с информацией о транспорте, маршрутах и расписании.
// Предоставляет возможность подписки на обновления местоположения.
type BusTracking struct {
	location  observer.Property[*model.BusTrackingInfo] // Свойство для хранения и обновления информации о местоположении
	route     *repository.Route                         // Репозиторий для работы с маршрутами
	transport *repository.Transport                     // Репозиторий для работы с транспортными средствами
	schedule  *repository.Schedule                      // Репозиторий для работы с расписанием
}

// New создает новый экземпляр сервиса BusTracking.
// Инициализирует сервис с указанными репозиториями.
//
// Параметры:
//   - route: репозиторий для работы с маршрутами
//   - transport: репозиторий для работы с транспортными средствами
//   - schedule: репозиторий для работы с расписанием
//
// Возвращает:
//   - *BusTracking: указатель на созданный сервис
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

// SubscribeLocation возвращает поток для подписки на обновления местоположения.
// Позволяет получать уведомления об изменении местоположения транспортных средств.
//
// Возвращает:
//   - observer.Stream[*model.BusTrackingInfo]: поток обновлений местоположения
func (s *BusTracking) SubscribeLocation() observer.Stream[*model.BusTrackingInfo] {
	return s.location.Observe()
}

// ProcessGPSData обрабатывает полученные GPS-данные.
// Сопоставляет данные с информацией о транспорте, маршруте и расписании.
// Обновляет информацию о местоположении и уведомляет подписчиков.
//
// Параметры:
//   - ctx: контекст для управления жизненным циклом
//   - gpsData: данные о местоположении транспортного средства
func (s *BusTracking) ProcessGPSData(ctx context.Context, gpsData model.GPS) {
	transport, err := s.transport.Get(gpsData.UID)
	if err != nil {
		l := slog.With(slog.String("uid", gpsData.UID))
		if errors.Is(err, repository.ErrNotFound) {
			l.WarnContext(ctx, "not found UID in transport")
			// TODO тут записываем что есть какой то UID сигнала по которому нет данных о привязке к транспорту
			return
		}
		l.ErrorContext(ctx, "get transport from UID", xslog.Error(err))
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
		l.ErrorContext(ctx, "get schedule for transport", xslog.Error(err))
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
		l.ErrorContext(ctx, "get route", xslog.Error(err))
		return
	}

	s.location.Update(&model.BusTrackingInfo{
		Route:     route,
		Transport: transport,
		Location:  gpsData,
		Schedule:  schedule,
	})
}
