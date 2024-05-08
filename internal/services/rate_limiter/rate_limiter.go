package rate_limiter

import (
	"fmt"
	"net/http"

	"github.com/itsindigo/reverse-proxy/internal/constants"
	"github.com/itsindigo/reverse-proxy/internal/repositories"
	"github.com/itsindigo/reverse-proxy/internal/repositories/token_bucket"
	"github.com/itsindigo/reverse-proxy/internal/services/ip_utils"
)

type RateLimiterService struct {
	TokenBucketRepository repository_token_bucket.TokenBucketRepository
}

func (rls *RateLimiterService) RegisterUserIpBucket(r *http.Request, method constants.HttpMethod, path string) error {
	requestKey, err := ip_utils.GetIpRequestKey(r)

	if err != nil {
		return err
	}

	// TODO Figure out this pattern + hash symbols
	requestKey = fmt.Sprintf("route_bucket:%s:%s:%s", requestKey, method, path)

	// bucket, err := rls.TokenBucketRepository.GetOrCreateTokenBucket(ctx, "some-key", route.RateLimit.RequestsPerMinute)
	return nil
}

func NewRateLimiterService(repositories *repositories.ApplicationRepositories) *RateLimiterService {
	return &RateLimiterService{TokenBucketRepository: *repositories.TokenBucket}
}
