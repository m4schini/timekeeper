package pages

import (
	"net/http"
	"raumzeitalpaka/app/auth/authz"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/app/database/query"
	"raumzeitalpaka/ports/www/components"
	"raumzeitalpaka/ports/www/middleware"
	"raumzeitalpaka/ports/www/render"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func CompactDayPage(event model.EventModel, data []model.TimeslotModel) Node {
	return components.ShellWithHead(event.Name, nil, []Node{},
		Main(
			components.CompactDay(data),
			components.ScriptReloadPageEveryMinute(),
		))
}

type EventScheduleDayRoute struct {
	GetEvent            query.GetEvent
	GetTimeslotsOfEvent query.GetTimeslotsOfEvent
	Authz               authz.Authorizer
}

func (l *EventScheduleDayRoute) Method() string {
	return http.MethodGet
}

func (l *EventScheduleDayRoute) Pattern() string {
	return "/event/{event}/schedule/{day}"
}

func (l *EventScheduleDayRoute) Handler() http.Handler {
	log := components.Logger(l)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var (
			eventParam  = strings.ToLower(chi.URLParam(request, "event"))
			dayParam    = strings.ToLower(chi.URLParam(request, "day"))
			isOrganizer = middleware.IsOrganizer(request, l.Authz)
			roles, _    = ParseRolesQuery(request.URL.Query(), isOrganizer)
		)

		eventId, err := strconv.ParseInt(eventParam, 10, 64)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid eventId", err)
			return
		}
		day, err := strconv.ParseInt(dayParam, 10, 64)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid day", err)
			return
		}

		event, err := l.GetEvent.Query(query.GetEventRequest{EventId: int(eventId)})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get event", err)
			return
		}

		timeslotsResponse, err := l.GetTimeslotsOfEvent.Query(query.GetTimeslotsOfEventRequest{
			EventId: int(eventId),
			Roles:   roles,
			Offset:  0,
			Limit:   1000,
		})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to retrieve day", err)
			return
		}
		timeslots := model.FilterTimeslotDay(timeslotsResponse.Timeslots, int(day))

		render.HTML(log, writer, request, CompactDayPage(event, timeslots))
	})
}
