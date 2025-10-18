package pages

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/middleware"
	"timekeeper/ports/www/render"
)

func CreateEventPage() Node {
	return Shell("",
		components.PageHeader(model.EventModel{}),
		Main(
			H2(Text("Create Event")),
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
	page := CreateEventPage()
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isOrganizer := middleware.IsOrganizer(request)
		if !isOrganizer {
			render.Error(log, writer, http.StatusUnauthorized, "user is not authorized", nil)
			return
		}

		render.HTML(log, writer, request, page)
	})
}
