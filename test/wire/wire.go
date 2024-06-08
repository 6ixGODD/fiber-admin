//go:build wireinject
// +build wireinject

package wire

import (
	"context"

	"fiber-admin/internal/pkg/config"
	"fiber-admin/internal/pkg/dao"
	daos "fiber-admin/internal/pkg/dao/mods"
	"fiber-admin/internal/pkg/service"
	adminservice "fiber-admin/internal/pkg/service/admin"
	adminservices "fiber-admin/internal/pkg/service/admin/mods"
	commonservice "fiber-admin/internal/pkg/service/common"
	commonservices "fiber-admin/internal/pkg/service/common/mods"
	sysservice "fiber-admin/internal/pkg/service/sys"
	sysservices "fiber-admin/internal/pkg/service/sys/mods"
	"fiber-admin/pkg/jwt"
	"fiber-admin/pkg/mongo"
	"fiber-admin/pkg/prometheus"
	"fiber-admin/pkg/redis"
	logging "fiber-admin/pkg/zap"
	"fiber-admin/test/mock"
	"github.com/casbin/casbin/v2"
	"github.com/google/wire"
)

type Injector struct {
	// Common
	Ctx        context.Context
	Config     *config.Config
	Cache      *dao.Cache
	Mongo      *mongo.Mongo
	Redis      *redis.Redis
	Zap        *logging.Zap
	Jwt        *jwt.Jwt
	Prometheus *prometheus.Prometheus

	// DAOs
	UserDao          daos.UserDao
	NoticeDao        daos.NoticeDao
	DocumentationDao daos.DocumentationDao
	LoginLogDao      daos.LoginLogDao
	OperationLogDao  daos.OperationLogDao

	// Mocks for DAOs
	UserDaoMock          *mock.UserDaoMock
	NoticeDaoMock        *mock.NoticeDaoMock
	DocumentationDaoMock *mock.DocumentationDaoMock
	LoginLogDaoMock      *mock.LoginLogDaoMock
	OperationLogDaoMock  *mock.OperationLogDaoMock

	// Services
	// Admin services
	AdminDocumentationService adminservices.DocumentationService
	AdminNoticeService        adminservices.NoticeService
	AdminLogsService          adminservices.LogsService
	AdminUserService          adminservices.UserService
	// Common services
	CommonAuthService          commonservices.AuthService
	CommonIdempotencyService   commonservices.IdempotencyService
	CommonDocumentationService commonservices.DocumentationService
	CommonNoticeService        commonservices.NoticeService
	CommonProfileService       commonservices.ProfileService
	// Sys services
	SysLogsService sysservices.LogsService

	// Casbin enforcer
	Enforcer *casbin.Enforcer
}

var (
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

	MockProviderSet = wire.NewSet(
		mock.NewUserDaoMockWithRandomData,
		mock.NewNoticeDaoMockWithRandomData,
		mock.NewLoginLogDaoMockWithRandomData,
		mock.NewOperationLogDaoMockWithRandomData,
		mock.NewDocumentationDaoMockWithRandomData,
	)
)

func InitializeTestInjector(ctx context.Context, config *config.Config, n int) (*Injector, error) {
	wire.Build(
		InitializeMongo,
		InitializeRedis,
		InitializeZap,
		InitializeJwt,
		InitializePrometheus,
		InitializeCasbinEnforcer,
		MockProviderSet,
		ServiceProviderSet,
		DaoProviderSet,
		wire.Struct(new(Injector), "*"),
	)
	return new(Injector), nil
}
