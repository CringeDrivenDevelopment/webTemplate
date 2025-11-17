package main

import (
	"backend/internal/infra"
	"backend/internal/service"
	"backend/internal/transport/api/handlers"
	"backend/internal/transport/api/middlewares"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// @title           Backend API
// @version         1.0

// @host      localhost:8080
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description "Type 'Bearer TOKEN' to correctly set the API Key"
func main() {
	// TODO: log db requests
	// TODO: add otel
	// TODO: add image proxy, DL
	// TODO: process errors like .Error(), if code is 500 - print stacktrace

	fx.New(
		fx.Provide(
			// REST API
			infra.NewEcho,
			middlewares.NewLogger,
			middlewares.NewAuth,
			handlers.NewAuth,

			// services and infra
			infra.NewLogger,
			infra.NewConfig,
			infra.NewPostgresConnection,
			service.NewAuth,
			service.NewUser,
		),
		fx.WithLogger(func(lc fx.Lifecycle, logger *zap.Logger) fxevent.Logger {
			return &infra.ZapFxLogger{Logger: logger}
		}),
		fx.Invoke(func(auth *handlers.Auth) {
			// need each of controllers, to register them

			// no need to call infra, apis and services, they're deps, started automatically
		}),
	).Run()
}
