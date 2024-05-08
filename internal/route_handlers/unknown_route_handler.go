package route_handlers

import (
	"log/slog"
	"net/http"
)

func UnknownRouteHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		slog.Info("No route registed at path", "path", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello proxy server"))
}
