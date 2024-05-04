package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/itsindigo/reverse-proxy/internal/app_config"
	"github.com/itsindigo/reverse-proxy/internal/connections"
	route_config "github.com/itsindigo/reverse-proxy/internal/route_config"
	"github.com/itsindigo/reverse-proxy/internal/route_handlers"
)

var ctx = context.Background()

func main() {
	mux := http.NewServeMux()

	config := app_config.NewConfig()

	rc := connections.CreateRedisClient(ctx, config.Redis)

	routes, err := route_config.Load("./route_definitions.yml")

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	mux.HandleFunc("/", route_handlers.UnknownRouteHandler)

	for _, route := range routes {
		route_handlers.RegisterProxyRoute(mux, rc, route)
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
