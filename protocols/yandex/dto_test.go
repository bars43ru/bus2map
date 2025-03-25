package yandex

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPointMarshalXml(t *testing.T) {
	const wantXml = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
		"<Point latitude=\"55.75363\" longitude=\"55.75363\" avg_speed=\"0\" direction=\"242\" time=\"10012009:172045\"></Point>"
	point := Point{
		Latitude:  55.753630,
		Longitude: 55.753630,
		AvgSpeed:  0,
		Direction: 242,
		Time:      CustomTime(time.Date(2009, 01, 10, 17, 20, 45, 0, time.UTC)),
	}
	xmlValue, err := xml.Marshal(point)
	assert.NoError(t, err)
	assert.Equal(t, wantXml, xml.Header+string(xmlValue))
}

func TestTrackMarshalXml(t *testing.T) {
	const wantXml = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
		"<Track uuid=\"123456789\" category=\"n\" route=\"145\" vehicle_type=\"bus\">" +
		"<point latitude=\"55.75363\" longitude=\"55.75363\" avg_speed=\"0\" direction=\"242\" time=\"10012009:142045\"></point>" +
		"</Track>"
	track := Track{
		UUID:        "123456789",
		Category:    NormalGpsSignal,
		Route:       "145",
		VehicleType: BusVehicleType,
		Point: Point{
			Latitude:  55.753630,
			Longitude: 55.753630,
			AvgSpeed:  0,
			Direction: 242,
			Time:      CustomTime(time.Date(2009, 01, 10, 14, 20, 45, 0, time.UTC)),
		},
	}
	xmlValue, err := xml.Marshal(track)
	assert.NoError(t, err)
	assert.Equal(t, wantXml, xml.Header+string(xmlValue))
}

func TestTracksMarshalXml(t *testing.T) {
	const wantXml = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<tracks clid=\"123\"></tracks>"
	tracks := tracks{
		Clid:   "123",
		Tracks: nil,
	}
	xmlValue, err := xml.Marshal(tracks)
	assert.NoError(t, err)
	assert.Equal(t, wantXml, xml.Header+string(xmlValue))
}

func TestFullMarshalXml(t *testing.T) {
	const wantXml = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
		"<tracks clid=\"123\">" +
		"<track uuid=\"123456789\" category=\"s\" route=\"145\" vehicle_type=\"tramway\">" +
		"<point latitude=\"55.75363\" longitude=\"55.75363\" avg_speed=\"0\" direction=\"242\" time=\"10012009:142045\"></point>" +
		"</track>" +
		"</tracks>"
	tracks := tracks{
		Clid: "123",
		Tracks: []Track{
			{
				UUID:        "123456789",
				Category:    SlowGpsSignal,
				Route:       "145",
				VehicleType: TramwayVehicleType,
				Point: Point{

					Latitude:  55.753630,
					Longitude: 55.753630,
					AvgSpeed:  0,
					Direction: 242,
					Time:      CustomTime(time.Date(2009, 01, 10, 14, 20, 45, 0, time.UTC)),
				},
			},
		},
	}
	xmlValue, err := xml.Marshal(tracks)
	assert.NoError(t, err)
	assert.Equal(t, wantXml, xml.Header+string(xmlValue))
}
