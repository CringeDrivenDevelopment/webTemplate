package infra

import (
	"errors"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	// DbUrl - Postgres Database connection string
	// Example - "postgres://username:password@localhost:5432/database_name"
	DbUrl string

	DbHost     string `env:"POSTGRES_HOST" env-default:"localhost"`
	DbPort     string `env:"POSTGRES_PORT" env-default:"5432"`
	DbPassword string `env:"POSTGRES_PASSWORD" env-default:"password"`
	DbUser     string `env:"POSTGRES_USER" env-default:"postgres"`
	DbName     string `env:"POSTGRES_DB" env-default:"backend"`

	JwtSecret string `env:"JWT_SECRET"`

	Debug bool `env:"DEBUG" env-default:"false"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	cfg.DbUrl = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)

	if cfg.JwtSecret == "" {
		return nil, errors.New("JWT_SECRET is REQUIRED not to be null")
	}

	return &cfg, nil
}
