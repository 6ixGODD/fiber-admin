package router

import (
	"fiber-admin/internal/pkg/api/v1"
	"fiber-admin/internal/pkg/router/v1/mods"
	"github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
)

const v1Prefix = "/v1"

type Router struct {
	ApiV1        *api.Api
	AdminRouter  *mods.AdminRouter
	CommonRouter *mods.CommonRouter
}

func (a *Router) RegisterRouter(
	router *fiber.Router, casbin *casbin.Middleware, idempotencyMiddleware, authMiddleware fiber.Handler,
) {
	a.registerV1Router(router, casbin, idempotencyMiddleware, authMiddleware)
}

func (a *Router) registerV1Router(
	router *fiber.Router, casbin *casbin.Middleware, idempotencyMiddleware, authMiddleware fiber.Handler,
) {
	v1Router := (*router).Group(v1Prefix)
	a.AdminRouter.RegisterAdminRouter(v1Router, a.ApiV1.AdminApi, casbin, idempotencyMiddleware, authMiddleware)
	a.CommonRouter.RegisterCommonRouter(v1Router, a.ApiV1.CommonApi, casbin, authMiddleware)
}
