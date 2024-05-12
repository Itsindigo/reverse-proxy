package rate_limiter

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/itsindigo/reverse-proxy/internal/constants"
	"github.com/itsindigo/reverse-proxy/internal/proxy_configuration"
	"github.com/itsindigo/reverse-proxy/internal/repositories"
	"github.com/itsindigo/reverse-proxy/internal/repositories/token_bucket"
	"github.com/itsindigo/reverse-proxy/internal/utils/crypto"
)

type BucketRefillTask struct {
	Pattern                string
	IncrementEveryNSeconds int
	MaxTokens              int
}

type RateLimiterService struct {
	TokenBucketRepository repository_token_bucket.TokenBucketRepository
}

// GetUserHttpRequestLimitKey creates a static Redis Key when given a UserIP and Route Info.
func (rls *RateLimiterService) GetUserHttpRequestLimitKey(ctx context.Context, userIP string, method constants.HttpMethod, path string) string {
	return fmt.Sprintf("%s:%s:%s:%s", constants.UserHttpRequestLimit, crypto.B64Encode(userIP), method, path)
}

// GetUserHttpRequestLimitKey creates a Redis Key Pattern to find requests against a given route.
func (rls *RateLimiterService) GetUserHttpRequestLimitKeyPattern(ctx context.Context, route proxy_configuration.Route) string {
	return fmt.Sprintf("%s:*:%s:%s", constants.UserHttpRequestLimit, route.Method, route.Path)
}

// GetTokenBucket returns a token bucket instance for a given request key.
func (rls *RateLimiterService) GetTokenBucket(ctx context.Context, requestKey string, limit int) (*repository_token_bucket.TokenBucket, error) {
	bucket, err := rls.TokenBucketRepository.GetOrCreateTokenBucket(ctx, requestKey, limit)

	if err != nil {
		return nil, err
	}

	return bucket, nil
}

func (rls *RateLimiterService) ApplyRequest(ctx context.Context, bucket *repository_token_bucket.TokenBucket) error {
	_, err := rls.TokenBucketRepository.ConsumeToken(ctx, bucket)

	if err != nil {
		return err
	}

	fmt.Printf("Token Count: %d\n", bucket.TokenCount)

	return nil
}

func (rls *RateLimiterService) CreateRefillTask(ctx context.Context, task BucketRefillTask) func() {
	// TODO:
	// Query redis for keys matching task pattern
	// Increment value by 1 if key value is < task.MaxTokens
	// Sleep until task start time + IncrementEveryNSeconds rounded down to nearest second
	return func() {
		for {
			continueAt := time.Now().Truncate(time.Second).Add(time.Duration(task.IncrementEveryNSeconds) * time.Second)

			increment := func(key string, value interface{}) error {
				valStr, ok := value.(string)

				if !ok {
					return fmt.Errorf("skipping redis key %q as value found was mistyped. expected: string, got: %T", key, value)
				}

				tokenCount, err := strconv.Atoi(valStr)
				if err != nil {
					return fmt.Errorf("skipping redis key %q as value could not be converted to int, err: %v", key, err)
				}

				if tokenCount >= task.MaxTokens {
					slog.Info("Bucket full")
					return nil
				}

				err = rls.TokenBucketRepository.SetKey(ctx, key, tokenCount+1, 0)

				if err != nil {
					return fmt.Errorf("could not set redis key %q, err: %v", key, err)
				}

				slog.Info(fmt.Sprintf("Incremented key %q by one, new value %d", key, tokenCount))

				return nil
			}

			rls.TokenBucketRepository.MapKeys(ctx, task.Pattern, increment)
			<-time.After(time.Until(continueAt))
		}
	}
}

func NewRateLimiterService(repositories *repositories.ApplicationRepositories) *RateLimiterService {
	return &RateLimiterService{TokenBucketRepository: *repositories.TokenBucket}
}
