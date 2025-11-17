package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// setup - добавить маршрут до эндпоинта
func (h *Auth) setup(router huma.API) {
	huma.Register(router, huma.Operation{
		OperationID: "auth",
		Path:        "/api/auth",
		Method:      http.MethodPost,
		Errors: []int{
			401,
			422,
			500,
		},
		Tags: []string{
			"auth",
		},
		Summary:     "Login",
		Description: "Получить токен для взаимодействия. Нуждается в логине и пароле, сразу выполняет и регистрацию, и вход. Действует 1 час",
	}, h.login)
}
