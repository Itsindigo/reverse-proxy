package route_handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/itsindigo/reverse-proxy/internal/proxy_configuration"
	"github.com/itsindigo/reverse-proxy/internal/repositories"
	"github.com/itsindigo/reverse-proxy/internal/services/ip_utils"
	"github.com/itsindigo/reverse-proxy/internal/services/rate_limiter"
)

func RegisterProxyRoute(ctx context.Context, mux *http.ServeMux, repos *repositories.ApplicationRepositories, route proxy_configuration.Route) {
	RateLimiterService := rate_limiter.NewRateLimiterService(repos)
	mux.HandleFunc(fmt.Sprintf("%s %s", route.Method, route.Path), func(w http.ResponseWriter, r *http.Request) {
		target := fmt.Sprintf("http://%s%s%s", route.Target.Host, route.Target.Port, route.Target.Path)
		parsedUrl, err := url.Parse(target)

		if err != nil {
			Handle500(w, err)
			return
		}

		userIP, err := ip_utils.GetIP(r.RemoteAddr, r.Header.Get("X-Forwarded-For"))

		if err != nil {
			Handle500(w, err)
			return
		}

		requestKey := RateLimiterService.GetUserRouteLevelRequestKey(ctx, userIP, route.Method, route.Target.Path)
		bucket, err := RateLimiterService.GetTokenBucket(ctx, requestKey, route.RateLimit.RequestsPerMinute)
		RateLimiterService.ApplyRequest(ctx, bucket)

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

	slog.Info("Registered route", "path", route.Path)
}
