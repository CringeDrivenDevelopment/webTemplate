package service

import (
	"backend/internal/infra/queries"
	"backend/internal/interfaces"
	"context"

	"github.com/alexedwards/argon2id"
	"github.com/oklog/ulid/v2"
)

type User struct {
	userRepo interfaces.UserRepository
}

func NewUser(userRepo interfaces.UserRepository) *User {
	return &User{userRepo: userRepo}
}

func (s *User) Create(ctx context.Context, email, password string) (string, error) {
	id := ulid.Make().String()

	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	if err := s.userRepo.CreateUser(ctx, queries.User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
	}); err != nil {
		return id, err
	}

	return id, nil
}

func (s *User) GetByID(ctx context.Context, id string) (queries.User, error) {

	return s.userRepo.GetUserByID(ctx, id)
}

func (s *User) GetByEmail(ctx context.Context, email string) (queries.User, error) {

	return s.userRepo.GetUserByEmail(ctx, email)
}
