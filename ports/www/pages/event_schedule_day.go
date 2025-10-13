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
	"timekeeper/app/database/model"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/middleware"
	"timekeeper/ports/www/render"
)

func DayPage(day int, event model.EventModel, data []model.TimeslotModel) Node {
	return Shell(
		Main(
			components.PageHeader(event),
			components.FullDay(day+1, event.Day(day), data),
			components.ScriptScrollSeperatorIntoView(),
			components.ScriptReloadPageEveryMinute(),
		))
}

func CompactDayPage(data []model.TimeslotModel) Node {
	return Shell(
		Main(
			components.CompactDay(data),
			components.ScriptScrollSeperatorIntoView(),
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
	log := zap.L().Named(l.Pattern())
	queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var (
			eventParam  = strings.ToLower(chi.URLParam(request, "event"))
			dayParam    = strings.ToLower(chi.URLParam(request, "day"))
			isOrganizer = middleware.IsOrganizer(request)
			roles, _    = ParseRolesQuery(request.URL.Query(), isOrganizer)
			useCompact  = request.URL.Query().Has("compact")
		)

		eventId, err := strconv.ParseInt(eventParam, 10, 64)
		if err != nil {
			render.RenderError(log, writer, http.StatusBadRequest, "invalid eventId", err)
			return
		}
		day, err := strconv.ParseInt(dayParam, 10, 64)
		if err != nil {
			render.RenderError(log, writer, http.StatusBadRequest, "invalid day", err)
			return
		}

		event, err := queries.GetEvent(int(eventId))
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to get event", err)
			return
		}

		timeslots, _, err := queries.GetTimeslotsOfEvent(int(eventId), 0, 100)
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to retrieve day", err)
			return
		}
		timeslots = model.FilterTimeslotDay(timeslots, int(day))
		timeslots = model.FilterTimeslotRoles(timeslots, roles)

		var page Node
		if useCompact {
			page = CompactDayPage(timeslots)
		} else {
			page = DayPage(int(day), event, timeslots)
		}

		render.Render(log, writer, request, page)
	})
}
