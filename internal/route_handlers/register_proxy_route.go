package route_handlers

import (
	"fmt"
	"log"
	http "net/http"
	"net/http/httputil"
	"net/url"

	"github.com/itsindigo/reverse-proxy/internal/route_config"
)

func RegisterProxyRoute(mux *http.ServeMux, route route_config.Route) {
	mux.HandleFunc(fmt.Sprintf("%s %s", route.Target.Method, route.Path), func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Receiving %v", route)
		target := fmt.Sprintf("http://%s%s%s", route.Target.Host, route.Target.Port, route.Target.Path)
		parsedUrl, err := url.Parse(target)

		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		proxy := &httputil.ReverseProxy{
			Director: func(pr *http.Request) {
				targetURL := url.URL{
					Scheme: parsedUrl.Scheme,
					Host:   parsedUrl.Host,
					Path:   parsedUrl.Path,
				}
				pr.URL = &targetURL
			},
		}

		proxy.ServeHTTP(w, r)
	})

	log.Printf("Registered route: %v", route.Path)
}
