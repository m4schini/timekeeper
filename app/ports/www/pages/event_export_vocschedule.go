package pages

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	export "timekeeper/app/export/voc"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/render"
)

type EventExportVocScheduleRoute struct {
	DB *database.Database
}

func (v *EventExportVocScheduleRoute) Method() string {
	return http.MethodGet
}

func (v *EventExportVocScheduleRoute) Pattern() string {
	return "/event/{event}/export/schedule.json"
}

func (v *EventExportVocScheduleRoute) Handler() http.Handler {
	log := components.Logger(v)
	queries := v.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		eventParam := chi.URLParam(request, "event")
		eventId, err := strconv.ParseInt(eventParam, 10, 64)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid event id", err)
			return
		}
		roles, _ := ParseRolesQuery(request.URL.Query(), false)

		event, err := queries.GetEvent(int(eventId))
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get event", err)
			return
		}

		timeslots, _, err := queries.GetTimeslotsOfEvent(int(eventId), 0, 1000)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get timeslots of event", err)
			return
		}
		timeslots = model.FilterTimeslotRoles(timeslots, roles)

		writer.Header().Set("Content-Type", "application/json")
		err = export.ExportVocScheduleTo(event, timeslots, writer)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to generate voc schedule", err)
			return
		}
	})
}
