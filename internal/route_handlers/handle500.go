package route_handlers

import (
	"fmt"
	"net/http"
)

func Handle500(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
}
