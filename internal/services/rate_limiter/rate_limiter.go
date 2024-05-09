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

func (rls *RateLimiterService) GetRequestKey(userIp string, method constants.HttpMethod, path string) string {
	return fmt.Sprintf("route_bucket:%s:%s:%s", crypto.B64Encode([]byte(userIp)), method, path)
}

func (rls *RateLimiterService) ApplyRequest(ctx context.Context, requestKey string) error {

	bucket, err := rls.TokenBucketRepository.GetOrCreateTokenBucket(ctx, requestKey, 10)

	if err != nil {
		return err
	}

	return nil
}

func NewRateLimiterService(repositories *repositories.ApplicationRepositories) *RateLimiterService {
	return &RateLimiterService{TokenBucketRepository: *repositories.TokenBucket}
}
