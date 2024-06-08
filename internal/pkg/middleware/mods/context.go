package mods

import (
	"fiber-admin/internal/pkg/config"
	logging "fiber-admin/pkg/zap"
	"github.com/gofiber/fiber/v2"
)

type ContextMiddleware struct {
	Zap *logging.Zap
}

func (m *ContextMiddleware) Register(app *fiber.App) {
	app.Use(m.contextMiddleware())
}

func (m *ContextMiddleware) contextMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		rid, ok := c.Locals(config.RequestIDKey).(string)
		if ok {
			ctx = m.Zap.SetRequestIDInContext(ctx, rid)
		}
		userID, ok := c.Locals(config.UserIDKey).(string)
		if ok {
			ctx = m.Zap.SetUserIDInContext(ctx, userID)
		}
		c.SetUserContext(ctx)
		return c.Next()
	}
}
