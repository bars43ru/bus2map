package config

import (
	"log/slog"
)

// gps bus tracking
type Config struct {
	Logger    slog.Level `envPrefix:"LOG_"`
	GRPC      GRPCServer `envPrefix:"GRPC_"`
	WialonIPS TCPServer  `envPrefix:"WIALON_IPS_"`
	EGTS      TCPServer  `envPrefix:"EGTS_"`
	TwoGIS    Yandex     `envPrefix:"TWOGIS_"`
	Yandex    Yandex     `envPrefix:"YANDEX_"`
}

type Logger struct {
	Level slog.Level `env:"LEVEL,required"`
}

type TCPServer struct {
	Enabled bool   `env:"ENABLED,required"`
	Addr    string `env:"LISTEN_ADDR,required"`
}

type Yandex struct {
	Enabled bool   `env:"ENABLED,required"`
	Clid    string `env:"CLID,required"`
	Url     string `env:"URL,required"`
}

type GRPCServer struct {
	ListenAddr    string `env:"LISTEN_ADDR,required"`
	UseReflection bool   `env:"REFLECTION,required"`
}
