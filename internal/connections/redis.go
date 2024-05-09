package connections

import (
	"context"
	"fmt"

	"github.com/itsindigo/reverse-proxy/internal/app_config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func CreateRedisClient(ctx context.Context, options app_config.RedisConfig) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", options.Host, options.Port),
		Password: options.Password, // no password set
		DB:       options.Database, // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	rdb.Del(ctx, "key")
	return &RedisClient{Client: rdb}
}
