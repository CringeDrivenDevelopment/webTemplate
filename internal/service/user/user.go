package user

import (
	"context"

	"github.com/alexedwards/argon2id"
	"github.com/oklog/ulid/v2"

	"backend/internal/infra/queries"
)

func (s *Service) Create(ctx context.Context, email, password string) (string, error) {
	id := ulid.Make().String()

	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	if err = s.repository.Create(ctx, queries.User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
	}); err != nil {
		return id, err
	}

	return id, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (queries.User, error) {
	return s.repository.GetUserByID(ctx, id)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (queries.User, error) {
	return s.repository.GetUserByEmail(ctx, email)
}
