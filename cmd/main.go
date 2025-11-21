package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/CringeDrivenDevelopment/webTemplate/internal/infra"
	"github.com/CringeDrivenDevelopment/webTemplate/internal/repository"
	userRepo "github.com/CringeDrivenDevelopment/webTemplate/internal/repository/user"
	"github.com/CringeDrivenDevelopment/webTemplate/internal/service/auth"
	"github.com/CringeDrivenDevelopment/webTemplate/internal/service/user"
	authV1 "github.com/CringeDrivenDevelopment/webTemplate/internal/transport/api/handlers/auth/v1"
	userV1 "github.com/CringeDrivenDevelopment/webTemplate/internal/transport/api/handlers/user/v1"
	"github.com/CringeDrivenDevelopment/webTemplate/internal/transport/api/middlewares"
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
	// TODO: add tracing, logging and metrics

	cfg, err := infra.NewConfig()
	if err != nil {
		panic(err)
	}

	logger, err := infra.NewLogger(cfg)
	if err != nil {
		panic(err)
	}

	fx.New(
		fx.Supply(logger.Zap, logger, cfg),
		fx.Provide(
			// REST API
			infra.NewEcho,
			middlewares.NewLogger,
			authV1.NewAuth,
			userV1.NewUser,

			// services and infra
			infra.NewPostgresConnection,
			fx.Annotate(
				userRepo.New,
				fx.As(new(repository.UserRepository)),
			),
			user.NewService,
			auth.NewService,
		),

		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			defer func(Zap *zap.Logger) {
				err := Zap.Sync()
				if err != nil {
					println(err)
				}
			}(logger.Zap)

			return &fxevent.ZapLogger{Logger: logger.Zap}
		}),

		// need each of controllers, to register them
		// no need to call infra, apis and services, they're deps, started automatically
		fx.Invoke(func(auth *authV1.Auth) {}, func(user *userV1.User) {}),
	).Run()
}
