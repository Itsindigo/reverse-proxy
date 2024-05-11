package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/itsindigo/reverse-proxy/internal/app_config"
	"github.com/itsindigo/reverse-proxy/internal/connections"
	"github.com/itsindigo/reverse-proxy/internal/proxy_configuration"
	"github.com/itsindigo/reverse-proxy/internal/repositories"
	"github.com/itsindigo/reverse-proxy/internal/route_handlers"
)

var ctx = context.Background()

func main() {
	mux := http.NewServeMux()

	config := app_config.NewConfig()

	rc := connections.CreateRedisClient(ctx, config.Redis)
	repositories := repositories.CreateApplicationRepositories(rc)

	routes, err := proxy_configuration.Load("../../RouteDefinitions.yml")

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	mux.HandleFunc("/", route_handlers.UnknownRouteHandler)

	for _, route := range routes {
		route_handlers.RegisterProxyRoute(ctx, mux, repositories, route)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.ProxyServer.Port),
		Handler: mux,
	}

	log.Printf("Server is listening on http://localhost:%s", config.ProxyServer.Port)
	err = server.ListenAndServe()

	if err != nil {
		fmt.Printf("Error starting server %s\n", err)
	}
}
