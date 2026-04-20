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
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func EventUpdateFormHandler(getEvent query.GetEvent, userHasRole query.UserHasRole, updateEvent command.UpdateEvent) *MultiMethodRoute {
	pattern := "/event/{eventID}/edit"
	log := zap.L().Named(pattern)
	return &MultiMethodRoute{
		Route: pattern,
		Get:   eventUpdateFormPage(log, getEvent),
		Post:  eventUpdateFormHandler(log, userHasRole, updateEvent),
	}
}

func eventUpdateFormHandler(log *zap.Logger, userHasRole query.UserHasRole, updateEvent command.UpdateEvent) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		form, err := components.DecodeEventForm(r, true)
		if err != nil {
			render.Error(log, w, http.StatusBadRequest, "failed to decode form", err)
			return
		}

		err = domain.UpdateEvent(ctx, userHasRole, updateEvent, command.UpdateEventRequest{
			ID:    form.Event,
			Name:  form.Name,
			Slug:  form.Slug,
			Start: time.Time(form.Start),
			End:   time.Time(form.End),
		})
		if err != nil {
			render.Error(log, w, http.StatusInternalServerError, "failed to update event", err)
			return
		}
		log.Debug("updated event", zap.Int("id", form.Event))

		http.Redirect(w, r, fmt.Sprintf("/event/%v", form.Event), http.StatusSeeOther)
	})
}

func eventUpdateFormPage(log *zap.Logger, getEvent query.GetEvent) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		eventParam := chi.URLParam(request, "eventID")
		eventId, err := strconv.ParseInt(eventParam, 10, 64)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to parse eventID", err)
			return
		}

		event, err := getEvent.Query(request.Context(), query.GetEventRequest{EventId: int(eventId)})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get event", err)
			return
		}

		render.HTML(log, writer, request, pages.EditEventPage(event))
	})
}
