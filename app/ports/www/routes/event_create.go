package routes

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/app/database/command"
	"raumzeitalpaka/app/database/query"
	"raumzeitalpaka/domain"
	"raumzeitalpaka/ports/www/components"
	"raumzeitalpaka/ports/www/pages"
	"raumzeitalpaka/ports/www/render"
	"time"

	"go.uber.org/zap"
)

func EventCreateFormHandler(getEvent query.GetEvent, createEvent command.CreateEvent) *MultiMethodRoute {
	pattern := "/event/new"
	log := zap.L().Named(pattern)
	return &MultiMethodRoute{
		Route: pattern,
		Get:   eventCreateFormPage(log),
		Post:  eventCreateFormHandler(log, createEvent),
	}
}

func eventCreateFormHandler(log *zap.Logger, createEvent command.CreateEvent) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		form, err := components.DecodeEventForm(r, false)
		if err != nil {
			render.Error(log, w, http.StatusBadRequest, "failed to decode form", err)
			return
		}

		id, err := domain.CreateEvent(ctx, createEvent, command.CreateEventRequest{
			Name:  form.Name,
			Slug:  form.Slug,
			Start: time.Time(form.Start),
			End:   time.Time(form.End),
		})
		if err != nil {
			render.Error(log, w, http.StatusInternalServerError, "failed to create event", err)
			return
		}
		log.Debug("created event", zap.Int("id", id))

		http.Redirect(w, r, fmt.Sprintf("/event/%v", id), http.StatusSeeOther)
	})
}

func eventCreateFormPage(log *zap.Logger) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		render.HTML(log, writer, request, pages.CreateEventPage())
	})
}
