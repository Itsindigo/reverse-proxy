package repository_token_bucket

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/itsindigo/reverse-proxy/internal/app_errors"
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

	if err = repo.rc.Client.Set(ctx, key, limit, 0).Err(); err != nil {
		return nil, fmt.Errorf("error setting key %q with limit %d: %w", key, limit, err)
	}

	if err = repo.rc.Client.Expire(ctx, key, 24*time.Hour).Err(); err != nil {
		return nil, fmt.Errorf("error setting expiry on key %q", key)
	}

	return &TokenBucket{Key: key, TokenCount: limit}, nil
}

func (repo *TokenBucketRepository) ConsumeToken(ctx context.Context, bucket *TokenBucket) (*TokenBucket, error) {
	val, err := repo.rc.Client.Get(ctx, bucket.Key).Result()

	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("could not find bucket %q", bucket.Key)
	}

	tokenCount, err := strconv.Atoi(val)

	if err != nil {
		return nil, fmt.Errorf("bucket did not contain a valid int, got [%T, %q]", val, val)
	}

	if tokenCount <= 0 {
		return nil, app_errors.BucketExhaustedError{}
	}

	if err = repo.rc.Client.Decr(ctx, bucket.Key).Err(); err != nil {
		return nil, fmt.Errorf("error occurred while decrememnting count: %v", err)
	}
	bucket.TokenCount -= 1

	return bucket, nil
}

func NewTokenBucketRepository(rc *connections.RedisClient) *TokenBucketRepository {
	return &TokenBucketRepository{rc: rc}
}
