package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/itsindigo/reverse-proxy/internal/app_config"
	"github.com/itsindigo/reverse-proxy/internal/connections"
	"github.com/itsindigo/reverse-proxy/internal/proxy_configuration"
	"github.com/itsindigo/reverse-proxy/internal/repositories"
	"github.com/itsindigo/reverse-proxy/internal/route_handlers"
)

func getRouteHandlers(ctx context.Context, repositories *repositories.ApplicationRepositories, routes []proxy_configuration.Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		route_handlers.RegisterProxyRoute(ctx, mux, repositories, route)
	}

	mux.HandleFunc("/", route_handlers.UnknownRouteHandler)
	return mux
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := app_config.NewConfig()
	rc := connections.CreateRedisClient(ctx, config.Redis)
	repositories := repositories.CreateApplicationRepositories(rc)
	routes, err := proxy_configuration.Load("./RouteDefinitions.yml")

	if err != nil {
		log.Fatalf("Error loading route configurations: %v", err)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.ProxyServer.Port),
		Handler: getRouteHandlers(ctx, repositories, routes),
	}

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP Server error: %v", err)
		}
	}()

	log.Printf("Server is listening on http://localhost:%s", config.ProxyServer.Port)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signals
	slog.Info("Received Shutdown Signal", slog.String("signal", sig.String()))

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	slog.Info("System shutdown complete, exiting.")
}
