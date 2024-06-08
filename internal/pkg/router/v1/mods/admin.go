package mods

import (
	"fiber-admin/internal/pkg/api/v1/admin"
	"fiber-admin/internal/pkg/config"
	"github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
)

const adminPrefix = "/admin"

type AdminRouter struct{}

// RegisterAdminRouter registers the admin router.
func (a *AdminRouter) RegisterAdminRouter(
	app fiber.Router, api *admin.Admin, casbin *casbin.Middleware, idempotencyMiddleware, authMiddleware fiber.Handler,
) {
	group := app.Group(adminPrefix)

	group.Post(
		"/notice",
		authMiddleware,
		casbin.RequiresRoles([]string{config.UserRoleAdmin}),
		api.NoticeApi.InsertNotice,
	)
	group.Put(
		"/notice",
		authMiddleware,
		casbin.RequiresRoles([]string{config.UserRoleAdmin}),
		api.NoticeApi.UpdateNotice,
	)
	group.Delete(
		"/notice",
		authMiddleware,
		casbin.RequiresRoles([]string{config.UserRoleAdmin}),
		api.NoticeApi.DeleteNotice,
	)

	group.Post(
		"/user",
		authMiddleware,
		idempotencyMiddleware, // An example of using idempotency middleware. Actually not necessary.
		casbin.RequiresRoles([]string{config.UserRoleAdmin}),
		api.UserApi.InsertUser,
	)
	group.Get(
		"/user",
		authMiddleware,
		casbin.RequiresRoles([]string{config.UserRoleAdmin}),
		api.UserApi.GetUser,
	)
	group.Get(
		"/user/list",
		authMiddleware,
		casbin.RequiresRoles([]string{config.UserRoleAdmin}),
		api.UserApi.GetUserList,
	)

	group.Put(
		"/user",
		authMiddleware,
		casbin.RequiresRoles([]string{config.UserRoleAdmin}),
		api.UserApi.UpdateUser,
	)
	group.Delete(
		"/user",
		authMiddleware,
		casbin.RequiresRoles([]string{config.UserRoleAdmin}),
		api.UserApi.DeleteUser,
	)
	group.Put(
		"/user/password",
		authMiddleware,
		casbin.RequiresRoles([]string{config.UserRoleAdmin}),
		api.UserApi.ChangeUserPassword,
	)

	group.Post(
		"/documentation",
		authMiddleware,
		casbin.RequiresRoles([]string{config.UserRoleAdmin}),
		api.DocumentationApi.InsertDocumentation,
	)
	group.Put(
		"/documentation",
		authMiddleware,
		casbin.RequiresRoles([]string{config.UserRoleAdmin}),
		api.DocumentationApi.UpdateDocumentation,
	)
	group.Delete(
		"/documentation",
		authMiddleware,
		casbin.RequiresRoles([]string{config.UserRoleAdmin}),
		api.DocumentationApi.DeleteDocumentation,
	)

	group.Get(
		"/login-log/list",
		authMiddleware,
		casbin.RequiresRoles([]string{config.UserRoleAdmin}),
		api.LogsApi.GetLoginLogList,
	)
	group.Get(
		"/operation-log/list",
		authMiddleware,
		casbin.RequiresRoles([]string{config.UserRoleAdmin}),
		api.LogsApi.GetOperationLogList,
	)
}
