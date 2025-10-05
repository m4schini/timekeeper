package www

import (
	"github.com/go-chi/chi/v5"
	"net"
	"net/http"
	"timekeeper/app/auth"
	"timekeeper/ports/www/middleware"
)

func Serve(listener net.Listener, authenticator auth.Authenticator, pages []Route, components []Route) error {
	r := chi.NewRouter()
	r.Use(middleware.AllowAllCORS, middleware.Log, middleware.UseGzip, middleware.UseAuth(authenticator))
	for _, route := range pages {
		HandleRoute(r, route)
	}
	r.Route("/_", func(r chi.Router) {
		//r.Use(middleware.AllowAllCORS, middleware.Log, middleware.UseGzip, middleware.UseAuth(authenticator))
		for _, route := range components {
			HandleRoute(r, route)
		}
	})

	return http.Serve(listener, r)
}
