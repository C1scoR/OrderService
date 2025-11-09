package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment string `env:"ENV" env-default:"development"`

	Port       int            `env:"PORT" env-default:"50051"`
	Timeout    time.Duration  `env:"HTTP_TIMEOUT" env-default:"30s"`
	PostgreSQL PostgresConfig `env-prefix:"POSTGRES_"`
}

type PostgresConfig struct {
	Host     string `env:"HOST" env-default:"localhost"`
	Port     int    `env:"PORT" env-default:"50000"`
	User     string `env:"USER" env-default:"postgres"`
	Password string `env:"PASSWORD" env-default:"postgres"`
	Database string `env:"DATABASE" env-default:"postgres"`
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

// ConnString() метод для конфигурации PostgresConfig, который будет возвращать строку подключения
func (pg *PostgresConfig) ConnString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		pg.Host,
		pg.Port,
		pg.User,
		pg.Password,
		pg.Database,
	)
}
