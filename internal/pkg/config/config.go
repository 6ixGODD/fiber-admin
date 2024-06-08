package config

import (
	"sync"

	"fiber-admin/internal/pkg/config/mods"
	"github.com/mcuadros/go-defaults"
)

var (
	configInstance = &Config{}
	once           sync.Once
)

type Config struct {
	BaseConfig        mods.BaseConfig        `mapstructure:"base" yaml:"base"`
	CasbinConfig      mods.CasbinConfig      `mapstructure:"casbin" yaml:"casbin"`
	FiberConfig       mods.FiberConfig       `mapstructure:"fiber" yaml:"fiber"`
	JWTConfig         mods.JWTConfig         `mapstructure:"jwt" yaml:"jwt"`
	MongoConfig       mods.MongoConfig       `mapstructure:"mongo" yaml:"mongo"`
	PrometheusConfig  mods.PrometheusConfig  `mapstructure:"prometheus" yaml:"prometheus"`
	MiddlewareConfig  mods.MiddlewareConfig  `mapstructure:"middleware" yaml:"middleware"`
	CacheConfig       mods.CacheConfig       `mapstructure:"cache" yaml:"cache"`
	TasksConfig       mods.TasksConfig       `mapstructure:"tasks" yaml:"tasks"`
	ZapConfig         mods.ZapConfig         `mapstructure:"zap" yaml:"zap"`
	IdempotencyConfig mods.IdempotencyConfig `mapstructure:"idempotency" yaml:"idempotency"`
}

// New returns instance of Config
func New() *Config {
	once.Do(
		func() {
			defaults.SetDefaults(configInstance)
		},
	)
	return configInstance
}

func Set(cfg *Config) {
	configInstance = cfg
}
