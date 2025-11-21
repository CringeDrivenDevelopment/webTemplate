package userRepo

import (
	"context"
	"strings"

	"backend/internal/infra/queries"
	"backend/pkg/utils"
)

func (ur *UserRepository) Create(ctx context.Context, user queries.User) error {
	rq := queries.New(ur.pgxpool)
	if _, err := rq.GetUserByEmail(ctx, strings.ToLower(user.Email)); err == nil {
		return utils.ErrEmailAlreadySignup
	}
	if err := utils.ExecInTx(ctx, ur.pgxpool, func(tq *queries.Queries) error {
		return tq.CreateUser(ctx, queries.CreateUserParams(user))
	}); err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (queries.User, error) {
	rq := queries.New(ur.pgxpool)
	return rq.GetUserByEmail(ctx, strings.ToLower(email))
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id string) (queries.User, error) {
	rq := queries.New(ur.pgxpool)
	return rq.GetUserByID(ctx, id)
}
