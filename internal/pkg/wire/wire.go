//go:build wireinject
// +build wireinject

package wire

import (
	"context"

	"fiber-admin/internal/app"
	adminapi "fiber-admin/internal/pkg/api/v1/admin"
	adminapis "fiber-admin/internal/pkg/api/v1/admin/mods"
	commonapi "fiber-admin/internal/pkg/api/v1/common"
	commonapis "fiber-admin/internal/pkg/api/v1/common/mods"
	"fiber-admin/internal/pkg/config"
	"fiber-admin/internal/pkg/dao"
	daos "fiber-admin/internal/pkg/dao/mods"
	"fiber-admin/internal/pkg/middleware"
	wares "fiber-admin/internal/pkg/middleware/mods"
	"fiber-admin/internal/pkg/router"
	routerv1 "fiber-admin/internal/pkg/router/v1"
	routers "fiber-admin/internal/pkg/router/v1/mods"
	"fiber-admin/internal/pkg/service"
	adminservice "fiber-admin/internal/pkg/service/admin"
	adminservices "fiber-admin/internal/pkg/service/admin/mods"
	commonservice "fiber-admin/internal/pkg/service/common"
	commonservices "fiber-admin/internal/pkg/service/common/mods"
	sysservice "fiber-admin/internal/pkg/service/sys"
	sysservices "fiber-admin/internal/pkg/service/sys/mods"
	"fiber-admin/internal/pkg/tasks"
	"fiber-admin/internal/pkg/validator"
	"github.com/google/wire"
)

var (
	RouterProviderSet = wire.NewSet(
		wire.Struct(new(routers.AdminRouter), "*"),
		wire.Struct(new(routers.CommonRouter), "*"),
		wire.Struct(new(routerv1.Router), "*"),
		wire.Struct(new(router.Router), "*"),
	)

	ApiProviderSet = wire.NewSet(
		wire.Struct(new(commonapis.AuthApi), "*"),
		wire.Struct(new(commonapis.ProfileApi), "*"),
		wire.Struct(new(commonapis.DocumentationApi), "*"),
		wire.Struct(new(commonapis.NoticeApi), "*"),
		wire.Struct(new(commonapis.IdempotencyApi), "*"),
		wire.Struct(new(adminapis.UserApi), "*"),
		wire.Struct(new(adminapis.DocumentationApi), "*"),
		wire.Struct(new(adminapis.NoticeApi), "*"),
		wire.Struct(new(adminapis.LogsApi), "*"),
		wire.Struct(new(commonapi.Common), "*"),
		wire.Struct(new(adminapi.Admin), "*"),
		wire.Struct(new(api.Api), "*"),
	)

	ValidatorProviderSet = wire.NewSet(
		validator.NewValidator,
	)

	ServiceProviderSet = wire.NewSet(
		service.NewCore,
		wire.Struct(new(adminservice.Admin), "*"),
		wire.Struct(new(commonservice.Common), "*"),
		wire.Struct(new(sysservice.Sys), "*"),
		adminservices.NewUserService,
		adminservices.NewNoticeService,
		adminservices.NewDocumentationService,
		adminservices.NewLogsService,
		commonservices.NewAuthService,
		commonservices.NewProfileService,
		commonservices.NewDocumentationService,
		commonservices.NewNoticeService,
		commonservices.NewIdempotencyService,
		sysservices.NewLogsService,
	)

	DaoProviderSet = wire.NewSet(
		dao.NewCore,
		dao.NewCache,
		daos.NewUserDao,
		daos.NewNoticeDao,
		daos.NewLoginLogDao,
		daos.NewOperationLogDao,
		daos.NewDocumentationDao,
	)

	MiddlewareProviderSet = wire.NewSet(
		wire.Struct(new(wares.LoggingMiddleware), "*"),
		wire.Struct(new(wares.PrometheusMiddleware), "*"),
		wire.Struct(new(wares.AuthMiddleware), "*"),
		wire.Struct(new(wares.ContextMiddleware), "*"),
		wire.Struct(new(wares.IdempotencyMiddleware), "*"),
		wire.Struct(new(middleware.Middleware), "*"),
	)

	SchedulerProviderSet = wire.NewSet(
		tasks.New,
	)
)

// InitializeApp initialize app
func InitializeApp(ctx context.Context) (*app.App, error) {
	wire.Build(
		config.New,
		InitializeMongo,
		InitializeRedis,
		InitializeZap,
		InitializeJwt,
		InitializePrometheus,
		InitializeCasbinEnforcer,
		DaoProviderSet,
		ServiceProviderSet,
		ValidatorProviderSet,
		ApiProviderSet,
		RouterProviderSet,
		MiddlewareProviderSet,
		SchedulerProviderSet,
		app.New,
	)
	return new(app.App), nil
}
