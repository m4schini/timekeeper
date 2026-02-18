package api

import (
	"net"
	"raumzeitalpaka/app/database"

	"github.com/go-chi/chi/v5"
)

func Serve(listener net.Listener, db *database.Database) {
	r := NewRouter(db)
}

func NewRouter(db *database.Database) chi.Router {
	mux := chi.NewMux()
	mux.Route("/event", func(r chi.Router) {
		r.Route("/{eventID}", func(r chi.Router) {
			r.Get("/", GetArticle)       // GET /articles/123
			r.Post("/", GetArticle)      // POST /articles/123
			r.Put("/", UpdateArticle)    // PUT /articles/123
			r.Delete("/", DeleteArticle) // DELETE /articles/123
		})

		r.Get("/schedule", GetSchedule)

		r.Get("/locations", GetLocationsOfevent)
		r.Route("/location/{locationID}", func(r chi.Router) {
			r.Put("/", UpdateLocationToEvent)
			r.Delete("/", RemoveLocationFromEvent)
		})

		r.Route("/room/{eventID}", func(r chi.Router) {
			r.Get("/", GetArticle)       // GET /articles/123
			r.Post("/", GetArticle)      // POST /articles/123
			r.Put("/", UpdateArticle)    // PUT /articles/123
			r.Delete("/", DeleteArticle) // DELETE /articles/123
		})
	})
	mux.Route("/room", func(r chi.Router) {
		r.Route("/{eventID}", func(r chi.Router) {
			r.Get("/", GetArticle)       // GET /articles/123
			r.Post("/", GetArticle)      // POST /articles/123
			r.Put("/", UpdateArticle)    // PUT /articles/123
			r.Delete("/", DeleteArticle) // DELETE /articles/123
		})

		r.Get("/schedule", GetSchedule)
	})
	mux.Route("/timeslot/{timeslotID}", func(r chi.Router) {
		r.Get("/", GetArticle)
		r.Post("/", GetArticle)
		r.Put("/", UpdateArticle)
	})

	return mux
}
