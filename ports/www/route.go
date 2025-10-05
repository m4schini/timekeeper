package www

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Route interface {
	Method() string
	Pattern() string
	Handler() http.Handler
}

func HandleRoute(router chi.Router, route Route) {
	router.Method(route.Method(), route.Pattern(), route.Handler())
}
