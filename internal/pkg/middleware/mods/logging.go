package mods

import (
	logging "fiber-admin/pkg/zap"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	Zap *logging.Zap
}

func (l *LoggingMiddleware) Register(app *fiber.App) {
	app.Use(l.loggingMiddleware())
}

func (l *LoggingMiddleware) loggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			// ctx                  = c.UserContext()
			sysLogger, reqLogger *zap.Logger
			err                  error
			sysCtx               = l.Zap.SetTagInContext(c.Context(), logging.SystemTag)
			reqCtx               = l.Zap.SetTagInContext(c.Context(), logging.RequestTag)
		)
		sysLogger, _ = l.Zap.GetLogger(sysCtx)
		reqLogger, _ = l.Zap.GetLogger(reqCtx)
		err = c.Next()
		if err != nil {
			sysLogger.Warn("Failed to execute request", zap.Error(err))
		}

		if c.Response().StatusCode() >= fiber.StatusInternalServerError {
			reqLogger.Error(
				"Request",
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
				zap.String("ip", c.IP()),
				zap.String("userAgent", c.Get(fiber.HeaderUserAgent)),
				zap.Any("query", c.Request().URI().QueryArgs()),
				zap.Any("form", c.Request().PostArgs()),
				zap.Any("body", c.Body()),
				zap.Int("status", c.Response().StatusCode()),
				zap.Any("response", c.Response().Body()),
			)
		} else if c.Response().StatusCode() >= fiber.StatusBadRequest {
			reqLogger.Warn(
				"Request",
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
				zap.String("ip", c.IP()),
				zap.String("userAgent", c.Get(fiber.HeaderUserAgent)),
				zap.Any("query", c.Request().URI().QueryArgs()),
				zap.Any("form", c.Request().PostArgs()),
				zap.Any("body", c.Body()),
				zap.Int("status", c.Response().StatusCode()),
				zap.Any("response", c.Response().Body()),
			)
		} else {
			reqLogger.Info(
				"Request",
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
				zap.String("ip", c.IP()),
				zap.String("userAgent", c.Get(fiber.HeaderUserAgent)),
				zap.Any("query", c.Request().URI().QueryArgs()),
				zap.Any("form", c.Request().PostArgs()),
				zap.Any("body", c.Body()),
				zap.Int("status", c.Response().StatusCode()),
				zap.Any("response", c.Response().Body()),
			)
		}

		return err
	}
}
