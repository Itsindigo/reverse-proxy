package connections

import (
	"context"
	"fmt"
	"strconv"

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
	return &RedisClient{Client: rdb}
}

type TokenBucket struct {
	Key        string
	TokenCount int
}

// Create a token bucket at the key string provided, with a limit value.
// TODO: Set a lifetime expire on the key, that is renewed only when a request is made.
func (rc *RedisClient) CreateNewTokenBucket(ctx context.Context, key string, limit int) (*TokenBucket, error) {
	val, err := rc.Client.Get(ctx, key).Result()

	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("error retrieving key %q: %w", key, err)
	}

	if val != "" {
		tokenCount, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("key %q does not contain an integer: %v", key, err)
		}

		return &TokenBucket{Key: key, TokenCount: tokenCount}, nil
	}

	if rc.Client.Set(ctx, key, limit, 0).Err(); err != nil {
		return nil, fmt.Errorf("error setting key &q with limit %d: %w", key, limit, err)
	}

	return &TokenBucket{Key: key, TokenCount: limit}, nil
}
