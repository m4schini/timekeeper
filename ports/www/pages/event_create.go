package pages

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/middleware"
	. "timekeeper/ports/www/render"
)

func CreateEventPage() Node {
	return Shell(
		components.PageHeader(model.EventModel{}),
		Main(
			Div(Text("Create Event")),
			components.EventForm(nil, "POST", "/_/event", "Create"),
		),
	)
}

type CreateEventPageRoute struct {
	DB *database.Database
}

func (l *CreateEventPageRoute) Method() string {
	return http.MethodGet
}

func (l *CreateEventPageRoute) Pattern() string {
	return "/event/new"
}

func (l *CreateEventPageRoute) Handler() http.Handler {
	log := components.Logger(l)
	//queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isOrganizer := middleware.IsOrganizer(request)
		if !isOrganizer {
			RenderError(log, writer, http.StatusUnauthorized, "user is not authorized", nil)
			return
		}

		Render(log, writer, request, CreateEventPage())
	})
}
