package prometheus

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	Namespace  string
	Subsystem  string
	MetricPath string
}

type Prometheus struct {
	PrometheusConfig *Config
	reqCount         *prometheus.CounterVec
	reqDuration      *prometheus.HistogramVec
}

func New(namespace string, subsystem string, metricPath string) *Prometheus {
	p := &Prometheus{
		PrometheusConfig: &Config{
			Namespace:  namespace,
			Subsystem:  subsystem,
			MetricPath: metricPath,
		},
	}
	p.init()
	return p
}

func (p *Prometheus) init() {
	p.reqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "requests_total",
			Namespace: p.PrometheusConfig.Namespace,
			Subsystem: p.PrometheusConfig.Subsystem,
			Help:      "Number of HTTP requests",
		},
		[]string{"status_code", "method", "path"},
	)

	p.reqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      "request_duration_seconds",
			Namespace: p.PrometheusConfig.Namespace,
			Subsystem: p.PrometheusConfig.Subsystem,
			Help:      "Duration of HTTP requests",
		}, []string{"method", "handler"},
	)
}

func (p *Prometheus) PrometheusFiberHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() == p.PrometheusConfig.MetricPath { // Skip metrics path
			err := c.Next()
			return err
		}

		start := time.Now()

		err := c.Next()

		r := c.Route()

		statusCode := strconv.Itoa(c.Response().StatusCode())
		elapsed := float64(time.Since(start).Nanoseconds()) / 1e9

		p.reqCount.With(
			prometheus.Labels{"status_code": statusCode, "method": c.Method(), "path": r.Path},
		).Inc()
		p.reqDuration.With(
			prometheus.Labels{"method": c.Method(), "handler": r.Path},
		).Observe(elapsed)
		return err
	}
}
