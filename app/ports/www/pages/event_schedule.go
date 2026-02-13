package pages

import (
	"fmt"
	"net/http"
	"net/url"
	"raumzeitalpaka/app/auth/authz"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/app/database/query"
	"raumzeitalpaka/app/export/md"
	"raumzeitalpaka/ports/www/components"
	"raumzeitalpaka/ports/www/middleware"
	"raumzeitalpaka/ports/www/render"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func SchedulePage(event model.EventModel, withActions bool, days [][]model.TimeslotModel, roles []model.Role) Node {
	dayNodes := make(Group, len(days))
	rolesStr := make([]string, len(roles))
	for i, r := range roles {
		rolesStr[i] = string(r)
	}

	for i, timeslots := range days {
		dayNodes[i] = components.Day(event.ID, i, event.Day(i), withActions, timeslots, fmt.Sprintf("/event/%v/schedule/%v?role=%v", event.ID, i, strings.Join(rolesStr, ",")))
	}

	return components.Shell(event.Name,
		Main(
			components.PageHeader(event),
			If(withActions, components.EventActions(event.ID)),
			Div(Class("days-container"), dayNodes),
		),
		components.ScriptScrollSeperatorIntoView(),
		If(!withActions, components.ScriptReloadPageEveryMinute()),
	)
}

func CompactSchedulePage(event model.EventModel, days [][]model.TimeslotModel) Node {
	dayNodes := make(Group, 0, len(days)*2)
	for i, timeslots := range days {
		h := H3(Text(fmt.Sprintf("Tag %v (%v)", i, md.Wochentag(event.Day(i).Weekday()))))

		dayNodes = append(dayNodes, h)
		dayNodes = append(dayNodes, components.CompactDay(timeslots))
	}

	return components.ShellWithHead(event.Name, nil, []Node{},
		Main(
			dayNodes,
		),
		components.ScriptScrollSeperatorIntoView(),
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

type SchedulePageRoute struct {
	GetEvent            query.GetEvent
	GetTimeslotsOfEvent query.GetTimeslotsOfEvent
	Authz               authz.Authorizer
}

func (l *SchedulePageRoute) Method() string {
	return http.MethodGet
}

func (l *SchedulePageRoute) Pattern() string {
	return "/event/{event}/schedule"
}

func (l *SchedulePageRoute) Handler() http.Handler {
	log := components.Logger(l)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		eventParam := chi.URLParam(request, "event")
		eventId, err := strconv.ParseInt(eventParam, 10, 64)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid event id", err)
			return
		}
		useCompact := request.URL.Query().Has("compact")
		isOrganizer := middleware.IsOrganizer(request, l.Authz)
		roles, _ := ParseRolesQuery(request.URL.Query(), isOrganizer)

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
			render.Error(log, writer, http.StatusInternalServerError, "failed to get timeslots", err)
			return
		}

		eventDays := model.MapTimeslotsToDays(timeslots.Timeslots)
		renderData := make([][]model.TimeslotModel, len(eventDays))
		for day, timeslotsOfDay := range eventDays {
			renderData[day] = timeslotsOfDay
		}

		if useCompact {
			render.HTML(log, writer, request, CompactSchedulePage(event, renderData))
		} else {
			render.HTML(log, writer, request, SchedulePage(event, isOrganizer, renderData, roles))
		}
	})
}
