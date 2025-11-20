package repository

import (
	"context"

	"backend/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) error
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	GetUserByID(ctx context.Context, id string) (model.User, error)
}
