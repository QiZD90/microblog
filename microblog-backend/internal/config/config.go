package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Port           string `env:"BACKEND_HTTP_PORT" envDefault:"8080"`
	PostgresPort   string `env:"POSTGRES_PORT"`
	PrometheusPort string `env:"PROMETHEUS_PORT"`
}

func New() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
