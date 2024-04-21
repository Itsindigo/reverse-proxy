package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /hello", helloHandler)
	mux.HandleFunc("GET /", homeHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server is listening on http://localhost:8080")
	err := server.ListenAndServe()

	if err != nil {
		fmt.Printf("Error starting server %s\n", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s", r.URL)

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server one is OK!"))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Hello string `json:"hello"`
	}{
		Hello: "World",
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)

	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}
