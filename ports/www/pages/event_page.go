package pages

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"strconv"
	"timekeeper/app/database"
	model "timekeeper/app/database/model"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/middleware"
	. "timekeeper/ports/www/render"
)

func EventPage(event model.EventModel, withActions bool, days [][]model.TimeslotModel) Node {
	g := make(Group, len(days))
	for i, timeslots := range days {
		g[i] = components.Day(event.ID, i, components.AddDays(event.Start, i), withActions, timeslots)
	}

	return Shell(
		Main(
			components.PageHeader(event),
			Div(Class("days-container"),
				g,
			),
		),
		Script(Raw(`
document.getElementById('separator').scrollIntoView({
            behavior: 'auto',
            block: 'center',
            inline: 'center'
        });
`)),
	)
}

type EventPageRoute struct {
	DB *database.Database
}

func (l *EventPageRoute) Method() string {
	return http.MethodGet
}

func (l *EventPageRoute) Pattern() string {
	return "/event/{event}"
}

func (l *EventPageRoute) Handler() http.Handler {
	log := zap.L().Named("www").Named("event")
	queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		eventParam := chi.URLParam(request, "event")
		eventId, err := strconv.ParseInt(eventParam, 10, 64)
		hasFilter := request.URL.Query().Has("role")
		filterRole := model.RoleFrom(request.URL.Query().Get("role"))
		isOrganizer := middleware.IsOrganizer(request)

		if err != nil {
			RenderError(log, writer, http.StatusBadRequest, "invalid event id", err)
			return
		}

		event, err := queries.GetEvent(int(eventId))
		if err != nil {
			RenderError(log, writer, http.StatusInternalServerError, "failed to get event", err)
			return
		}

		timeslots, total, err := queries.GetTimeslotsOfEvent(int(eventId), 0, 1000)
		if err != nil {
			RenderError(log, writer, http.StatusInternalServerError, "failed to get timeslots", err)
			return
		}
		log.Debug("retrieved timeslots", zap.Int("total", total), zap.Int("count", len(timeslots)))

		tsMap := make(map[int][]model.TimeslotModel)
		for _, timeslot := range timeslots {
			ts, _ := tsMap[timeslot.Day]
			if ts == nil {
				ts = make([]model.TimeslotModel, 0)
			}

			ts = append(ts, timeslot)

			tsMap[timeslot.Day] = ts
		}

		renderData := make([][]model.TimeslotModel, len(tsMap))
		for i, models := range tsMap {
			var filteredModels []model.TimeslotModel
			if hasFilter {
				for _, timeslotModel := range models {
					if timeslotModel.Role == filterRole {
						filteredModels = append(filteredModels, timeslotModel)
					}
				}
			} else {
				filteredModels = models
			}

			renderData[i] = filteredModels
		}

		Render(writer, request, EventPage(event, isOrganizer, renderData))
	})
}
