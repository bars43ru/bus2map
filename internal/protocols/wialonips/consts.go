// Package wialonips содержит реализацию парсера протокола WialonIPS.
// Предоставляет функции для разбора и валидации данных телеметрии в формате WialonIPS.
package wialonips

const (
	// delimiter определяет символ-разделитель между пакетами данных
	delimiter byte = '\n'

	// Имена полей в регулярных выражениях
	uidField    = "uid"    // Идентификатор устройства
	dateField   = "date"   // Дата в формате DDMMYY
	timeField   = "time"   // Время в формате HHMMSS
	lat1Field   = "lat1"   // Первая часть широты (градусы)
	lat2Field   = "lat2"   // Вторая часть широты (направление)
	lon1Field   = "lon1"   // Первая часть долготы (градусы)
	lon2Field   = "lon2"   // Вторая часть долготы (направление)
	speedField  = "speed"  // Скорость в км/ч
	courseField = "course" // Курс в градусах
	altField    = "alt"    // Высота над уровнем моря
	satsField   = "sats"   // Количество спутников

	// patternL - шаблон для разбора пакета с идентификатором устройства
	patternL = `#L#(?P<` + uidField + `>\w+);`

	// patternD - шаблон для разбора пакета с данными телеметрии
	patternD = `#D#` +
		`(?P<` + dateField + `>\d+);` +
		`(?P<` + timeField + `>\d+);` +
		`(?P<` + lat1Field + `>\d+\.\d+);` +
		`(?P<` + lat2Field + `>\w+);` +
		`(?P<` + lon1Field + `>\d+\.\d+);` +
		`(?P<` + lon2Field + `>\w+);` +
		`(?P<` + speedField + `>\d+);` +
		`(?P<` + courseField + `>\d+);` +
		`(?P<` + altField + `>\d+\.\d+);` +
		`(?P<` + satsField + `>\d+)`

	// layoutTime определяет формат времени в пакетах WialonIPS (DDMMYYHHMMSS)
	layoutTime = "020106150405"
)
