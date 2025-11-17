package interfaces

import (
	"backend/internal/infra/queries"
	"context"
)

type UserService interface {
	Create(ctx context.Context, email, password string) (string, error)
	GetByID(ctx context.Context, id string) (queries.User, error)
	GetByEmail(ctx context.Context, email string) (queries.User, error)
}
type UserRepository interface {
	CreateUser(ctx context.Context, user queries.User) error
	GetUserByID(ctx context.Context, id string) (queries.User, error)
	GetUserByEmail(ctx context.Context, email string) (queries.User, error)
}
type AuthService interface {
	VerifyPassword(user queries.User, password string) error
	VerifyToken(authHeader string) (string, error)
	GenerateToken(userID string) (string, error)
}
