package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/itsindigo/reverse-proxy/internal/app_config"
	"github.com/itsindigo/reverse-proxy/internal/connections"
	"github.com/itsindigo/reverse-proxy/internal/proxy_configuration"
	"github.com/itsindigo/reverse-proxy/internal/repositories"
	"github.com/itsindigo/reverse-proxy/internal/services/rate_limiter"
)

func start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	config := app_config.NewConfig()
	rc := connections.CreateRedisClient(ctx, config.Redis)
	repositories := repositories.CreateApplicationRepositories(rc)
	rls := rate_limiter.NewRateLimiterService(repositories)
	routes, err := proxy_configuration.Load("./RouteDefinitions.yml")

	if err != nil {
		log.Fatalf("Error loading route configurations: %v", err)
	}

	refillTasks := createRefillTasks(ctx, rls, routes)

	for _, refillTask := range refillTasks {
		wg.Add(1)
		go refillTask(ctx, wg)
	}

	<-ctx.Done()
}

func createRefillTasks(ctx context.Context, rls *rate_limiter.RateLimiterService, routes []proxy_configuration.Route) []func(ctx context.Context, wg *sync.WaitGroup) {
	refillTasks := make([]func(ctx context.Context, wg *sync.WaitGroup), 0, len(routes))

	for _, route := range routes {
		pattern := rls.GetUserHttpRequestLimitKeyPattern(ctx, route)
		slog.Info("Starting Refill Task For Pattern", slog.String("pattern", rls.GetUserHttpRequestLimitKeyPattern(ctx, route)))

		refillTask := rls.CreateRefillTask(ctx, rate_limiter.BucketRefillTask{
			Pattern:                   pattern,
			IncrementNTokensPerSecond: route.RateLimit.RequestsPerMinute / 60,
			MaxTokens:                 route.RateLimit.RequestsPerMinute,
		})

		refillTasks = append(refillTasks, refillTask)
	}
	return refillTasks
}

func main() {
	var ctx, cancel = context.WithCancel(context.Background())
	var wg sync.WaitGroup

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		slog.Info("Received Shutdown Signal", slog.String("signal", sig.String()))
		cancel()
	}()

	wg.Add(1)
	go start(ctx, &wg)
	wg.Wait()
	slog.Info("All goroutines closed, exiting.")
}
