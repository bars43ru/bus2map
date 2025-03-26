package receiver

import (
	"context"

	"github.com/bars43ru/bus2map/internal/model"
)

type GPSLocator interface {
	ProcessGPSData(ctx context.Context, gps model.GPS)
}
