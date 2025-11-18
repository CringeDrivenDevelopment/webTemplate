package utils

import (
	"backend/internal/infra"
	"errors"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrContextUserNotFound = errors.New("user not found in context")
)

func Convert(functionError error, logger *infra.Logger) error {
	if errors.Is(functionError, gorm.ErrRecordNotFound) {
		return echo.ErrNotFound
	}

	if errors.Is(functionError, ErrInvalidToken) {
		return echo.ErrUnauthorized
	}

	if errors.Is(functionError, ErrInvalidPassword) {
		return echo.ErrUnauthorized
	}

	logger.Error("500 error stacktrace", zap.Error(functionError))

	return echo.ErrInternalServerError
}
