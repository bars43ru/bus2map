// Package receiver содержит компоненты для приема и обработки данных GPS от различных протоколов.
// Предоставляет интерфейсы и функции для работы с данными телеметрии.
package receiver

import (
	"context"

	"github.com/bars43ru/bus2map/internal/model"
)

// GPSLocator определяет интерфейс для обработки данных GPS.
// Используется для передачи данных о местоположении транспортных средств в систему.
type GPSLocator interface {
	// ProcessGPSData обрабатывает данные о местоположении транспортного средства.
	// Принимает контекст и структуру GPS с информацией о координатах, скорости и курсе.
	ProcessGPSData(ctx context.Context, gps model.GPS)
}
