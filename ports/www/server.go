package www

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net"
	"net/http"
	"timekeeper/app/auth"
	"timekeeper/config"
	"timekeeper/ports/www/middleware"
)

func Serve(listener net.Listener, authenticator auth.Authenticator, pages []Route, components []Route) error {
	r := chi.NewRouter()
	r.Use(
		middleware.AllowAllCORS,
		middleware.Log,
		//middleware.UseGzip,
		middleware.UseAuth(authenticator),
	)
	for _, route := range pages {
		HandleRoute(r, route)
	}
	r.Route("/_", func(r chi.Router) {
		//r.Use(middleware.AllowAllCORS, middleware.Log, middleware.UseGzip, middleware.UseAuth(authenticator))
		for _, route := range components {
			HandleRoute(r, route)
		}
	})

	if config.TelemetryEnabled() {
		zap.L().Named("telemetry").Info("telemetry is enabled")
		r.Handle("/metrics", promhttp.Handler())
	} else {
		zap.L().Named("telemetry").Info("telemetry is disabled")
	}

	return http.Serve(listener, r)
}
