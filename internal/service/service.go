package service

import (
	"context"

	"backend/internal/model"
)

// AuthService defines auth service interface

type AuthService interface {
	VerifyToken(authHeader string) (string, error)

	VerifyPassword(user model.User, password string) error

	GenerateToken(userID string) (string, error)
}

// UserService defines user service interface

type UserService interface {
	Create(ctx context.Context, email, password string) (string, error)

	GetByID(ctx context.Context, id string) (model.User, error)

	GetByEmail(ctx context.Context, email string) (model.User, error)
}
