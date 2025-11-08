package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"booker/pkg/metrics"
)

func MetricsMiddleware(serviceName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			duration := time.Since(start).Seconds()
			status := strconv.Itoa(c.Response().Status)
			method := c.Request().Method
			path := c.Path()

			metrics.HTTPRequestsTotal.WithLabelValues(method, path, status, serviceName).Inc()
			metrics.HTTPRequestDuration.WithLabelValues(method, path, status, serviceName).Observe(duration)

			return err
		}
	}
}

