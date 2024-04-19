package main

import (
	"fmt"
	"log"
	"net/http"

	proxy_config "github.com/itsindigo/reverse-proxy/internal/proxy_config"
	route_handlers "github.com/itsindigo/reverse-proxy/internal/route_handlers"
)

func main() {
	mux := http.NewServeMux()

	routes, err := proxy_config.Parse("./proxy_config.yml")

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for _, route := range routes {
		route_handlers.RegisterProxyRoute(mux, route)
	}

	server := &http.Server{
		Addr:    ":6666",
		Handler: mux,
	}

	log.Println("Server is listening on http://localhost:6666")
	err = server.ListenAndServe()

	if err != nil {
		fmt.Printf("Error starting server %s\n", err)
	}
}
