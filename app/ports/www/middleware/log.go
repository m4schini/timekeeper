package middleware

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/config"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

func Log(next http.Handler) http.Handler {
	log := zap.L().Named("www").WithOptions(zap.AddCallerSkip(1))
	telemtryEnabled := config.TelemetryEnabled()
	counter := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "raumzeitalpaka",
		Subsystem: "www",
		Name:      "requests",
	}, []string{"method", "status", "route"})

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

		sr := RecordStatus(writer)
		next.ServeHTTP(sr, request)

		d := time.Since(start)
		if telemtryEnabled {
			counter.WithLabelValues(request.Method, fmt.Sprintf("%v", sr.status), request.URL.Path).Inc()
		}

		log = log.With(zap.Int("status", sr.status), zap.Duration("duration", d))
		if d > 100*time.Millisecond {
			log.Warn("handled www request")
		} else {
			log.Info("handled www request")
		}
	})
}

func RecordStatus(w http.ResponseWriter) *statusRecorder {
	return &statusRecorder{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

// statusRecorder wraps http.ResponseWriter to capture the status code
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
