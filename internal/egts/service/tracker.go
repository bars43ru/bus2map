package service

import (
	"github.com/bars43ru/bus2map/internal/ctype"
	"github.com/imkira/go-observer/v2"
)

type Tracker observer.Property[ctype.RawGPSData]
