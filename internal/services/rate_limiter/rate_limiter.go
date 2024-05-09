package rate_limiter

import (
	"context"
	"fmt"

	"github.com/itsindigo/reverse-proxy/internal/constants"
	"github.com/itsindigo/reverse-proxy/internal/repositories"
	"github.com/itsindigo/reverse-proxy/internal/repositories/token_bucket"
	"github.com/itsindigo/reverse-proxy/internal/utils/crypto"
)

type RateLimiterService struct {
	TokenBucketRepository repository_token_bucket.TokenBucketRepository
}

func (rls *RateLimiterService) GetUserRouteLevelRequestKey(ctx context.Context, userIP string, method constants.HttpMethod, path string) string {
	fmt.Printf("IP: %v\nEncoded: %v\n", userIP, crypto.B64Encode(userIP))
	return fmt.Sprintf("user_route_requests:%s:%s:%s", crypto.B64Encode(userIP), method, path)
}

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

func NewRateLimiterService(repositories *repositories.ApplicationRepositories) *RateLimiterService {
	return &RateLimiterService{TokenBucketRepository: *repositories.TokenBucket}
}
