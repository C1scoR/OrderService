package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment string `env:"ENV" env-default:"development"`

	Port    int           `env:"PORT" env-default:"50051"`
	Timeout time.Duration `env:"HTTP_TIMEOUT" env-default:"30s"`
}

// ParseConfig загружает переменные в окружение и возвращает экземпляр структуры Config
func ParseConfig() (*Config, error) {
	//Сначала выгружаем все переменные из .env файла в окружение
	_ = godotenv.Load(".env")
	cfg := &Config{}
	//читаем эти переменные из окружения, и ими же заполняем экземпляр конфига
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
