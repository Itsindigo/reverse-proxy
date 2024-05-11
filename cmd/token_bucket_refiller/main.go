package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/itsindigo/reverse-proxy/internal/app_config"
	"github.com/itsindigo/reverse-proxy/internal/connections"
	"github.com/itsindigo/reverse-proxy/internal/proxy_configuration"
	"github.com/itsindigo/reverse-proxy/internal/repositories"
	"github.com/itsindigo/reverse-proxy/internal/services/rate_limiter"
)

var ctx = context.Background()

func start() {
	/**
	* 1. Establish routes that are monitored for bucket refill
	* 2. Query Buckets that match those route patterns every * n seconds (sleep go routine per route interval?)
	* 3. Increment counters
	 */

	config := app_config.NewConfig()
	rc := connections.CreateRedisClient(ctx, config.Redis)
	repositories := repositories.CreateApplicationRepositories(rc)
	rls := rate_limiter.NewRateLimiterService(repositories)
	routes, err := proxy_configuration.Load("../../RouteDefinitions.yml")

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	refillTasks := make([]rate_limiter.BucketRefillTask, 0)

	for _, route := range routes {
		refillTasks = append(refillTasks, rate_limiter.BucketRefillTask{
			Pattern:                rls.GetUserHttpRequestLimitKeyPattern(ctx, route),
			IncrementEveryNSeconds: 60 / route.RateLimit.RequestsPerMinute,
			MaxTokens:              route.RateLimit.RequestsPerMinute,
		})
	}

	for {
		time.Sleep(time.Second * 1)
	}
}

func main() {
	go start()
	select {}
}
