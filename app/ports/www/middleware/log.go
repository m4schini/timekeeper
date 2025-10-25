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

// statusRecorder wraps http.ResponseWriter to capture the status code
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func Log(next http.Handler) http.Handler {
	log := zap.L().Named("www").WithOptions(zap.AddCallerSkip(1))
	telemtryEnabled := config.TelemetryEnabled()
	counter := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "timekeeper",
		Subsystem: "www",
		Name:      "requests",
	}, []string{"method", "status", "route"})
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

		sr := &statusRecorder{
			ResponseWriter: writer,
			status:         http.StatusOK,
		}
		next.ServeHTTP(sr, request)

		d := time.Since(start)
		if telemtryEnabled {
			counter.WithLabelValues(request.Method, fmt.Sprintf("%v", sr.status), request.URL.Path).Inc()
			requestDuration.Observe(float64(d.Milliseconds()))
		}

		if d > 100*time.Millisecond {
			log.Warn("handled www request", zap.Duration("duration", d))
		} else if request.Method != http.MethodGet {
			log.Info("handled www request", zap.Duration("duration", d))
		} else {
			log.Debug("handled www request", zap.Duration("duration", d))
		}
	})
}
