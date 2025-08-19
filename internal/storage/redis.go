package storage

import (
	"context"
	"fmt"

	"github.com/SOMTHING-ITPL/ITPL-server/config"
	"github.com/go-redis/redis/v8"
)

func InitRedis(cfg config.RedisDBConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	return rdb, nil
}
