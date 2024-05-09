package route_handlers

import (
	"net/http"
)

func Handle429(w http.ResponseWriter) {
	http.Error(w, "Rate Limit Exceeded", http.StatusTooManyRequests)
}
