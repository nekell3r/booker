package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"booker/pkg/metrics"
)

func MetricsMiddleware(serviceName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Log incoming request
			log.Info().
				Str("service", serviceName).
				Str("method", c.Request().Method).
				Str("path", c.Request().URL.Path).
				Str("remote_addr", c.Request().RemoteAddr).
				Msg("HTTP request received")

			err := next(c)

			duration := time.Since(start).Seconds()
			status := strconv.Itoa(c.Response().Status)
			method := c.Request().Method
			path := c.Path()

			// Log response
			if err != nil {
				log.Error().
					Err(err).
					Str("service", serviceName).
					Str("method", method).
					Str("path", path).
					Int("status", c.Response().Status).
					Float64("duration_seconds", duration).
					Msg("HTTP request failed")
			} else {
				log.Info().
					Str("service", serviceName).
					Str("method", method).
					Str("path", path).
					Int("status", c.Response().Status).
					Float64("duration_seconds", duration).
					Msg("HTTP request completed")
			}

			metrics.HTTPRequestsTotal.WithLabelValues(method, path, status, serviceName).Inc()
			metrics.HTTPRequestDuration.WithLabelValues(method, path, status, serviceName).Observe(duration)

			return err
		}
	}
}

