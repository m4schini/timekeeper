package handler

import (
	"net/http"
	"raumzeitalpaka/app/database/command"
	"raumzeitalpaka/app/database/query"
	"raumzeitalpaka/domain"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func CreateEvent(createEvent command.CreateEvent) http.HandlerFunc {
	type Response struct {
		Event int `json:"event"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := requireAuthentication(ctx)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		request, err := Decode[command.CreateEventRequest](r)
		if err != nil {
			http.Error(w, "failed to decode request", http.StatusBadRequest)
			return
		}

		eventId, err := domain.CreateEvent(ctx, createEvent, request)
		if err != nil {
			InternalServerErr(w, err)
			return
		}

		Encode(w, Response{Event: eventId})
	}
}

func GetEvent(getEvent query.GetEvent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		eventParam := chi.URLParam(r, "eventID")
		eventId, err := strconv.ParseInt(eventParam, 10, 64)
		if err != nil {
			http.Error(w, "failed to decode request", http.StatusBadRequest)
			return
		}

		event, err := domain.GetEvent(ctx, getEvent, query.GetEventRequest{
			EventId: int(eventId),
		})
		if err != nil {
			InternalServerErr(w, err)
			return
		}

		Encode(w, event)
	}
}

func GetEvents(getEvents query.GetEvents) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		events, err := domain.GetEvents(ctx, getEvents)
		if err != nil {
			InternalServerErr(w, err)
			return
		}

		Encode(w, events)
	}
}

func UpdateEvent(updateEvent command.UpdateEvent, userHasRole query.UserHasRole) http.HandlerFunc {
	type Response struct {
		Status string `json:"status"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := requireAuthentication(ctx)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		request, err := Decode[command.UpdateEventRequest](r)
		if err != nil {
			http.Error(w, "failed to decode request", http.StatusBadRequest)
			return
		}

		err = domain.UpdateEvent(ctx, userHasRole, updateEvent, request)
		if err != nil {
			InternalServerErr(w, err)
			return
		}

		Encode(w, Response{Status: "ok"})
	}
}

func AddEventLocation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := requireAuthentication(ctx)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}
}

func RemoveEventLocation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := requireAuthentication(ctx)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}
}
