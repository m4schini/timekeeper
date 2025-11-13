package www

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net"
	"net/http"
	"raumzeitalpaka/app/auth"
	"raumzeitalpaka/config"
	"raumzeitalpaka/ports/www/middleware"
	"strings"
)

func Serve(listener net.Listener, authenticator auth.Authenticator, pages []Route, components []Route) error {
	r := chi.NewRouter()
	r.Use(
		http.NewCrossOriginProtection().Handler,
		middleware.UseAuth(authenticator),
		middleware.Log,
	)
	for _, route := range pages {
		HandleRoute(r, route)
	}
	r.Route("/_", func(r chi.Router) {
		for _, route := range components {
			HandleRoute(r, route)
		}
	})

	if config.TelemetryEnabled() {
		EnableMetricsEndpoint(r)
	} else {
		zap.L().Named("telemetry").Info("metrics endpoint is disabled")
	}

	return http.Serve(listener, r)
}

func EnableMetricsEndpoint(r chi.Router) {
	zap.L().Named("telemetry").Info("metrics endpoint is enabled", zap.String("route", "/metrics"))
	token := config.MetricsEndpointToken()
	next := promhttp.Handler()
	r.Handle("/metrics", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if token == "" {
			next.ServeHTTP(writer, request)
			return
		}

		authHeader := request.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer") && strings.HasSuffix(authHeader, token) {
			next.ServeHTTP(writer, request)
			return
		}

		http.Error(writer, "unauthorized", http.StatusUnauthorized)
	}))
}
