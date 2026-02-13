package pages

import (
	"net/http"
	"raumzeitalpaka/app/auth/authz"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/ports/www/components"
	"raumzeitalpaka/ports/www/middleware"
	"raumzeitalpaka/ports/www/render"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func CreateEventPage() Node {
	return components.Shell("",
		components.PageHeader(model.EventModel{}),
		Main(
			H2(Text("Create Event")),
			components.EventCreateForm(),
		),
	)
}

type CreateEventPageRoute struct {
	Authz authz.Authorizer
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
		isOrganizer := middleware.IsOrganizer(request, l.Authz)
		if !isOrganizer {
			render.Error(log, writer, http.StatusUnauthorized, "user is not authorized", nil)
			return
		}

		render.HTML(log, writer, request, page)
	})
}
