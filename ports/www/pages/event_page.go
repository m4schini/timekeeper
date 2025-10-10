package pages

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"strconv"
	"strings"
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
			components.PageHeader(event, withActions),
			Div(Class("days-container"),
				g,
			),
		),
		components.ScriptScrollSeperatorIntoView(),
		If(!withActions, components.ScriptReloadPageEveryMinute()),
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
		roles := strings.Split(request.URL.Query().Get("role"), ",")
		isOrganizer := middleware.IsOrganizer(request)

		if !hasFilter {
			if isOrganizer {
				roles = []string{string(model.RoleOrganizer), string(model.RoleMentor), string(model.RoleParticipant)}
			} else {
				roles = []string{string(model.RoleParticipant)}
			}
		}

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

		filterRoles := make([]model.Role, len(roles))
		for i, role := range roles {
			filterRoles[i] = model.RoleFrom(role)
		}

		renderData := make([][]model.TimeslotModel, len(tsMap))
		for i, models := range tsMap {
			var filteredModels []model.TimeslotModel
			for _, timeslotModel := range models {
				for _, role := range filterRoles {
					if timeslotModel.Role == role {
						filteredModels = append(filteredModels, timeslotModel)
						break
					}
				}
			}

			renderData[i] = filteredModels
		}

		Render(writer, request, EventPage(event, isOrganizer, renderData))
	})
}
