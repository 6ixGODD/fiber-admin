package mods

import (
	"fmt"

	"fiber-admin/internal/pkg/config"
	"fiber-admin/internal/pkg/service/common/mods"
	"fiber-admin/pkg/errors"
	"github.com/gofiber/fiber/v2"
)

type IdempotencyMiddleware struct {
	IdempotencyService mods.IdempotencyService
	Config             *config.Config
}

func (m *IdempotencyMiddleware) IdempotencyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		token := c.Get(m.Config.IdempotencyConfig.IdempotencyTokenHeader)
		if token != "" {
			err := m.IdempotencyService.CheckIdempotencyToken(ctx, token)
			if err != nil {
				return err
			} else {
				return c.Next()
			}
		}
		return errors.Idempotency(fmt.Errorf("idempotency token missed"))
	}
}
