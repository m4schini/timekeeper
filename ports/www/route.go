package www

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

type Route interface {
	Method() string
	Pattern() string
	UseCache() bool
	Handler() http.Handler
}

func HandleRoute(router chi.Router, route Route) {
	router.Method(route.Method(), route.Pattern(), route.Handler())

	zap.L().Named("ports").Named("www").
		Debug("added route",
			zap.String("method", route.Method()),
			zap.String("route", route.Pattern()),
			zap.Bool("useCache", route.UseCache()),
		)
}
