package v1

import (
	"backend/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"

	"backend/internal/infra"
	"backend/internal/service"
	"backend/internal/service/auth"
	"backend/internal/transport/api/dto"
)

type Auth struct {
	authService service.AuthService
	logger      *infra.Logger
}

// NewAuth - создать новый экземпляр обработчика
func NewAuth(authService *auth.Service, logger *infra.Logger, router *echo.Echo) *Auth {
	result := &Auth{
		authService: authService,
		logger:      logger,
	}

	router.POST("/api/login", result.login)

	return result
}

// login godoc
// @Summary      Login
// @Description  Вход в аккаунт
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

	token, err := h.authService.Login(ctx, data.Email, data.Password)
	if err != nil {
		return utils.Convert(err, h.logger)
	}

	tokenData := dto.Token{
		Token: token,
	}
	return echoCtx.JSON(http.StatusOK, tokenData)
}
