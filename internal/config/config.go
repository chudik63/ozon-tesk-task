package config

import (
	"errors"
	"fmt"
	"ozon-tesk-task/internal/database/postgres"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	postgres.Config

	MigrationsPath string `env:"MIGRATIONS_PATH"`
	ServicePort    string `env:"SERVICE_PORT"`
}

func New() (*Config, error) {
	cfg := Config{}

	err := cleanenv.ReadEnv(&cfg)

	if cfg == (Config{}) {
		return nil, errors.New("config is empty")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}
