package repository_token_bucket

import (
	"context"
	"fmt"
	"strconv"

	"github.com/itsindigo/reverse-proxy/internal/connections"
	"github.com/redis/go-redis/v9"
)

type TokenBucketRepository struct {
	rc *connections.RedisClient
}

type TokenBucket struct {
	Key        string
	TokenCount int
}

// Create a token bucket at the key string provided, with a limit value.
// TODO: Set a lifetime expire on the key, that is renewed only when a request is made.
func (repo *TokenBucketRepository) GetOrCreateTokenBucket(ctx context.Context, key string, limit int) (*TokenBucket, error) {
	val, err := repo.rc.Client.Get(ctx, key).Result()

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

	if repo.rc.Client.Set(ctx, key, limit, 0).Err(); err != nil {
		return nil, fmt.Errorf("error setting key &q with limit %d: %w", key, limit, err)
	}

	return &TokenBucket{Key: key, TokenCount: limit}, nil
}

func NewTokenBucketRepository(rc *connections.RedisClient) *TokenBucketRepository {
	return &TokenBucketRepository{rc: rc}
}
