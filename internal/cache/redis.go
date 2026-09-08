package cache

import (
	"context"
	"crud_service/internal/config"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{Addr: cfg.Addr()})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed connect to redis")
	}

	return client, nil
}
