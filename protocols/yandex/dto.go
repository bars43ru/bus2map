package yandex

import (
	"time"
)

// GpsSignal категория GPS-сигнала:
//
//	⦁ Slow - "медленный" (устанавливается для треков общественного транспорта);
//	⦁ Normal - обычный.
type GpsSignal string

const (
	SlowGpsSignal   GpsSignal = "s"
	NormalGpsSignal GpsSignal = "n"
)

// VehicleType тип общественного транспортного средства:
//
//	⦁ Bus - автобус;
//	⦁ Trolleybus - троллейбус;
//	⦁ Tramway - трамвай;
//	⦁ Minibus - маршрутное такси.
type VehicleType string

const (
	BusVehicleType        VehicleType = "bus"
	TrolleybusVehicleType VehicleType = "trolleybus"
	TramwayVehicleType    VehicleType = "tramway"
	MinibusVehicleType    VehicleType = "minibus"
)

// tracks пакет передаваемых данных.
type tracks struct {
	// Clid идентификатор участника программы.
	//	 Длина идентификатора не должна превышать 32 символа и содержать только символы латинского алфавита и цифры.
	Clid string `xml:"clid,attr"`
	// Tracks перечень машин по которым будем высылать информацию
	Tracks []Track `xml:"track"`
}

// Track данные о транспортном средстве и маршруте по которому он движется.
type Track struct {
	// Uuid идентификатор движущегося объекта (транспортного средства).
	//  Длина идентификатора не должна превышать 32 символа и содержать только символы латинского алфавита и цифры.
	//  uuid=0d63b6deacb91b00e46194fac325b72a
	UUID string `xml:"uuid,attr"`
	// Category категория GPS-сигнала
	Category GpsSignal `xml:"category,attr"`
	// Route идентификатор маршрута.
	//	route=190Б
	Route string `xml:"route,attr"`
	// VehicleType тип общественного транспортного средства:
	VehicleType VehicleType `xml:"vehicle_type,attr"`
	// Point данные об последнем актуальном местоположении данного транспортного средства
	Point Point `xml:"point"`
}

// Point данные о местоположении общественного транспорта (только для треков общественного транспорта).
type Point struct {
	// Latitude долгота точки в градусах. В качестве десятичного разделителя используется точка.
	//	longitude=37.620070
	Latitude float64 `xml:"latitude,attr"`
	// Longitude широта точки в градусах. В качестве десятичного разделителя используется точка.
	//	latitude=55.753630
	Longitude float64 `xml:"longitude,attr"`
	// AvgSpeed мгновенная скорость транспортного средства, полученная от приемника GPS, км/ч.
	//	avg_speed=53
	AvgSpeed uint `xml:"avg_speed,attr"`
	// Direction направление движения в градусах (направление на север - 0 градусов). Диапазон значений 0-360.
	// 	direction=242
	Direction uint `xml:"direction,attr"`
	// Time дата и время получения координат точки от GPS-приемника (по Гринвичу). Формат: ДДММГГГГ:ччммсс
	//  time=10012009:142045
	Time CustomTime `xml:"time,attr"`
}

type CustomTime time.Time

func (t CustomTime) MarshalText() ([]byte, error) {
	return ([]byte)((time.Time(t)).Format("02012006:150405")), nil
}
