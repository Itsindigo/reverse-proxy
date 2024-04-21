package main

import (
	"fmt"
	"log"
	"net/http"

	app_config "github.com/itsindigo/reverse-proxy/internal/config"
	route_config "github.com/itsindigo/reverse-proxy/internal/route_config"
	"github.com/itsindigo/reverse-proxy/internal/route_handlers"
)

func main() {
	mux := http.NewServeMux()

	config := app_config.NewConfig()

	fmt.Println(config)

	routes, err := route_config.Load("./route_definitions.yml")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	mux.HandleFunc("/", route_handlers.UnknownRouteHandler)

	for _, route := range routes {
		route_handlers.RegisterProxyRoute(mux, route)
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
