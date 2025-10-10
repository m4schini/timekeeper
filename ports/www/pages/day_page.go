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
			components.PageHeader(event, false),
			components.FullDay(day+1, components.AddDays(event.Start, day), data),
			Div(Style("margin-top: 0.3rem; margin-bottom: -0.7rem"),
				Text("Export: "),
				components.ExportMarkdownButton(event.ID, day),
			),
			components.ScriptScrollSeperatorIntoView(),
			components.ScriptReloadPageEveryMinute(),
		))
}

func CompactDayPage(day int, event model.EventModel, data []model.TimeslotModel) Node {
	return Shell(
		Main(
			components.CompactDay(event.ID, day, data),
			//Div(Style("margin-top: 0.3rem; margin-bottom: -0.7rem"),
			//	Text("Export: "),
			//	components.ExportMarkdownButton(event.ID, day),
			//),
			components.ScriptScrollSeperatorIntoView(),
			components.ScriptReloadPageEveryMinute(),
		))
}

type DayPageRoute struct {
	DB *database.Database
}

func (l *DayPageRoute) Method() string {
	return http.MethodGet
}

func (l *DayPageRoute) Pattern() string {
	return "/event/{event}/{day}"
}

func (l *DayPageRoute) Handler() http.Handler {
	log := zap.L().Named(l.Pattern())
	queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var (
			eventParam  = strings.ToLower(chi.URLParam(request, "event"))
			dayParam    = strings.ToLower(chi.URLParam(request, "day"))
			hasFilter   = request.URL.Query().Has("role")
			roles       = strings.Split(request.URL.Query().Get("role"), ",")
			isOrganizer = middleware.IsOrganizer(request)
			useCompact  = request.URL.Query().Has("compact")
		)

		if !hasFilter {
			if isOrganizer {
				roles = []string{string(model.RoleOrganizer), string(model.RoleMentor), string(model.RoleParticipant)}
			} else {
				roles = []string{string(model.RoleParticipant)}
			}
		}

		log.Debug("rendering day", zap.String("eventParam", eventParam), zap.String("dayParam", dayParam))
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

		filterRoles := make([]model.Role, len(roles))
		for i, role := range roles {
			filterRoles[i] = model.RoleFrom(role)
		}

		dayData := make([]model.TimeslotModel, 0, len(timeslots))
		for _, timeslot := range timeslots {
			if timeslot.Day == int(day) {
				for _, role := range filterRoles {
					if timeslot.Role == role {
						dayData = append(dayData, timeslot)
						break
					}
				}
			}
		}

		var page Node
		if useCompact {
			page = CompactDayPage(int(day), event, dayData)
		} else {
			page = DayPage(int(day), event, dayData)
		}

		err = render.Render(writer, request, page)
		if err != nil {
			log.Error("failed to render dayParam", zap.Error(err))
		}
	})
}
