package main

import (
	"context"
	"fmt"
	"log"

	"github.com/itsindigo/reverse-proxy/internal/app_config"
	"github.com/itsindigo/reverse-proxy/internal/connections"
	"github.com/itsindigo/reverse-proxy/internal/proxy_configuration"
	"github.com/itsindigo/reverse-proxy/internal/repositories"
	"github.com/itsindigo/reverse-proxy/internal/services/rate_limiter"
)

var ctx = context.Background()

func start() {
	config := app_config.NewConfig()
	rc := connections.CreateRedisClient(ctx, config.Redis)
	repositories := repositories.CreateApplicationRepositories(rc)
	rls := rate_limiter.NewRateLimiterService(repositories)
	routes, err := proxy_configuration.Load("../../RouteDefinitions.yml")

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	refillTasks := make([]func(), 0)

	for _, route := range routes {
		fmt.Printf("Pattern: %s", rls.GetUserHttpRequestLimitKeyPattern(ctx, route))
		refillTasks = append(refillTasks, rls.CreateRefillTask(ctx, rate_limiter.BucketRefillTask{
			Pattern:                rls.GetUserHttpRequestLimitKeyPattern(ctx, route),
			IncrementEveryNSeconds: 60 / route.RateLimit.RequestsPerMinute,
			MaxTokens:              route.RateLimit.RequestsPerMinute,
		}))
	}

	for _, refillTask := range refillTasks {
		go refillTask()
	}
}

func main() {
	go start()
	select {}
}
