package service

import (
	"backend/internal/infra/queries"
	"backend/pkg/utils"
	"context"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

type User struct {
	pool *pgxpool.Pool
}

func NewUser(pool *pgxpool.Pool) *User {
	return &User{pool: pool}
}

func (s *User) Create(ctx context.Context, email, password string) (string, error) {
	id := ulid.Make().String()

	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	rq := queries.New(s.pool)
	if _, err := rq.GetUserById(ctx, id); err == nil {
		return id, nil
	}

	if err := utils.ExecInTx(ctx, s.pool, func(tq *queries.Queries) error {
		return tq.CreateUser(ctx, queries.CreateUserParams{
			ID:           id,
			Email:        email,
			PasswordHash: passwordHash,
		})
	}); err != nil {
		return id, err
	}

	return id, nil
}

func (s *User) GetByID(ctx context.Context, id string) (queries.User, error) {
	rq := queries.New(s.pool)

	return rq.GetUserById(ctx, id)
}

func (s *User) GetByEmail(ctx context.Context, email string) (queries.User, error) {
	rq := queries.New(s.pool)

	return rq.GetUserByEmail(ctx, email)
}
