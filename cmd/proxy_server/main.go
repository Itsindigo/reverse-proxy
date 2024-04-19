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

	mux.HandleFunc("/", unknownRouteHandler)

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

func unknownRouteHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Printf("No route registed at: %v", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello proxy server"))
}
