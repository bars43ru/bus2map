// Package wialonips реализует парсинг протокола WialonIPS для получения GPS-данных от трекеров.
// Протокол используется для передачи телеметрии в формате:
// #L#<IMEI>;<Data>#<CRC>#
// Где:
// - IMEI - уникальный идентификатор устройства
// - Data - данные в формате CSV
// - CRC - контрольная сумма

package wialonips

import (
	"math"
	"time"
)

// Coordinate представляет координату в формате протокола WialonIPS.
// Координаты передаются в формате DDMM.MMMM, где:
// - DD - градусы
// - MM.MMMM - минуты с десятичной дробью
type Coordinate float64

// ToWgs84 конвертирует координату из формата WialonIPS в стандартный формат WGS84.
// Формат WialonIPS: DDMM.MMMM
// Формат WGS84: DD.DDDDDD
//
//nolint:gomnd // перевод координат в wgs84
func (c Coordinate) ToWgs84() float64 {
	ratio := math.Pow(10, 6)
	degrees := math.Trunc(float64(c) / 100)
	remain := float64(c) - degrees*100
	return degrees + math.Round(remain/60*ratio)/ratio
}

// messageL представляет заголовок сообщения WialonIPS.
// Формат: #L#<IMEI>;
type messageL struct {
	// UID идентификатор транспортного средства в системе мониторинга
	UID string
}

// messageD представляет данные GPS в сообщении WialonIPS.
// Формат: #D#<date>;<time>;<lat1>;<lat2>;<lon1>;<lon2>;<speed>;<course>;<alt>;<sats>;
type messageD struct {
	// Time дата и время сообщения в формате DDMMYYHHMMSS
	Time time.Time
	/// Latitude широта в формате DDMM.MMMM
	Latitude Coordinate
	// Longitude долгота в формате DDMM.MMMM
	Longitude Coordinate
	// Speed Скорость в км/ч
	Speed uint
	// Course курс в градусах (0-359)
	Course uint8
	// Alt высота. Если отсутствует, значение null.
	// Alt *int
	// Sats Количество спутников. Если отсутствует, значение null.
	// Sats *int
}

// Point объединяет заголовок и данные GPS в единую структуру.
// Используется для передачи полной информации о местоположении транспортного средства.
type Point struct {
	messageL
	messageD
}
