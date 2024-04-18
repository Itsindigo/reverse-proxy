package main

import (
	"fmt"
	route_map_parser "github.com/itsindigo/reverse-proxy/internal/route_map"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":6666",
		Handler: mux,
	}

	route_map_parser.Parse("./route_map.yml")
	fmt.Println("Server is listening on http://localhost:6666")
	err := server.ListenAndServe()

	if err != nil {
		fmt.Printf("Error starting server %s\n", err)
	}

}
