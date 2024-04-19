package route_handlers

import (
	proxy_config "github.com/itsindigo/reverse-proxy/internal/proxy_config"
	"log"
	http "net/http"
)

func RegisterProxyRoute(mux *http.ServeMux, route proxy_config.Route) {
	mux.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Redirect to METHOD: %s, HOST: %s, PORT: %s PATH: %s\n", route.Target.Method, route.Target.Host, route.Target.Port, route.Target.Path)
	})
	log.Printf("Registered route: %v", route.Path)
}
