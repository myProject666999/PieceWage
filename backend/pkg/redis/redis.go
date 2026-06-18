package redis

import (
	"context"
	"fmt"

	"piece-wage/internal/config"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client
var Ctx = context.Background()

func Init(cfg *config.RedisConfig) error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	if err := RDB.Ping(Ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect redis: %w", err)
	}
	return nil
}

func Close() error {
	if RDB != nil {
		return RDB.Close()
	}
	return nil
}
