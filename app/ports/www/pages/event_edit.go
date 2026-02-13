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

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func EditEventPage(event model.EventModel) Node {
	return components.Shell(event.Name,
		components.PageHeader(event),
		Main(
			H2(Text("Edit Event")),
			components.EventUpdateForm(event),
		),
	)
}

type EditEventPageRoute struct {
	GetEvent query.GetEvent
	Authz    authz.Authorizer
}

func (l *EditEventPageRoute) Method() string {
	return http.MethodGet
}

func (l *EditEventPageRoute) Pattern() string {
	return "/event/{event}/edit"
}

func (l *EditEventPageRoute) Handler() http.Handler {
	log := components.Logger(l)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isOrganizer := middleware.IsOrganizer(request, l.Authz)
		if !isOrganizer {
			render.Error(log, writer, http.StatusUnauthorized, "user is not authorized", nil)
			return
		}

		var (
			eventParam   = chi.URLParam(request, "event")
			eventId, err = strconv.ParseInt(eventParam, 10, 64)
		)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid eventId", err)
			return
		}

		event, err := l.GetEvent.Query(query.GetEventRequest{EventId: int(eventId)})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get event", err)
			return
		}

		render.HTML(log, writer, request, EditEventPage(event))
	})
}
