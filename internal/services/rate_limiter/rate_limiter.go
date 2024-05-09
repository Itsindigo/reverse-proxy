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
	return fmt.Sprintf("user_route_requests:%s:%s:%s", crypto.B64Encode([]byte(userIP)), method, path)
}

// TODO: Track global requests as well as route level requests.
func (rls *RateLimiterService) GetUserGlobalRequestKey(ctx context.Context, userIP string) string {
	return fmt.Sprintf("user_global_requests:%s", crypto.B64Encode([]byte(userIP)))
}

func (rls *RateLimiterService) GetTokenBucket(ctx context.Context, requestKey string, limit int) (*repository_token_bucket.TokenBucket, error) {
	bucket, err := rls.TokenBucketRepository.GetOrCreateTokenBucket(ctx, requestKey, limit)

	if err != nil {
		return nil, err
	}

	return bucket, nil
}

func (rls *RateLimiterService) ApplyRequest(ctx context.Context, bucket *repository_token_bucket.TokenBucket) error {
	return nil
}

func NewRateLimiterService(repositories *repositories.ApplicationRepositories) *RateLimiterService {
	return &RateLimiterService{TokenBucketRepository: *repositories.TokenBucket}
}
