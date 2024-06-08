package middleware

import (
	"time"
)

type LimiterConfig struct {
	Max        int           `mapstructure:"limiter_max" yaml:"limiter_max" default:"20"`
	Expiration time.Duration `mapstructure:"limiter_expiration" yaml:"limiter_expiration" default:"30s"`
}
