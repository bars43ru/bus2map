package model

import (
	"time"

	"github.com/bars43ru/bus2map/internal/model/transport_type"
)

type Route struct {
	Number       RouteNumber
	YandexNumber string
	TwoGISNumber string
}

type Transport struct {
	GUID        string
	StateNumber StateNumber
	Type        transport_type.Type
}

type Schedule struct {
	Number      RouteNumber
	StateNumber StateNumber
	From        time.Time
	To          time.Time
}

// BusTrackingInfo содержит информацию о маршруте, текущем положении и активном расписании автобуса
type BusTrackingInfo struct {
	Route     Route     // Информация о маршруте
	Transport Transport // Информация об автобусе
	Location  GPS       // Текущие GPS-координаты
	Schedule  Schedule  // Данные из расписания, по которому автобус движется в данный момент
}

type GPS struct {
	UID       string    // идентификатор транспортного средства в системе мониторинга
	Time      time.Time // дата и время сообщения
	Latitude  float64   // широта
	Longitude float64   // долгота
	Speed     uint32    // скорость
	Course    uint32    // курс
}
