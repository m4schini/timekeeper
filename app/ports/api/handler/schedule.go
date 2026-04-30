package handler

import (
	"net/http"
	"raumzeitalpaka/app/database/query"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetSchedule(schedule query.GetTimeslotsOfEvent) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		eventParam := chi.URLParam(request, "eventID")
		eventID, err := strconv.ParseInt(eventParam, 10, 64)
		if err != nil {
			http.Error(writer, "invalid eventID", http.StatusBadRequest)
			return
		}

		roles, _ := ParseRolesQuery(request.URL.Query(), false)

		resp, err := schedule.Query(ctx, query.GetTimeslotsOfEventRequest{
			EventId: int(eventID),
			Roles:   roles,
			Offset:  0,
			Limit:   1000,
		})
		if err != nil {
			http.Error(writer, "failed to get data", http.StatusInternalServerError)
			return
		}

		Encode(writer, resp)
	}

}

func GetRoomSchedule() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func GetOrgSchedule() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func GetEventsSchedule(events query.GetEvents) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		resp, err := events.Query(ctx, query.GetEventsRequest{
			Offset: 0,
			Limit:  1000,
		})
		if err != nil {
			http.Error(w, "failed to get data", http.StatusInternalServerError)
			return
		}

		Encode(w, resp)
	}
}
