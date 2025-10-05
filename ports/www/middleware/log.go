package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

func Log(next http.Handler) http.Handler {
	log := zap.L().Named("api")
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		path := request.URL.Path
		if len(path) > 96 {
			path = path[:16] + "..." + path[len(path)-16:]
		}

		log := log.With(zap.String("method", request.Method), zap.String("path", path))
		log.Debug("received api request")

		next.ServeHTTP(writer, request)

		log.Debug("handled api request", zap.Duration("duration", time.Since(start)))
	})
}
