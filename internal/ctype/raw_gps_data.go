package ctype

import "time"

// RawGPSData представляет необработанные данные о местоположении транспортного средства используемые в получателях от
// провайдеров.
type RawGPSData struct {
	UID       string    // Идентификатор транспортного средства в системе мониторинга
	Time      time.Time // Дата и время получения координат
	Latitude  float64   // Широта в градусах (WGS84)
	Longitude float64   // Долгота в градусах (WGS84)
	Speed     uint32    // Скорость движения в км/ч
	Course    uint32    // Курс движения в градусах (0-359)
}

func (d *RawGPSData) IsEmpty() bool {
	return d.UID == "" && d.Time.IsZero() && d.Latitude == 0 && d.Longitude == 0 && d.Speed == 0 && d.Course == 0
}
