package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

func New() (Config, error) {
	config, err := env.ParseAs[Config]()
	if err != nil {
		return config, fmt.Errorf("parse env: %w", err)
	}
	return config, err
}
