package user

import (
	"context"
	"errors"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"

	"github.com/CringeDrivenDevelopment/webTemplate/internal/infra/queries"
	"github.com/CringeDrivenDevelopment/webTemplate/pkg/utils"
)

func (s *Service) Register(ctx context.Context, email, password string) (string, error) {
	if _, err := s.repository.GetUserByEmail(ctx, email); err == nil {
		return "", utils.ErrEmailAlreadySignup
	} else if !errors.Is(err, pgx.ErrNoRows) {
		// Some other error occurred
		return "", err
	}

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
		return "", err
	}

	return id, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (queries.User, error) {
	return s.repository.GetUserByID(ctx, id)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (queries.User, error) {
	return s.repository.GetUserByEmail(ctx, email)
}
