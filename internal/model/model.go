// Package model содержит основные структуры данных для работы с транспортом и маршрутами.
// Определяет типы данных для хранения информации о маршрутах, транспорте, расписании и отслеживании.
package model

import (
	"time"

	"github.com/bars43ru/bus2map/internal/model/transport_type"
)

// Route представляет информацию о маршруте общественного транспорта.
// Содержит номера маршрута в различных системах (внутренний, Яндекс, 2ГИС).
type Route struct {
	Number       RouteNumber
	YandexNumber string
	TwoGISNumber string
}

// Transport представляет информацию о транспортном средстве.
// Содержит уникальный идентификатор, государственный номер и тип транспорта.
type Transport struct {
	GUID        string
	StateNumber StateNumber
	Type        transport_type.Type
}

// Schedule представляет информацию о расписании движения транспорта.
// Содержит номер маршрута, государственный номер транспорта и временной интервал работы.
type Schedule struct {
	Number      RouteNumber
	StateNumber StateNumber
	From        time.Time
	To          time.Time
}

// BusTrackingInfo содержит полную информацию о движущемся транспортном средстве.
// Объединяет данные о маршруте, транспорте, текущем местоположении и активном расписании.
type BusTrackingInfo struct {
	Route     Route     // Информация о маршруте
	Transport Transport // Информация об автобусе
	Location  GPS       // Текущие GPS-координаты
	Schedule  Schedule  // Данные из расписания, по которому автобус движется в данный момент
}

// GPS представляет данные о местоположении транспортного средства.
// Содержит информацию о координатах, скорости и курсе движения.
type GPS struct {
	UID       string    // Идентификатор транспортного средства в системе мониторинга
	Time      time.Time // Дата и время получения координат
	Latitude  float64   // Широта в градусах (WGS84)
	Longitude float64   // Долгота в градусах (WGS84)
	Speed     uint32    // Скорость движения в км/ч
	Course    uint32    // Курс движения в градусах (0-359)
}
