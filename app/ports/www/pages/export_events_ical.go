package pages

import (
	"net/http"
	"raumzeitalpaka/app/database/query"
	export "raumzeitalpaka/app/export/ical"
	"raumzeitalpaka/ports/www/components"
	"raumzeitalpaka/ports/www/render"
)

type EventsExportIcalRoute struct {
	GetEvents query.GetEvents
}

func (v *EventsExportIcalRoute) Method() string {
	return http.MethodGet
}

func (v *EventsExportIcalRoute) Pattern() string {
	return "/export/events.ics"
}

func (v *EventsExportIcalRoute) Handler() http.Handler {
	log := components.Logger(v)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		events, err := v.GetEvents.Query(query.GetEventsRequest{
			Offset: 0,
			Limit:  1000,
		})
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
		writer.Header().Set("Content-Disposition", "inline; filename=raumzeitalpaka_events.ics")
		writer.Write([]byte(cal))
	})
}
