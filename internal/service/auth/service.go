package auth

import (
	"time"

	"github.com/CringeDrivenDevelopment/webTemplate/internal/infra"
	"github.com/CringeDrivenDevelopment/webTemplate/internal/repository"
)

type Service struct {
	secret     string
	expires    time.Duration
	repository repository.UserRepository
}

// NewService - создать новый экземпляр сервиса авторизации
func NewService(cfg *infra.Config, userRepository repository.UserRepository) *Service {
	return &Service{
		secret:     cfg.JwtSecret,
		expires:    time.Hour,
		repository: userRepository,
	}
}
