package pages

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/app/database/query"
	export "raumzeitalpaka/app/export/md"
	"raumzeitalpaka/ports/www/components"
	"raumzeitalpaka/ports/www/render"
	"slices"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Day struct {
	Index     int
	Timeslots []model.TimeslotModel
}

type EventScheduleExportMarkdownRoute struct {
	GetEvent            query.GetEvent
	GetTimeslotsOfEvent query.GetTimeslotsOfEvent
}

func (l *EventScheduleExportMarkdownRoute) Method() string {
	return http.MethodGet
}

func (l *EventScheduleExportMarkdownRoute) Pattern() string {
	return "/event/{event}/export/schedule.md"
}

func (l *EventScheduleExportMarkdownRoute) Handler() http.Handler {
	log := components.Logger(l)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var (
			eventParam = strings.ToLower(chi.URLParam(request, "event"))
		)
		log.Debug("export event markdown", zap.String("event", eventParam))
		eventId, err := strconv.ParseInt(eventParam, 10, 64)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid eventId", err)
			return
		}
		roles, _ := ParseRolesQuery(request.URL.Query(), false)

		event, err := l.GetEvent.Query(query.GetEventRequest{EventId: int(eventId)})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get event", err)
			return
		}

		timeslots, err := l.GetTimeslotsOfEvent.Query(query.GetTimeslotsOfEventRequest{
			EventId: int(eventId),
			Roles:   roles,
			Offset:  0,
			Limit:   1000,
		})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to retrieve day", err)
			return
		}

		eventDays := model.MapTimeslotsToDays(timeslots.Timeslots)
		days := make([]Day, 0, len(eventDays))
		for i, models := range eventDays {
			days = append(days, Day{
				Index:     i,
				Timeslots: models,
			})
		}
		slices.SortFunc(days, func(a, b Day) int {
			return a.Index - b.Index
		})

		var fullExport string
		for _, day := range days {
			table, err := export.ExportMarkdownTimetable(day.Timeslots)
			if err != nil {
				render.Error(log, writer, http.StatusInternalServerError, "failed to render markdown", err)
				return
			}

			fullExport += fmt.Sprintf("**%v**", export.Wochentag(event.Day(day.Index).Weekday()))
			fullExport += "\n"
			fullExport += table
			fullExport += "\n\n\n"
		}

		writer.Header().Set("Content-Type", "text/markdown; charset=utf8")
		writer.Write([]byte(fullExport))
	})
}
