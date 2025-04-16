// Package receiver содержит компоненты для приема и обработки данных GPS от различных протоколов.
// Предоставляет интерфейсы и функции для работы с данными телеметрии.
package receiver

import (
	"context"
	"fmt"
	"io"

	"github.com/bars43ru/bus2map/internal/model"
	"github.com/bars43ru/bus2map/internal/protocols/egts"
	"github.com/bars43ru/bus2map/pkg/tcp"
)

// BridgeEGTS создает обработчик соединения для протокола EGTS.
// Преобразует данные из формата EGTS в внутренний формат GPS и передает их в GPSLocator.
//
// Параметры:
//   - gpsLocator: интерфейс для обработки данных GPS
//
// Возвращает:
//   - tcp.ConnectionHandlerFunc: функция-обработчик TCP-соединения
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
