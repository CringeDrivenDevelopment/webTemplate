package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"backend/internal/infra"
	"backend/internal/service"
	"backend/internal/service/auth"
	"backend/internal/service/user"
	"backend/internal/transport/api/dto"
	"backend/pkg/utils"
)

type Auth struct {
	userService service.UserService
	authService service.AuthService

	logger *infra.Logger
}

// NewAuth - создать новый экземпляр обработчика
func NewAuth(userService *user.Service, authService *auth.Service, logger *infra.Logger, router *echo.Echo) *Auth {
	result := &Auth{
		userService: userService,
		authService: authService,
		logger:      logger,
	}

	router.POST("/api/login", result.login)

	return result
}

// login godoc
// @Summary      Login
// @Description  Войти в аккаунт, также выполняет функцию регистрации
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.AuthData  true  "Auth data"
// @Success      200  {object}  dto.Token
// @Failure      400  {object}  dto.ApiError
// @Failure      401  {object}  dto.ApiError
// @Failure      500  {object}  dto.ApiError
// @Router       /api/login [post]
func (h *Auth) login(echoCtx echo.Context) error {
	var data dto.AuthData
	if err := echoCtx.Bind(&data); err != nil {
		return err
	}

	ctx := echoCtx.Request().Context()

	h.logger.Info("login: " + data.Email)

	var err error
	var userID string

	user, err := h.userService.GetByEmail(ctx, data.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Warn(fmt.Sprintf("login error: email - %s, error - %s", data.Email, err.Error()))

			return utils.Convert(err, h.logger)
		}

		userID, err = h.userService.Create(ctx, data.Email, data.Password)
		if err != nil {
			h.logger.Warn(fmt.Sprintf("login error: email - %s, error - %s", data.Email, err.Error()))

			return utils.Convert(err, h.logger)
		}
	} else {
		if err = h.authService.VerifyPassword(user, data.Password); err != nil {
			h.logger.Warn(fmt.Sprintf("login error: email - %s, error - %s", data.Email, err.Error()))

			return utils.Convert(err, h.logger)
		}

		userID = user.ID
	}

	token, err := h.authService.GenerateToken(userID)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("login error: email - %s, error - %s", data.Email, err.Error()))

		return utils.Convert(err, h.logger)
	}

	tokenData := dto.Token{
		Token: token,
	}

	return echoCtx.JSON(http.StatusOK, tokenData)
}
