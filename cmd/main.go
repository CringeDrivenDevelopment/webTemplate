package main

import (
	"backend/internal/infra"
	"backend/internal/infra/queries"
	"backend/internal/service"
	"backend/internal/transport/api/handlers"
	"backend/internal/transport/api/middlewares"

	"go.uber.org/fx"
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

	fx.New(
		fx.Provide(
			// REST API
			infra.NewEcho,
			middlewares.NewLogger,
			handlers.NewAuth,

			// services and infra
			infra.NewLogger,
			infra.NewConfig,
			infra.NewPostgresConnection,
			queries.NewUserRepo,
			service.NewAuth,
			service.NewUser,
		),
		/*
			fx.WithLogger(func(lc fx.Lifecycle, logger *infra.Logger) fxevent.Logger {
				return &infra.ZapFxLogger{Logger: logger.Zap}
			}),
		*/
		fx.Invoke(func(auth *handlers.Auth) {
			// need each of controllers, to register them

			// no need to call infra, apis and services, they're deps, started automatically
		}),
	).Run()
}
