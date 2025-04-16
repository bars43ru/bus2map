// Package model содержит основные структуры данных для работы с транспортом и маршрутами.
package model

// RouteNumber представляет номер маршрута общественного транспорта.
// Используется для идентификации маршрутов в системе.
type RouteNumber string

// StateNumber представляет государственный номер транспортного средства.
// Используется для идентификации транспортных средств.
type StateNumber string

// String возвращает строковое представление номера маршрута.
func (s RouteNumber) String() string {
	return string(s)
}

// String возвращает строковое представление государственного номера.
func (s StateNumber) String() string {
	return string(s)
}
