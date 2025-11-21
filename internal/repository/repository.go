package repository

import (
	"context"

	"github.com/CringeDrivenDevelopment/webTemplate/internal/infra/queries"
)

type UserRepository interface {
	Create(ctx context.Context, user queries.User) error
	GetUserByID(ctx context.Context, id string) (queries.User, error)
	GetUserByEmail(ctx context.Context, email string) (queries.User, error)
}
