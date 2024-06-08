package mods

import (
	"fiber-admin/internal/pkg/config/mods/middleware"
)

type MiddlewareConfig struct {
	LimiterConfig middleware.LimiterConfig `mapstructure:"limiter" yaml:"limiter"`
	CorsConfig    middleware.CorsConfig    `mapstructure:"cors" yaml:"cors"`
}
