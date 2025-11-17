package handlers

import (
	"backend/internal/interfaces"
	"backend/internal/service"
	"backend/internal/transport/api/dto"
	"backend/pkg/utils"
	"context"
	"errors"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Auth struct {
	userService interfaces.UserService
	authService interfaces.AuthService

	logger *zap.Logger
}

// NewAuth - создать новый экземпляр обработчика
func NewAuth(userService *service.User, authService *service.Auth, logger *zap.Logger, api huma.API) *Auth {
	result := &Auth{
		userService: userService,
		authService: authService,
		logger:      logger,
	}

	result.setup(api)

	return result
}

// login - Получить токен для взаимодействия. Нуждается в Raw строке из Telegram Mini App. Действует 1 час
func (h *Auth) login(ctx context.Context, input *dto.AuthInputStruct) (*dto.AuthOutputStruct, error) {
	data := input.Body

	h.logger.Info("login: " + data.Email)

	var err error
	var userID string

	user, err := h.userService.GetByEmail(ctx, data.Email)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			h.logger.Warn(fmt.Sprintf("login error: email - %s, error - %s", input.Body.Email, err.Error()))

			return nil, utils.Convert(err, h.logger)
		}

		userID, err = h.userService.Create(ctx, data.Email, data.Password)
		if err != nil {
			h.logger.Warn(fmt.Sprintf("login error: email - %s, error - %s", data.Email, err.Error()))

			return nil, utils.Convert(err, h.logger)
		}
	} else {
		if err := h.authService.VerifyPassword(user, data.Password); err != nil {
			h.logger.Warn(fmt.Sprintf("login error: email - %s, error - %s", data.Email, err.Error()))

			return nil, utils.Convert(err, h.logger)
		}

		userID = user.ID
	}

	token, err := h.authService.GenerateToken(userID)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("login error: email - %s, error - %s", data.Email, err.Error()))

		return nil, utils.Convert(err, h.logger)
	}

	tokenData := dto.Token{
		Token: token,
	}

	return &dto.AuthOutputStruct{Body: tokenData}, nil
}
