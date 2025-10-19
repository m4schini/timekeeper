package middleware

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"net/http"
	"time"
	"timekeeper/config"
)

func Log(next http.Handler) http.Handler {
	log := zap.L().Named("www").WithOptions(zap.AddCallerSkip(1))
	telemtryEnabled := config.TelemetryEnabled()
	counter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "timekeeper",
		Subsystem: "www",
		Name:      "requests",
	})
	requestDuration := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "timekeeper",
		Subsystem: "www",
		Name:      "request_durations",
		Buckets:   []float64{1, 5, 10, 20, 40, 60, 100, 150, 200, 250, 300, 600, 1200, 6000},
	})

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()

		logFields := []zap.Field{
			zap.String("route", fmt.Sprintf("%v %v", request.Method, request.URL.Path)),
		}

		userId, role, isAuthenticated := LoadUser(request)
		logFields = append(logFields, zap.Bool("is_authenticated", isAuthenticated))
		if isAuthenticated {
			logFields = append(logFields, zap.Int("user_id", userId), zap.Any("role", role))
		}

		log := log.With(logFields...)
		log.Debug("received www request")

		next.ServeHTTP(writer, request)

		d := time.Since(start)
		if telemtryEnabled {
			counter.Inc()
			requestDuration.Observe(float64(d.Milliseconds()))
		}

		log.Info("handled www request", zap.Duration("duration", d))
	})
}
