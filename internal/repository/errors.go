// Package repository содержит реализацию репозиториев для работы с данными.
// Предоставляет интерфейсы и структуры для хранения и доступа к информации о маршрутах, транспорте и расписании.
package repository

import "errors"

// ErrNotFound возвращается при отсутствии запрашиваемых данных в репозитории.
// Используется для обработки ситуаций, когда данные не найдены.
var ErrNotFound = errors.New("not found")

const (
	// FileDatasourceRoute путь к файлу с данными о маршрутах.
	// Содержит информацию о номерах маршрутов и их идентификаторах в различных системах.
	FileDatasourceRoute = "./datasource/route.txt"

	// FileDatasourceSchedule путь к файлу с данными о расписании.
	// Содержит информацию о времени работы транспорта на маршрутах.
	FileDatasourceSchedule = "./datasource/schedule.txt"

	// FileDatasourceTransport путь к файлу с данными о транспорте.
	// Содержит информацию о транспортных средствах и их идентификаторах.
	FileDatasourceTransport = "./datasource/transport.txt"
)
