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
	"github.com/itsindigo/reverse-proxy/internal/utils/math_utils"
)

type BucketRefillTask struct {
	Pattern                   string
	IncrementNTokensPerSecond int
	MaxTokens                 int
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

	return nil
}

func (rls *RateLimiterService) CreateRefillTask(ctx context.Context, task BucketRefillTask) func() {
	return func() {
		for {
			continueAt := time.Now().Truncate(time.Second).Add(1 * time.Second)

			increment := func(key string, value interface{}) error {
				valStr, ok := value.(string)

				if !ok {
					return fmt.Errorf("skipping redis key %q as value found was mistyped. expected: string, got: %T", key, value)
				}
				tokenCount, err := strconv.Atoi(valStr)
				if err != nil {
					return fmt.Errorf("skipping redis key %q as value could not be converted to int, err: %v", key, err)
				}

				newTokenCount := tokenCount + task.IncrementNTokensPerSecond
				newTokenCount = math_utils.Min(newTokenCount, task.MaxTokens)

				if tokenCount == newTokenCount {
					return nil
				}

				err = rls.TokenBucketRepository.SetKey(ctx, key, newTokenCount, 0)

				if err != nil {
					return fmt.Errorf("could not set redis key %s, err: %v", key, err)
				}

				slog.Info("Refilling Bucket Key", slog.String("key", key), slog.Int("new_value", newTokenCount))
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
