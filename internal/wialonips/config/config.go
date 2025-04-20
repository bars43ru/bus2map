package config

import (
	"log/slog"
)

type Config struct {
	Logger         Logger    `envPrefix:"LOG_"`
	WialonIPS      TCPServer `envPrefix:"WIALON_IPS_"`
	CoordinateAddr string    `env:"COORDINATE_ADDR,required"`
	Source         string    `env:"SOURCE,required"`
}

type Logger struct {
	Level slog.Level `env:"LEVEL,required"`
}

type TCPServer struct {
	Addr string `env:"LISTEN_ADDR,required"`
}
