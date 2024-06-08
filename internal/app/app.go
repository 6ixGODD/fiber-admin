package app

import (
	"context"
	"fmt"

	"fiber-admin/internal/pkg/config"
	"fiber-admin/internal/pkg/errors"
	"fiber-admin/internal/pkg/middleware"
	"fiber-admin/internal/pkg/router"
	"fiber-admin/internal/pkg/tasks"
	e "fiber-admin/pkg/errors"
	"fiber-admin/pkg/mongo"
	"fiber-admin/pkg/redis"
	logging "fiber-admin/pkg/zap"
	"github.com/casbin/mongodb-adapter/v3"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type App struct {
	App        *fiber.App
	Zap        *logging.Zap
	Logger     *zap.Logger
	Config     *config.Config
	Router     *router.Router
	Middleware *middleware.Middleware
	Tasks      *tasks.Tasks
	Mongo      *mongo.Mongo
	Redis      *redis.Redis
	Ctx        context.Context
}

// New factory function that initializes the application and returns a fiber.App instance.
func New(
	ctx context.Context, zap *logging.Zap, config *config.Config, router *router.Router,
	middleware *middleware.Middleware, tasks *tasks.Tasks, mongo *mongo.Mongo, redis *redis.Redis,
) (*App, error) {
	app := &App{
		Zap:        zap,
		Config:     config,
		Router:     router,
		Middleware: middleware,
		Tasks:      tasks,
		Mongo:      mongo,
		Redis:      redis,
		Ctx:        ctx,
	}

	if err := app.Init(ctx); err != nil {
		return nil, err
	}
	return app, nil
}

func (a *App) Init(ctx context.Context) error {
	// Set logger
	ctx = a.Zap.SetTagInContext(ctx, logging.SystemTag)
	logger, err := a.Zap.GetLogger(ctx)
	if err != nil {
		return err
	}
	a.Logger = logger

	// Set Fiber app
	app := fiber.New(
		fiber.Config{
			Prefork:                 a.Config.FiberConfig.Prefork,
			ServerHeader:            a.Config.FiberConfig.ServerHeader,
			BodyLimit:               a.Config.FiberConfig.BodyLimit,
			Concurrency:             a.Config.FiberConfig.Concurrency,
			ReadTimeout:             a.Config.FiberConfig.ReadTimeout,
			WriteTimeout:            a.Config.FiberConfig.WriteTimeout,
			IdleTimeout:             a.Config.FiberConfig.IdleTimeout,
			ReadBufferSize:          a.Config.FiberConfig.ReadBufferSize,
			WriteBufferSize:         a.Config.FiberConfig.WriteBufferSize,
			ProxyHeader:             a.Config.FiberConfig.ProxyHeader,
			DisableStartupMessage:   a.Config.FiberConfig.DisableStartupMessage,
			AppName:                 a.Config.BaseConfig.AppName,
			ReduceMemoryUsage:       a.Config.FiberConfig.ReduceMemoryUsage,
			EnableTrustedProxyCheck: a.Config.FiberConfig.EnableTrustedProxyCheck,
			TrustedProxies:          a.Config.FiberConfig.TrustedProxies,
			EnablePrintRoutes:       a.Config.FiberConfig.EnablePrintRoutes,
			ErrorHandler:            errors.ErrorHandler, // Custom error handler
			JSONDecoder:             json.Unmarshal,      // Use go-json for enhanced JSON decoding performance
			JSONEncoder:             json.Marshal,
		},
	)

	// Register Middleware
	if err := a.Middleware.Register(app); err != nil {
		return err
	}

	// Set hooks
	app.Hooks().OnShutdown(
		func() error {
			return ShutdownHandler(a.Ctx, a)
		},
	)

	// Ping

	// Set Casbin
	adapter, err := mongodbadapter.NewAdapter(a.Config.CasbinConfig.PolicyAdapterUrl)
	if err != nil {
		return err
	}
	c := casbin.New(
		casbin.Config{
			ModelFilePath: a.Config.CasbinConfig.ModelPath,
			PolicyAdapter: adapter,
			Lookup: func(c *fiber.Ctx) string {
				return c.Locals(config.UserIDKey).(string)
			},
			Forbidden: func(ctx *fiber.Ctx) error {
				return e.PermissionDeny(fmt.Errorf("permission deny")) // will capture by error handler
			},
		},
	)

	// Register routers
	a.Router.RegisterRouter(
		app, c, a.Middleware.IdempotencyMiddleware.IdempotencyMiddleware(),
		a.Middleware.AuthMiddleware.AuthMiddleware(),
	)

	// Set app
	a.App = app

	// Start scheduled tasks
	if err := a.Tasks.Start(); err != nil {
		return err
	}
	return nil
}
