package mods

import (
	"time"

	"fiber-admin/internal/pkg/config/mods/cache"
)

type CacheConfig struct {
	DefaultTTL            time.Duration     `mapstructure:"default_ttl" yaml:"default_ttl" default:"5m"`
	UserCacheTTL          time.Duration     `mapstructure:"user_cache_ttl" yaml:"user_cache_ttl" default:"5m"`
	NoticeCacheTTL        time.Duration     `mapstructure:"notice_cache_ttl" yaml:"notice_cache_ttl" default:"5m"`
	DocumentationCacheTTL time.Duration     `mapstructure:"documentation_cache_ttl" yaml:"documentation_cache_ttl" default:"5m"`
	TokenBlacklistTTL     time.Duration     `mapstructure:"token_blacklist_ttl" yaml:"token_blacklist_ttl" default:"1h"`
	RedisConfig           cache.RedisConfig `mapstructure:"redis" yaml:"redis"`
}
