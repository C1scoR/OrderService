package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string `env:"ENV" env-default:"development"`

	Port    int           `env:"PORT" env-default:"8080"`
	Timeout time.Duration `env:"HTTP_TIMEOUT" env-default:"30s"`

	DebeziumBaseURL string `env:"DEBEZIUM_BASE_URL" env-default:"http://localhost:8080"`
}

func ParseConfig(path string) (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}
