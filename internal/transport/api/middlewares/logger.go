package middlewares

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"backend/internal/infra"
)

func NewLogger(log *infra.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()

			res := c.Response()

			fields := map[string]interface{}{
				"remote_ip": c.RealIP(),

				"latency": time.Since(start).String(),

				"host": req.Host,

				"request": fmt.Sprintf("%s %s", req.Method, req.RequestURI),

				"status": res.Status,

				"size": res.Size,

				"user_agent": req.UserAgent(),
			}

			id := req.Header.Get(echo.HeaderXRequestID)

			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			fields["request_id"] = id

			n := res.Status

			switch {

			case n >= 500:

				log.With(zap.Error(err)).Error("Server error", fields)

			case n >= 400:

				log.With(zap.Error(err)).Warn("Client error", fields)

			case n >= 300:

				log.Info("Redirection", fields)

			default:

				log.Info("Success", fields)

			}

			return nil
		}
	}
}
