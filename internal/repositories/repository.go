package repositories

import (
	"github.com/itsindigo/reverse-proxy/internal/connections"
	"github.com/itsindigo/reverse-proxy/internal/repositories/token_bucket"
)

type ApplicationRepositories struct {
	TokenBucket *repository_token_bucket.TokenBucketRepository
}

func CreateApplicationRepositories(rc *connections.RedisClient) *ApplicationRepositories {
	return &ApplicationRepositories{
		TokenBucket: repository_token_bucket.NewTokenBucketRepository(rc),
	}
}
