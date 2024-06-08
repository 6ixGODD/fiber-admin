package mods

import (
	"time"
)

type FiberConfig struct {
	Prefork                 bool          `mapstructure:"prefork" yaml:"prefork" default:"false"`
	ServerHeader            string        `mapstructure:"server_header" yaml:"server_header" default:""`
	BodyLimit               int           `mapstructure:"body_limit" yaml:"body_limit" default:"4194304"`  // 4 * 1024 * 1024
	Concurrency             int           `mapstructure:"concurrency" yaml:"concurrency" default:"262144"` // 256 * 1024
	ReadTimeout             time.Duration `mapstructure:"read_timeout" yaml:"read_timeout" default:"10s"`
	WriteTimeout            time.Duration `mapstructure:"write_timeout" yaml:"write_timeout" default:"10s"`
	IdleTimeout             time.Duration `mapstructure:"idle_timeout" yaml:"idle_timeout" default:"2m"`
	ReadBufferSize          int           `mapstructure:"read_buffer_size" yaml:"read_buffer_size" default:"4096"`
	WriteBufferSize         int           `mapstructure:"write_buffer_size" yaml:"write_buffer_size" default:"4096"`
	ProxyHeader             string        `mapstructure:"proxy_header" yaml:"proxy_header" default:"X-Forwarded-For"`
	DisableKeepalive        bool          `mapstructure:"disable_keepalive" yaml:"disable_keepalive" default:"false"`
	DisableStartupMessage   bool          `mapstructure:"disable_startup_message" yaml:"disable_startup_message" default:"true"`
	ReduceMemoryUsage       bool          `mapstructure:"reduce_memory_usage" yaml:"reduce_memory_usage" default:"false"`
	EnableTrustedProxyCheck bool          `mapstructure:"enable_trusted_proxy_check" yaml:"enable_trusted_proxy_check" default:"false"`
	TrustedProxies          []string      `mapstructure:"trusted_proxies" yaml:"trusted_proxies" default:""`
	EnablePrintRoutes       bool          `mapstructure:"enable_print_routes" yaml:"enable_print_routes" default:"true"`
}
