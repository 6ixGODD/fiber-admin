// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"context"
	"fiber-admin/internal/app"
	"fiber-admin/internal/pkg/api/v1"
	"fiber-admin/internal/pkg/api/v1/admin"
	mods4 "fiber-admin/internal/pkg/api/v1/admin/mods"
	"fiber-admin/internal/pkg/api/v1/common"
	mods6 "fiber-admin/internal/pkg/api/v1/common/mods"
	"fiber-admin/internal/pkg/config"
	"fiber-admin/internal/pkg/dao"
	"fiber-admin/internal/pkg/dao/mods"
	"fiber-admin/internal/pkg/middleware"
	mods8 "fiber-admin/internal/pkg/middleware/mods"
	router2 "fiber-admin/internal/pkg/router"
	"fiber-admin/internal/pkg/router/v1"
	mods7 "fiber-admin/internal/pkg/router/v1/mods"
	"fiber-admin/internal/pkg/service"
	admin2 "fiber-admin/internal/pkg/service/admin"
	mods2 "fiber-admin/internal/pkg/service/admin/mods"
	common2 "fiber-admin/internal/pkg/service/common"
	mods5 "fiber-admin/internal/pkg/service/common/mods"
	"fiber-admin/internal/pkg/service/sys"
	mods3 "fiber-admin/internal/pkg/service/sys/mods"
	"fiber-admin/internal/pkg/tasks"
	"fiber-admin/internal/pkg/validator"
	"github.com/google/wire"
)

// Injectors from wire.go:

// InitializeApp initialize app
func InitializeApp(ctx context.Context) (*app.App, error) {
	configConfig := config.New()
	zap, err := InitializeZap(configConfig)
	if err != nil {
		return nil, err
	}
	core, err := service.NewCore(ctx, configConfig, zap)
	if err != nil {
		return nil, err
	}
	mongo, err := InitializeMongo(ctx, configConfig)
	if err != nil {
		return nil, err
	}
	daoCore, err := dao.NewCore(ctx, mongo, zap, configConfig)
	if err != nil {
		return nil, err
	}
	redis, err := InitializeRedis(ctx, configConfig)
	if err != nil {
		return nil, err
	}
	cache := dao.NewCache(redis, configConfig)
	userDao, err := mods.NewUserDao(ctx, daoCore, cache)
	if err != nil {
		return nil, err
	}
	enforcer, err := InitializeCasbinEnforcer(configConfig)
	if err != nil {
		return nil, err
	}
	userService := mods2.NewUserService(core, userDao, enforcer)
	loginLogDao, err := mods.NewLoginLogDao(ctx, daoCore, cache, userDao)
	if err != nil {
		return nil, err
	}
	operationLogDao, err := mods.NewOperationLogDao(ctx, daoCore, cache, userDao)
	if err != nil {
		return nil, err
	}
	logsService := mods3.NewLogsService(core, loginLogDao, operationLogDao)
	validate, err := validator.NewValidator()
	if err != nil {
		return nil, err
	}
	userApi := &mods4.UserApi{
		UserService: userService,
		LogsService: logsService,
		Validator:   validate,
	}
	noticeDao, err := mods.NewNoticeDao(ctx, daoCore, cache)
	if err != nil {
		return nil, err
	}
	noticeService := mods2.NewNoticeService(core, noticeDao)
	noticeApi := &mods4.NoticeApi{
		NoticeService: noticeService,
		LogsService:   logsService,
		Validator:     validate,
	}
	documentationDao, err := mods.NewDocumentationDao(ctx, daoCore, cache)
	if err != nil {
		return nil, err
	}
	documentationService := mods2.NewDocumentationService(core, documentationDao)
	documentationApi := &mods4.DocumentationApi{
		DocumentationService: documentationService,
		LogsService:          logsService,
		Validator:            validate,
	}
	modsLogsService := mods2.NewLogsService(core, loginLogDao, operationLogDao)
	logsApi := &mods4.LogsApi{
		LogsService: modsLogsService,
		Validator:   validate,
	}
	adminAdmin := &admin.Admin{
		UserApi:          userApi,
		NoticeApi:        noticeApi,
		DocumentationApi: documentationApi,
		LogsApi:          logsApi,
	}
	jwt, err := InitializeJwt(configConfig)
	if err != nil {
		return nil, err
	}
	authService := mods5.NewAuthService(core, userDao, cache, jwt)
	authApi := &mods6.AuthApi{
		AuthService: authService,
		LogsService: logsService,
		Validator:   validate,
	}
	profileService := mods5.NewProfileService(core, userDao)
	profileApi := &mods6.ProfileApi{
		ProfileService: profileService,
	}
	modsDocumentationService := mods5.NewDocumentationService(core, documentationDao)
	modsDocumentationApi := &mods6.DocumentationApi{
		DocumentationService: modsDocumentationService,
		Validator:            validate,
	}
	modsNoticeService := mods5.NewNoticeService(core, noticeDao)
	modsNoticeApi := &mods6.NoticeApi{
		NoticeService: modsNoticeService,
		Validator:     validate,
	}
	idempotencyService := mods5.NewIdempotencyService(core, cache)
	idempotencyApi := &mods6.IdempotencyApi{
		IdempotencyService: idempotencyService,
	}
	commonCommon := &common.Common{
		AuthApi:          authApi,
		ProfileApi:       profileApi,
		DocumentationApi: modsDocumentationApi,
		NoticeApi:        modsNoticeApi,
		IdempotencyApi:   idempotencyApi,
	}
	apiApi := &api.Api{
		AdminApi:  adminAdmin,
		CommonApi: commonCommon,
	}
	adminRouter := &mods7.AdminRouter{}
	commonRouter := &mods7.CommonRouter{
		Config: configConfig,
	}
	routerRouter := &router.Router{
		ApiV1:        apiApi,
		AdminRouter:  adminRouter,
		CommonRouter: commonRouter,
	}
	router3 := &router2.Router{
		RouterV1: routerRouter,
	}
	authMiddleware := &mods8.AuthMiddleware{
		Jwt:    jwt,
		Cache:  cache,
		Config: configConfig,
	}
	loggingMiddleware := &mods8.LoggingMiddleware{
		Zap: zap,
	}
	prometheus := InitializePrometheus(configConfig)
	prometheusMiddleware := &mods8.PrometheusMiddleware{
		Prometheus: prometheus,
		Zap:        zap,
	}
	contextMiddleware := &mods8.ContextMiddleware{
		Zap: zap,
	}
	idempotencyMiddleware := &mods8.IdempotencyMiddleware{
		IdempotencyService: idempotencyService,
		Config:             configConfig,
	}
	middlewareMiddleware := &middleware.Middleware{
		AuthMiddleware:        authMiddleware,
		LoggingMiddleware:     loggingMiddleware,
		PrometheusMiddleware:  prometheusMiddleware,
		ContextMiddleware:     contextMiddleware,
		IdempotencyMiddleware: idempotencyMiddleware,
		Config:                configConfig,
	}
	tasksTasks, err := tasks.New(ctx, configConfig, loginLogDao, operationLogDao, jwt, zap)
	if err != nil {
		return nil, err
	}
	appApp, err := app.New(ctx, zap, configConfig, router3, middlewareMiddleware, tasksTasks, mongo, redis)
	if err != nil {
		return nil, err
	}
	return appApp, nil
}

// wire.go:

var (
	RouterProviderSet = wire.NewSet(wire.Struct(new(mods7.AdminRouter), "*"), wire.Struct(new(mods7.CommonRouter), "*"), wire.Struct(new(router.Router), "*"), wire.Struct(new(router2.Router), "*"))

	ApiProviderSet = wire.NewSet(wire.Struct(new(mods6.AuthApi), "*"), wire.Struct(new(mods6.ProfileApi), "*"), wire.Struct(new(mods6.DocumentationApi), "*"), wire.Struct(new(mods6.NoticeApi), "*"), wire.Struct(new(mods6.IdempotencyApi), "*"), wire.Struct(new(mods4.UserApi), "*"), wire.Struct(new(mods4.DocumentationApi), "*"), wire.Struct(new(mods4.NoticeApi), "*"), wire.Struct(new(mods4.LogsApi), "*"), wire.Struct(new(common.Common), "*"), wire.Struct(new(admin.Admin), "*"), wire.Struct(new(api.Api), "*"))

	ValidatorProviderSet = wire.NewSet(validator.NewValidator)

	ServiceProviderSet = wire.NewSet(service.NewCore, wire.Struct(new(admin2.Admin), "*"), wire.Struct(new(common2.Common), "*"), wire.Struct(new(sys.Sys), "*"), mods2.NewUserService, mods2.NewNoticeService, mods2.NewDocumentationService, mods2.NewLogsService, mods5.NewAuthService, mods5.NewProfileService, mods5.NewDocumentationService, mods5.NewNoticeService, mods5.NewIdempotencyService, mods3.NewLogsService)

	DaoProviderSet = wire.NewSet(dao.NewCore, dao.NewCache, mods.NewUserDao, mods.NewNoticeDao, mods.NewLoginLogDao, mods.NewOperationLogDao, mods.NewDocumentationDao)

	MiddlewareProviderSet = wire.NewSet(wire.Struct(new(mods8.LoggingMiddleware), "*"), wire.Struct(new(mods8.PrometheusMiddleware), "*"), wire.Struct(new(mods8.AuthMiddleware), "*"), wire.Struct(new(mods8.ContextMiddleware), "*"), wire.Struct(new(mods8.IdempotencyMiddleware), "*"), wire.Struct(new(middleware.Middleware), "*"))

	SchedulerProviderSet = wire.NewSet(tasks.New)
)
