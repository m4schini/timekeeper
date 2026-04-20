package www

import (
	"net/http"
	"raumzeitalpaka/ports/www/middleware"
	"reflect"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Route interface {
	Method() string
	Pattern() string
	Handler() http.Handler
}

func HandleRoute(router chi.Router, route Route) {
	if route.Method() == "" {
		router.Handle(route.Pattern(), middleware.UseGzip(route.Handler()))
	} else {
		router.Method(route.Method(), route.Pattern(), middleware.UseGzip(route.Handler()))
	}

	zap.L().Named("ports").Named("www").
		Info("added route",
			zap.String("method", route.Method()),
			zap.String("route", route.Pattern()),
			zap.Any("handler", reflect.TypeOf(route)),
		)
}
