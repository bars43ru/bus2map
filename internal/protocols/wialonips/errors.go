// Package wialonips содержит реализацию парсера протокола WialonIPS.
// Предоставляет функции для разбора и валидации данных телеметрии в формате WialonIPS.
package wialonips

import "errors"

// ErrFormat возвращается при обнаружении некорректного формата данных.
// Используется для обработки ошибок парсинга и валидации пакетов WialonIPS.
var ErrFormat = errors.New("incorrect format")
