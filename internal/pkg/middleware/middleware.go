package middleware

import (
	"fmt"

	"fiber-admin/internal/pkg/config"
	ware "fiber-admin/internal/pkg/middleware/mods"
	"fiber-admin/pkg/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Middleware struct {
	AuthMiddleware        *ware.AuthMiddleware
	LoggingMiddleware     *ware.LoggingMiddleware
	PrometheusMiddleware  *ware.PrometheusMiddleware
	ContextMiddleware     *ware.ContextMiddleware
	IdempotencyMiddleware *ware.IdempotencyMiddleware
	Config                *config.Config
}

func (m *Middleware) Register(app *fiber.App) error {
	// Register RequestID Middleware
	app.Use(
		requestid.New(
			requestid.Config{
				Header:     fiber.HeaderXRequestID,
				ContextKey: config.RequestIDKey,
			},
		),
	)

	// Register Logging Middleware
	m.LoggingMiddleware.Register(app)

	// Register CORS Middleware
	if m.Config.BaseConfig.EnableCors {
		app.Use(
			cors.New(
				cors.Config{
					// Next:             nil,
					AllowOrigins:     m.Config.MiddlewareConfig.CorsConfig.AllowOrigins,
					AllowMethods:     m.Config.MiddlewareConfig.CorsConfig.AllowMethods,
					AllowHeaders:     m.Config.MiddlewareConfig.CorsConfig.AllowHeaders,
					AllowCredentials: m.Config.MiddlewareConfig.CorsConfig.AllowCredentials,
					ExposeHeaders:    m.Config.MiddlewareConfig.CorsConfig.ExposeHeaders,
					MaxAge:           m.Config.MiddlewareConfig.CorsConfig.MaxAge,
				},
			),
		)
	}

	// Register limiter Middleware
	app.Use(
		limiter.New(
			limiter.Config{
				Max:               m.Config.MiddlewareConfig.LimiterConfig.Max,
				Expiration:        m.Config.MiddlewareConfig.LimiterConfig.Expiration,
				LimiterMiddleware: limiter.SlidingWindow{},
				LimitReached: func(c *fiber.Ctx) error {
					return errors.ServerBusy(fmt.Errorf("too many requests"))
				},
			},
		),
	)

	// Register Auth Middleware
	// m.AuthMiddleware.Register(app)

	// Register Prometheus Middleware
	if m.Config.BaseConfig.EnablePrometheus {
		if err := m.PrometheusMiddleware.Register(app); err != nil {
			return err
		}
	}
	// Register Context Middleware
	m.ContextMiddleware.Register(app)

	return nil
}
