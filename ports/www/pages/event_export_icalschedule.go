package pages

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	export "timekeeper/app/export/ical"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/render"
)

type EventExportIcalScheduleRoute struct {
	DB *database.Database
}

func (v *EventExportIcalScheduleRoute) Method() string {
	return http.MethodGet
}

func (v *EventExportIcalScheduleRoute) Pattern() string {
	return "/event/{event}/export/schedule.ics"
}

func (v *EventExportIcalScheduleRoute) Handler() http.Handler {
	queries := v.DB.Queries
	log := components.Logger(v)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		eventParam := chi.URLParam(request, "event")
		eventId, err := strconv.ParseInt(eventParam, 10, 64)
		if err != nil {
			render.RenderError(log, writer, http.StatusBadRequest, "invalid event id", err)
			return
		}
		roles, _ := ParseRolesQuery(request.URL.Query(), false)

		event, err := queries.GetEvent(int(eventId))
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to get event", err)
			return
		}

		timeslots, _, err := queries.GetTimeslotsOfEvent(int(eventId), 0, 1000)
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to get timeslots of event", err)
			return
		}
		timeslots = model.FilterTimeslotRoles(timeslots, roles)

		writer.Header().Set("Content-Type", "text/calendar; charset=utf-8")
		writer.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=timekeeper_event_%v.ics", event.ID))

		cal, err := export.ExportCalendarSchedule(event, timeslots)
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to generate ical schedule", err)
			return
		}

		writer.Write([]byte(cal))
	})
}
