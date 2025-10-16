package pages

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"strconv"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/middleware"
	. "timekeeper/ports/www/render"
)

func EditEventPage(event model.EventModel) Node {
	return Shell(
		components.PageHeader(model.EventModel{}),
		Main(
			Div(Text("Edit Event")),
			components.EventForm(&event, "POST", "/_/event/edit", "Update"),
		),
	)
}

type EditEventPageRoute struct {
	DB *database.Database
}

func (l *EditEventPageRoute) Method() string {
	return http.MethodGet
}

func (l *EditEventPageRoute) Pattern() string {
	return "/event/{event}/edit"
}

func (l *EditEventPageRoute) Handler() http.Handler {
	log := zap.L().Named("www").Named("event")
	queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isOrganizer := middleware.IsOrganizer(request)
		if !isOrganizer {
			RenderError(log, writer, http.StatusUnauthorized, "user is not authorized", nil)
			return
		}

		var (
			eventParam   = chi.URLParam(request, "event")
			eventId, err = strconv.ParseInt(eventParam, 10, 64)
		)
		if err != nil {
			RenderError(log, writer, http.StatusBadRequest, "invalid eventId", err)
			return
		}

		event, err := queries.GetEvent(int(eventId))
		if err != nil {
			RenderError(log, writer, http.StatusInternalServerError, "failed to get event", err)
			return
		}

		Render(log, writer, request, EditEventPage(event))
	})
}
