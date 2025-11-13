package www

import (
	"net/http"
	"raumzeitalpaka/ports/www/middleware"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Route interface {
	Method() string
	Pattern() string
	Handler() http.Handler
}

func HandleRoute(router chi.Router, route Route) {
	router.Method(route.Method(), route.Pattern(), middleware.UseGzip(route.Handler()))

	zap.L().Named("ports").Named("www").
		Info("added route",
			zap.String("method", route.Method()),
			zap.String("route", route.Pattern()),
		)
}
