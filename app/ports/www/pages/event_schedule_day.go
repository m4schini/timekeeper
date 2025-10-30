package pages

import (
	"net/http"
	"strconv"
	"strings"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/middleware"
	"timekeeper/ports/www/render"

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func CompactDayPage(event model.EventModel, data []model.TimeslotModel) Node {
	return ShellWithHead(event.Name, nil, []Node{},
		Main(
			components.CompactDay(data),
			components.ScriptReloadPageEveryMinute(),
		))
}

type EventScheduleDayRoute struct {
	DB *database.Database
}

func (l *EventScheduleDayRoute) Method() string {
	return http.MethodGet
}

func (l *EventScheduleDayRoute) Pattern() string {
	return "/event/{event}/schedule/{day}"
}

func (l *EventScheduleDayRoute) Handler() http.Handler {
	log := components.Logger(l)
	queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var (
			eventParam  = strings.ToLower(chi.URLParam(request, "event"))
			dayParam    = strings.ToLower(chi.URLParam(request, "day"))
			isOrganizer = middleware.IsOrganizer(request)
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

		event, err := queries.GetEvent(int(eventId))
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get event", err)
			return
		}

		timeslots, _, err := queries.GetTimeslotsOfEvent(int(eventId), 0, 100)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to retrieve day", err)
			return
		}
		timeslots = model.FilterTimeslotDay(timeslots, int(day))
		timeslots = model.FilterTimeslotRoles(timeslots, roles)

		render.HTML(log, writer, request, CompactDayPage(event, timeslots))
	})
}
