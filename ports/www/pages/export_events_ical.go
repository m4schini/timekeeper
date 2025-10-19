package pages

import (
	"net/http"
	"timekeeper/app/database"
	export "timekeeper/app/export/ical"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/render"
)

type EventsExportIcalRoute struct {
	DB *database.Database
}

func (v *EventsExportIcalRoute) Method() string {
	return http.MethodGet
}

func (v *EventsExportIcalRoute) Pattern() string {
	return "/export/events.ics"
}

func (v *EventsExportIcalRoute) Handler() http.Handler {
	log := components.Logger(v)
	queries := v.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		events, err := queries.GetEvents(0, 100)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get events", err)
			return
		}

		cal, err := export.ExportEventCalendar(events)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to generate ical schedule", err)
			return
		}

		writer.Header().Set("Content-Type", "text/calendar; charset=utf-8")
		writer.Header().Set("Content-Disposition", "inline; filename=timekeeper_events.ics")
		writer.Write([]byte(cal))
	})
}
