package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string `env:"ENV" env-default:"development"`

	Port    int           `env:"PORT" env-default:"50051"`
	Timeout time.Duration `env:"HTTP_TIMEOUT" env-default:"30s"`
}

func ParseConfig(path string) (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}
