package pages

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/middleware"
	. "timekeeper/ports/www/render"
)

func EventPage(event model.EventModel, withActions bool, days [][]model.TimeslotModel) Node {
	dayNodes := make(Group, len(days))
	for i, timeslots := range days {
		dayNodes[i] = components.Day(event.ID, i, event.Day(i), withActions, timeslots)
	}

	return Shell(
		Main(
			components.PageHeader(event, withActions),
			If(withActions, components.EventActions(event.ID)),
			Div(Class("days-container"), dayNodes),
		),
		components.ScriptScrollSeperatorIntoView(),
		If(!withActions, components.ScriptReloadPageEveryMinute()),
	)
}

func ParseRolesQuery(query url.Values, userIsOrganizer bool) (roles []model.Role, hasRoles bool) {
	hasRoles = query.Has("role")
	rolesStrs := strings.Split(query.Get("role"), ",")

	if !hasRoles {
		if userIsOrganizer {
			rolesStrs = []string{string(model.RoleOrganizer), string(model.RoleMentor), string(model.RoleParticipant)}
		} else {
			rolesStrs = []string{string(model.RoleParticipant)}
		}
	}

	roles = make([]model.Role, len(rolesStrs))
	for i, role := range rolesStrs {
		roles[i] = model.RoleFrom(role)
	}

	return roles, hasRoles
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

func (l *EventPageRoute) UseCache() bool {
	return true
}

func (l *EventPageRoute) Handler() http.Handler {
	log := zap.L().Named("www").Named("event")
	queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		eventParam := chi.URLParam(request, "event")
		eventId, err := strconv.ParseInt(eventParam, 10, 64)
		if err != nil {
			RenderError(log, writer, http.StatusBadRequest, "invalid event id", err)
			return
		}
		isOrganizer := middleware.IsOrganizer(request)
		roles, _ := ParseRolesQuery(request.URL.Query(), isOrganizer)

		event, err := queries.GetEvent(int(eventId))
		if err != nil {
			RenderError(log, writer, http.StatusInternalServerError, "failed to get event", err)
			return
		}

		timeslots, _, err := queries.GetTimeslotsOfEvent(int(eventId), 0, 1000)
		if err != nil {
			RenderError(log, writer, http.StatusInternalServerError, "failed to get timeslots", err)
			return
		}

		eventDays := model.MapTimeslotsToDays(timeslots)
		renderData := make([][]model.TimeslotModel, len(eventDays))
		for day, timeslotsOfDay := range eventDays {
			renderData[day] = model.FilterTimeslotRoles(timeslotsOfDay, roles)
		}

		Render(log, writer, request, EventPage(event, isOrganizer, renderData))
	})
}
