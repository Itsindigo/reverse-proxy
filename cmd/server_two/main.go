package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("GET /goodbye", goodbyeHandler)

	server := &http.Server{
		Addr:    ":9090",
		Handler: mux,
	}

	fmt.Println("Server is listening on http://localhost:9090")
	err := server.ListenAndServe()

	if err != nil {
		fmt.Printf("Error starting server %s\n", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server two is OK!"))
}

func goodbyeHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Goodbye string `json:"goodbye"`
	}{
		Goodbye: "world",
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)

	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}
