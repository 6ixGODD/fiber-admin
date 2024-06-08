package redis

import (
	"sync"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

type Redis struct {
	RedisClient *redis.Client
	redisConfig *Config
	mu          sync.Mutex
}

type Config struct {
	redisOptions *redis.Options
}

var (
	redisInstance *Redis // Singleton
	once          sync.Once
)

func New(ctx context.Context, options *redis.Options) (*Redis, error) {
	var err error
	once.Do(
		func() {
			r := &Redis{
				redisConfig: &Config{
					redisOptions: options,
				},
			}
			if err = r.Init(ctx); err != nil {
				return
			}
			redisInstance = r
		},
	)
	return redisInstance, err
}

func Update(ctx context.Context, options *redis.Options) error {
	var err error
	if err := redisInstance.RedisClient.Close(); err != nil {
		return err
	}
	*redisInstance = Redis{
		redisConfig: &Config{
			redisOptions: options,
		},
	}
	if err = redisInstance.Init(ctx); err != nil {
		return err
	}
	return nil
}

func (r *Redis) Init(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.RedisClient != nil {
		return nil
	}
	client := redis.NewClient(r.redisConfig.redisOptions)
	if _, err := client.Ping(ctx).Result(); err != nil {
		return err
	}
	r.RedisClient = client
	return nil
}

func (r *Redis) GetClient() (client *redis.Client, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.RedisClient == nil {
		if err = r.Init(context.Background()); err != nil {
			return nil, err
		}
	}
	return r.RedisClient, nil
}

func (r *Redis) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.RedisClient == nil {
		return nil
	}
	err := r.RedisClient.Close()
	r.RedisClient = nil
	return err
}
