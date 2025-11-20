package user

import (
	"backend/internal/repository"
)

type Service struct {
	repository repository.UserRepository
}

func NewService(repository repository.UserRepository) *Service {
	return &Service{repository: repository}
}
