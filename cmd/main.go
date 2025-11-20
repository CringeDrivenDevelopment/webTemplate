package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"backend/internal/infra"
	"backend/internal/repository"
	userRepo "backend/internal/repository/user"
	"backend/internal/service/auth"
	"backend/internal/service/user"
	"backend/internal/transport/api/handlers"
	"backend/internal/transport/api/middlewares"
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
			fx.Annotate(
				userRepo.NewRepository,
				fx.As(new(repository.UserRepository)),
			),
			user.NewService,
			auth.NewService,
		),
		fx.WithLogger(func(lc fx.Lifecycle, logger *infra.Logger) fxevent.Logger {
			return &infra.ZapFxLogger{Logger: logger.Zap}
		}),
		fx.Invoke(func(auth *handlers.Auth) {
			// need each of controllers, to register them

			// no need to call infra, apis and services, they're deps, started automatically
		}),
	).Run()
}
