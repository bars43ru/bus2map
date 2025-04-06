package wialonips

const (
	delimiter byte = '\n'

	uidField    = "uid"
	dateField   = "date"
	timeField   = "time"
	lat1Field   = "lat1"
	lat2Field   = "lat2"
	lon1Field   = "lon1"
	lon2Field   = "lon2"
	speedField  = "speed"
	courseField = "course"
	altField    = "alt"
	satsField   = "sats"

	patternL = `#L#(?P<` + uidField + `>\w+);`
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

	layoutTime = "020106150405"
)
