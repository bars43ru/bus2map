// Package transport_type определяет типы транспортных средств, используемых в системе.
// Поддерживает различные виды общественного транспорта: автобусы, троллейбусы, трамваи и маршрутки.
package transport_type

//go:generate go tool go-enum --lower --marshal --names --values transport_type

// Type представляет тип транспортного средства.
// Используется для классификации транспорта в системе.
// ENUM(
//
//	BUS,        // Автобус
//	TROLLEYBUS, // Троллейбус
//	TRAMWAY,    // Трамвай
//	MINIBUS,    // Маршрутное такси
//
// )
type Type int
