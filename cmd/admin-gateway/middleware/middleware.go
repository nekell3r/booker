package middleware

import (
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"

	"booker/cmd/admin-gateway/config"
	"booker/pkg/redis"
)

type Middleware struct {
	redisClient *redis.Client
	cfg         *config.Config
}

func New(redisClient *redis.Client, cfg *config.Config) *Middleware {
	return &Middleware{
		redisClient: redisClient,
		cfg:         cfg,
	}
}

func (m *Middleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		log.Info().
			Str("path", c.Path()).
			Str("method", c.Request().Method).
			Str("auth_header", authHeader).
			Msg("AuthMiddleware: checking request")
		
		if authHeader == "" {
			log.Warn().
				Str("path", c.Path()).
				Str("method", c.Request().Method).
				Msg("AuthMiddleware: missing authorization header")
			return c.JSON(401, map[string]string{"error": "missing authorization header"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Warn().
				Str("path", c.Path()).
				Str("method", c.Request().Method).
				Str("auth_header", authHeader).
				Msg("AuthMiddleware: invalid authorization header format")
			return c.JSON(401, map[string]string{"error": "invalid authorization header"})
		}

		token := parts[1]

		// TODO: Validate JWT token
		// For now, just extract admin_id from token or use dummy
		adminID := "admin-1" // This should come from JWT claims

		log.Info().
			Str("path", c.Path()).
			Str("method", c.Request().Method).
			Str("admin_id", adminID).
			Msg("AuthMiddleware: request authorized")

		c.Set("admin_id", adminID)
		c.Set("token", token)

		return next(c)
	}
}

func (m *Middleware) RateLimitMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			adminID := c.Get("admin_id")
			if adminID == nil {
				return next(c)
			}

			key := "rl:" + adminID.(string)
			limit := 100 // requests per minute

			count, err := m.redisClient.Incr(c.Request().Context(), key)
			if err != nil {
				log.Error().Err(err).Msg("Rate limit check failed")
				return next(c) // Fail open
			}

			if count == 1 {
				m.redisClient.Expire(c.Request().Context(), key, time.Minute)
			}

			if count > int64(limit) {
				return c.JSON(429, map[string]string{"error": "rate limit exceeded"})
			}

			return next(c)
		}
	}
}

func (m *Middleware) SetupMiddleware(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
}


