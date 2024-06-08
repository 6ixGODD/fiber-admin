package mods

import (
	e "errors"
	"fmt"

	"fiber-admin/internal/pkg/config"
	"fiber-admin/internal/pkg/dao"
	"fiber-admin/pkg/errors"
	auth "fiber-admin/pkg/jwt"
	"fiber-admin/pkg/utils/check"
	"fiber-admin/pkg/utils/crypt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type AuthMiddleware struct {
	Jwt    *auth.Jwt
	Cache  *dao.Cache
	Config *config.Config
}

func (a *AuthMiddleware) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get(fiber.HeaderAuthorization)
		if token == "" {
			return errors.TokenMissed(fmt.Errorf("token missed"))
		}
		if !check.IsBearerToken(token) {
			return errors.TokenInvalid(fmt.Errorf("token should be bearer token (start with 'Bearer ' or 'bearer ')"))
		}
		token = token[7:] // remove 'Bearer '
		blacklistKey := fmt.Sprintf("%s:%s", config.TokenBlacklistCachePrefix, crypt.MD5(token))
		if ok, err := a.Cache.Get(c.Context(), blacklistKey); err == nil && *ok == config.CacheTrue {
			return errors.TokenInvalid(fmt.Errorf("token has been revoked"))
		}
		sub, err := a.Jwt.VerifyAccessToken(token)
		if err != nil {
			var ve *jwt.ValidationError
			if e.As(err, &ve) {
				if ve.Errors == jwt.ValidationErrorExpired {
					return errors.TokenExpired(fmt.Errorf("token is expired"))
				}
			}
			return errors.TokenInvalid(fmt.Errorf("token invalid"))
		}
		c.Locals(config.UserIDKey, sub)
		return c.Next()
	}
}
