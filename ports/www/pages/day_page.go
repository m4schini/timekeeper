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
	"timekeeper/ports/www/render"
)

func DayPage(day int, event model.EventModel, data []model.TimeslotModel) Node {
	return Shell(
		Main(
			components.PageHeader(event),
			components.FullDay(day+1, components.AddDays(event.Start, day), data),
			Div(Style("margin-top: 0.3rem; margin-bottom: -0.7rem"),
				Text("Export: "),
				components.ExportMarkdownButton(event.ID, day),
			),
			Script(Raw(`
document.getElementById('separator').scrollIntoView({
            behavior: 'auto',
            block: 'center',
            inline: 'center'
        });
`)),
			Script(Raw(`
console.log('starting reloader')
setInterval(() => {
    console.log('reloading page')
    window.location.reload()
    }, 60*1000)
`)),
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
			eventParam = strings.ToLower(chi.URLParam(request, "event"))
			dayParam   = strings.ToLower(chi.URLParam(request, "day"))
		)
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

		dayData := make([]model.TimeslotModel, 0, len(timeslots))
		for _, timeslot := range timeslots {
			if timeslot.Day == int(day) {
				dayData = append(dayData, timeslot)
			}
		}

		err = render.Render(writer, request, DayPage(int(day), event, dayData))
		if err != nil {
			log.Error("failed to render dayParam", zap.Error(err))
		}
	})
}
