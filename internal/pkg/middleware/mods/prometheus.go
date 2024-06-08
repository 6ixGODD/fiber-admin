package mods

import (
	"context"

	logging "fiber-admin/pkg/zap"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"go.uber.org/zap"

	"fiber-admin/pkg/prometheus"
)

type PrometheusMiddleware struct {
	Prometheus *prometheus.Prometheus
	Zap        *logging.Zap
}

func (p *PrometheusMiddleware) Register(app *fiber.App) error {
	if err := p.setupPath(app); err != nil {
		return err
	}
	app.Use(p.Prometheus.PrometheusFiberHandler())
	return nil
}

func (p *PrometheusMiddleware) setupPath(app *fiber.App) error {
	ctx := context.Background()
	p.Zap.SetTagInContext(ctx, logging.SystemTag)
	logger, err := p.Zap.GetLogger(ctx)
	if err != nil {
		return err
	}
	logger.Info(
		"Prometheus middleware setup path", zap.String("path", p.Prometheus.PrometheusConfig.MetricPath),
	)
	app.Get(
		p.Prometheus.PrometheusConfig.MetricPath, func(c *fiber.Ctx) error {
			h := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
			h(c.Context())
			return nil
		},
	)
	return nil
}
