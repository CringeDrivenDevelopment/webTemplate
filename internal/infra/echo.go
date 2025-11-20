package infra

import (
	"context"
	"errors"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type sonicJSONSerializer struct{}

func (s *sonicJSONSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	enc := sonic.ConfigStd.NewEncoder(c.Response())

	if indent != "" {
		enc.SetIndent("", indent)
	}

	return enc.Encode(i)
}

func (s *sonicJSONSerializer) Deserialize(c echo.Context, i interface{}) error {
	return sonic.ConfigStd.NewDecoder(c.Request().Body).Decode(i)
}

func NewEcho(lc fx.Lifecycle, cfg *Config, logger *Logger, loggerWare echo.MiddlewareFunc) *echo.Echo {
	router := echo.New()

	if !cfg.Debug {
		router.Use(middleware.Recover())
	}

	router.JSONSerializer = &sonicJSONSerializer{}

	router.GET("/api/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	router.HideBanner = true

	router.HidePort = true

	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"*",
		},
	}))

	router.Use(loggerWare)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("starting server on :8080")

			go func() {
				err := router.Start(":8080")

				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Fatal("stopping server, cause: error", zap.Error(err))
				}
			}()

			return nil
		},

		OnStop: func(ctx context.Context) error {
			logger.Info("stopped server")

			return router.Shutdown(ctx)
		},
	})

	return router
}
