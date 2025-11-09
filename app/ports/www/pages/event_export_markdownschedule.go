package pages

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	export "timekeeper/app/export/md"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/render"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Day struct {
	Index     int
	Timeslots []model.TimeslotModel
}

type EventScheduleExportMarkdownRoute struct {
	DB *database.Database
}

func (l *EventScheduleExportMarkdownRoute) Method() string {
	return http.MethodGet
}

func (l *EventScheduleExportMarkdownRoute) Pattern() string {
	return "/event/{event}/export/schedule.md"
}

func (l *EventScheduleExportMarkdownRoute) Handler() http.Handler {
	log := components.Logger(l)
	queries := l.DB.Queries
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

		event, err := queries.GetEvent(int(eventId))
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get event", err)
			return
		}

		timeslots, _, err := queries.GetTimeslotsOfEvent(int(eventId), roles, 0, 100)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to retrieve day", err)
			return
		}

		eventDays := model.MapTimeslotsToDays(timeslots)
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
