package mods

import (
	"fiber-admin/internal/pkg/api/v1/common"
	"fiber-admin/internal/pkg/config"
	"github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
)

type CommonRouter struct {
	Config *config.Config
}

// RegisterCommonRouter registers the common router.
func (c *CommonRouter) RegisterCommonRouter(
	app fiber.Router, api *common.Common, casbin *casbin.Middleware, authMiddleware fiber.Handler,
) {
	app.Get(
		"/ping", func(c *fiber.Ctx) error { return c.SendString("pong") },
	)
	app.Get(
		"/idempotency-token",
		authMiddleware,
		casbin.RequiresRoles([]string{config.UserRoleAdmin, config.UserRoleUser}),
		api.IdempotencyApi.GenerateIdempotencyToken,
	)
	app.Get(
		"/profile",
		authMiddleware,
		casbin.RequiresRoles([]string{config.UserRoleAdmin, config.UserRoleUser}),
		api.ProfileApi.GetProfile,
	)
	app.Put(
		"/change-password",
		authMiddleware,
		casbin.RequiresRoles([]string{config.UserRoleAdmin, config.UserRoleUser}),
		api.AuthApi.ChangePassword,
	)

	authGroup := app.Group("/auth")
	authGroup.Post(
		"/login",
		api.AuthApi.Login,
	)
	authGroup.Get(
		"/logout",
		authMiddleware,
		api.AuthApi.Logout,
	)
	authGroup.Post(
		"/refresh",
		authMiddleware,
		api.AuthApi.RefreshToken,
	)

	noticeGroup := app.Group("/notice")
	noticeGroup.Get(
		"/",
		authMiddleware,
		api.NoticeApi.GetNotice,
	)
	noticeGroup.Get(
		"/list",
		authMiddleware,
		api.NoticeApi.GetNoticeList,
	)

	documentationGroup := app.Group("/documentation")
	documentationGroup.Get(
		"/",
		authMiddleware,
		api.DocumentationApi.GetDocumentation,
	)
	documentationGroup.Get(
		"/list",
		authMiddleware,
		api.DocumentationApi.GetDocumentationList,
	)
}
