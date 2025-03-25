package egts

import (
	"time"
)

type Point struct {
	PacketID  uint32    // id пакета данных
	Client    uint32    // id устройства
	Time      time.Time // дата и время сообщения
	Latitude  float64   // широта
	Longitude float64   // долгота
	Speed     uint16    // скорость
	Course    uint8     // курс
}
