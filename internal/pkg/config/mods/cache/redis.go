package cache

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr            string        `mapstructure:"redis_addr" yaml:"redis_addr" default:"localhost:6379"`
	ClientName      string        `mapstructure:"redis_client_name" yaml:"redis_client_name" default:""`
	Username        string        `mapstructure:"redis_username" yaml:"redis_username" default:""`
	Password        string        `mapstructure:"redis_password" yaml:"redis_password" default:""`
	DB              int           `mapstructure:"redis_db" yaml:"redis_db" default:"0"`
	MaxRetries      int           `mapstructure:"redis_max_retries" yaml:"redis_max_retries" default:"5"`
	MinRetryBackoff time.Duration `mapstructure:"redis_min_retry_backoff" yaml:"redis_min_retry_backoff" default:"8ms"`
	MaxRetryBackoff time.Duration `mapstructure:"redis_max_retry_backoff" yaml:"redis_max_retry_backoff" default:"512ms"`
	DialTimeout     time.Duration `mapstructure:"redis_dial_timeout" yaml:"redis_dial_timeout" default:"5s"`
	ReadTimeout     time.Duration `mapstructure:"redis_read_timeout" yaml:"redis_read_timeout" default:"3s"`
	WriteTimeout    time.Duration `mapstructure:"redis_write_timeout" yaml:"redis_write_timeout" default:"3s"`
	PoolSize        int           `mapstructure:"redis_pool_size" yaml:"redis_pool_size" default:"10"`
	PoolTimeout     time.Duration `mapstructure:"redis_pool_timeout" yaml:"redis_pool_timeout" default:"4s"`
	MinIdleConns    int           `mapstructure:"redis_min_idle_conns" yaml:"redis_min_idle_conns" default:"0"`
	MaxIdleConns    int           `mapstructure:"redis_max_idle_conns" yaml:"redis_max_idle_conns" default:"0"`
	MaxActiveConns  int           `mapstructure:"redis_max_active_conns" yaml:"redis_max_active_conns" default:"0"`
	ConnMaxIdleTime time.Duration `mapstructure:"redis_conn_max_idle_time" yaml:"redis_conn_max_idle_time" default:"30m"`
	ConnMaxLifetime time.Duration `mapstructure:"redis_conn_max_lifetime" yaml:"redis_conn_max_lifetime" default:"-1"`
}

func (c *RedisConfig) GetRedisOptions() *redis.Options {
	return &redis.Options{
		Addr:            c.Addr,
		ClientName:      c.ClientName,
		Username:        c.Username,
		Password:        c.Password,
		DB:              c.DB,
		MaxRetries:      c.MaxRetries,
		MinRetryBackoff: c.MinRetryBackoff,
		MaxRetryBackoff: c.MaxRetryBackoff,
		DialTimeout:     c.DialTimeout,
		ReadTimeout:     c.ReadTimeout,
		WriteTimeout:    c.WriteTimeout,
		PoolSize:        c.PoolSize,
		PoolTimeout:     c.PoolTimeout,
		MinIdleConns:    c.MinIdleConns,
		MaxIdleConns:    c.MaxIdleConns,
		MaxActiveConns:  c.MaxActiveConns,
		ConnMaxIdleTime: c.ConnMaxIdleTime,
		ConnMaxLifetime: c.ConnMaxLifetime,
	}
}
