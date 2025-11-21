package utils

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/CringeDrivenDevelopment/webTemplate/internal/infra"
)

var (
	ErrInvalidToken = errors.New("invalid token")

	ErrInvalidPassword = errors.New("invalid password")

	ErrInvalidUser = errors.New("invalid user")

	ErrContextUserNotFound = errors.New("user not found in context")

	ErrEmailAlreadySignup = errors.New("email already signup")
)

func Convert(functionError error, logger *infra.Logger) error {
	if errors.Is(functionError, pgx.ErrNoRows) {
		return echo.ErrNotFound
	}

	if errors.Is(functionError, ErrInvalidToken) {
		return echo.ErrUnauthorized
	}

	if errors.Is(functionError, ErrEmailAlreadySignup) {
		return echo.ErrConflict
	}

	if errors.Is(functionError, ErrInvalidPassword) {
		return echo.ErrUnauthorized
	}

	if errors.Is(functionError, ErrInvalidUser) {
		return echo.ErrUnauthorized
	}

	logger.Error("500 error stacktrace", zap.Error(functionError))

	return echo.ErrInternalServerError
}
