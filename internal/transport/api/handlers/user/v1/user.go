package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/CringeDrivenDevelopment/webTemplate/internal/infra"
	"github.com/CringeDrivenDevelopment/webTemplate/internal/service"
	"github.com/CringeDrivenDevelopment/webTemplate/internal/service/user"
	"github.com/CringeDrivenDevelopment/webTemplate/internal/transport/api/dto"
	"github.com/CringeDrivenDevelopment/webTemplate/pkg/utils"
)

type User struct {
	userService service.UserService
	logger      *infra.Logger
}

// NewUser - создать новый экземпляр обработчика
func NewUser(userService *user.Service, logger *infra.Logger, router *echo.Echo) *User {
	result := &User{
		userService: userService,
		logger:      logger,
	}

	router.POST("/api/register", result.register)
	return result
}

// register godoc
// @Summary      Register
// @Description  Регистрация
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.AuthData  true  "Auth data"
// @Success      200  {object}  dto.Token
// @Failure      400  {object}  dto.ApiError
// @Failure      401  {object}  dto.ApiError
// @Failure      500  {object}  dto.ApiError
// @Router       /api/user/v1/register [post]
func (h *User) register(echoCtx echo.Context) error {
	var data dto.AuthData
	if err := echoCtx.Bind(&data); err != nil {
		return err
	}

	ctx := echoCtx.Request().Context()
	h.logger.Info("register: " + data.Email)

	token, err := h.userService.Register(ctx, data.Email, data.Password)
	if err != nil {
		return utils.Convert(err, h.logger)
	}

	tokenData := dto.Token{
		Token: token,
	}
	return echoCtx.JSON(http.StatusOK, tokenData)
}
