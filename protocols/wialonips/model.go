package wialonips

import (
	"math"
	"time"
)

type Coordinate float64

//nolint:gomnd // перевод координат в wgs84
func (c Coordinate) ToWgs84() float64 {
	ratio := math.Pow(10, 6)
	degrees := math.Trunc(float64(c) / 100)
	remain := float64(c) - degrees*100
	return degrees + math.Round(remain/60*ratio)/ratio
}

type messageL struct {
	// UID идентификатор транспортного средства в системе мониторинга
	UID string
}

type messageD struct {
	// Time дата и время сообщения
	Time time.Time
	/// Latitude широта
	Latitude Coordinate
	// Longitude долгота
	Longitude Coordinate
	// Speed Скорость
	Speed uint
	// Course курс
	Course uint8
	// Alt высота. Если отсутствует, значение null.
	// Alt *int
	// Sats Количество спутников. Если отсутствует, значение null.
	// Sats *int
}

type Point struct {
	messageL
	messageD
}
