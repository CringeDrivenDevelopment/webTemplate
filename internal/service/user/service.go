package user

import (
	"github.com/CringeDrivenDevelopment/webTemplate/internal/repository"
)

type Service struct {
	repository repository.UserRepository
}

func NewService(repository repository.UserRepository) *Service {
	return &Service{repository: repository}
}
