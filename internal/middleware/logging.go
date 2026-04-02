package middleware

import (
	"time"
	"url-shortener/internal/infrastructure/logger"

	"github.com/labstack/echo/v4"
)

func LoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			ctx := logger.WithTraceID(c.Request().Context(), requestID)
			c.SetRequest(c.Request().WithContext(ctx))
			log := logger.FromContext(ctx)

			log.Info("Incoming request",
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"client_ip", c.RealIP(),
				"user_agent", c.Request().UserAgent(),
			)

			err := next(c)

			latency := time.Since(start)
			attrs := []any{
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", c.Response().Status,
				"latency_ms", latency.Milliseconds(),
			}

			if err != nil {
				attrs = append(attrs, "error", err)
				log.Error("Request failed", attrs...)
			} else if c.Response().Status >= 500 {
				log.Error("Request completed with server error", attrs...)
			} else if c.Response().Status >= 400 {
				log.Warn("Request completed with client error", attrs...)
			} else {
				log.Info("Request completed", attrs...)
			}

			return err
		}
	}
}
