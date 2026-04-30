package api

import (
	"raumzeitalpaka/app/auth"
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/ports/api/handler"

	"github.com/go-chi/chi/v5"
)

func NewRouter(db *database.Database) chi.Router {
	c := db.Commands
	q := db.Queries

	mux := chi.NewMux()
	mux.Use(auth.UseBearerToken())
	mux.Get("/schedule", handler.GetEventsSchedule(q.Events))
	mux.Route("/org", func(r chi.Router) {
		r.Get("/schedule", handler.GetOrgSchedule())

		r.Post("/", handler.CreateOrg(c.CreateOrganisation))
		r.Route("/{org}", func(r chi.Router) {
			r.Get("/", handler.GetOrg())
			r.Put("/", handler.UpdateOrg())
			r.Delete("/", handler.DeleteOrg())
		})
	})

	mux.Get("/events", handler.GetEvents(q.Events))
	mux.Route("/event", func(r chi.Router) {

		r.Post("/", handler.CreateEvent(c.CreateEvent))
		r.Route("/{eventID}", func(r chi.Router) {
			r.Get("/schedule", handler.GetSchedule(q.TimeslotsOfEvent))
			r.Get("/", handler.GetEvent(q.Event))
			r.Put("/", handler.UpdateEvent(c.UpdateEvent, q.UserHasRole))
			//r.Delete("/", DeleteEvent)

			r.Route("/timeslot", func(r chi.Router) {
				r.Post("/", handler.CreateTimeslot(c.CreateTimeslot))
				mux.Route("/{timeslotID}", func(r chi.Router) {
					r.Get("/", handler.GetTimeslot(q.Timeslot))
					r.Put("/", handler.UpdateTimeslot(c.UpdateTimeslot))
				})
			})
		})

		r.Route("/location/{locationID}", func(r chi.Router) {
			r.Put("/", handler.AddEventLocation())
			r.Delete("/", handler.RemoveEventLocation())
		})
	})

	mux.Route("/location", func(r chi.Router) {
		r.Post("/", handler.CreateLocation())
		r.Route("/{locationID}", func(r chi.Router) {
			r.Get("/", handler.GetLocation())
			r.Put("/", handler.UpdateLocation())
		})
	})

	mux.Route("/room", func(r chi.Router) {
		r.Get("/schedule", handler.GetRoomSchedule())

		r.Post("/", handler.CreateRoom())
		r.Route("/{roomID}", func(r chi.Router) {
			r.Get("/", handler.GetRoom())
			r.Put("/", handler.UpdateRoom())
			r.Delete("/", handler.DeleteRoom())
		})

	})

	return mux
}
