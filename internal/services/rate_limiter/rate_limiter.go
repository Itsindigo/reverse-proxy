package rate_limiter

import (
	"context"
	"fmt"
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
	bucket, err := rls.TokenBucketRepository.ConsumeToken(ctx, bucket)

	if err != nil {
		return err
	}

	fmt.Printf("Token Count: %d\n", bucket.TokenCount)

	return nil
}

func (rls *RateLimiterService) CreateRefillTask(task BucketRefillTask) func() {
	return func() {
		for {
			fmt.Printf("Executing callback for: %s\n", task.Pattern)
			time.Sleep(time.Second * 5)
		}
	}
}

func NewRateLimiterService(repositories *repositories.ApplicationRepositories) *RateLimiterService {
	return &RateLimiterService{TokenBucketRepository: *repositories.TokenBucket}
}
