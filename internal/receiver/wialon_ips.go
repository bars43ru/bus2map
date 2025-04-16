// Package receiver содержит компоненты для приема и обработки данных GPS от различных протоколов.
// Предоставляет интерфейсы и функции для работы с данными телеметрии.
package receiver

import (
	"context"
	"fmt"
	"io"

	"github.com/bars43ru/bus2map/internal/model"
	"github.com/bars43ru/bus2map/internal/protocols/wialonips"
	"github.com/bars43ru/bus2map/pkg/tcp"
)

// BridgeWialonIPS создает обработчик соединения для протокола WialonIPS.
// Преобразует данные из формата WialonIPS в внутренний формат GPS и передает их в GPSLocator.
// Выполняет конвертацию координат из формата WialonIPS в WGS84.
//
// Параметры:
//   - gpsLocator: интерфейс для обработки данных GPS
//
// Возвращает:
//   - tcp.ConnectionHandlerFunc: функция-обработчик TCP-соединения
//
// Ошибки:
//   - Возвращает ошибку при неудачной инициализации парсера
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
