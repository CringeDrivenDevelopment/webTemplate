package auth

import (
	"time"

	"backend/internal/infra"
)

type Service struct {
	secret  string
	expires time.Duration
}

// NewService - создать новый экземпляр сервиса авторизации
func NewService(cfg *infra.Config) *Service {
	return &Service{
		secret:  cfg.JwtSecret,
		expires: time.Hour,
	}
}
