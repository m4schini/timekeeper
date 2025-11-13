package pages

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/ports/www/components"
	"raumzeitalpaka/ports/www/middleware"
	"raumzeitalpaka/ports/www/render"
)

func CreateLocationPage() Node {
	return Shell("",
		components.PageHeader(model.EventModel{}),
		Main(
			H2(Text("New Location")),
			components.CreateLocationForm(),
		),
	)
}

type CreateLocationPageRoute struct {
	DB *database.Database
}

func (l *CreateLocationPageRoute) Method() string {
	return http.MethodGet
}

func (l *CreateLocationPageRoute) Pattern() string {
	return "/location/new"
}

func (l *CreateLocationPageRoute) Handler() http.Handler {
	log := components.Logger(l)
	page := CreateLocationPage()
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isOrganizer := middleware.IsOrganizer(request)
		if !isOrganizer {
			render.Error(log, writer, http.StatusUnauthorized, "user is not authorized", nil)
			return
		}

		render.HTML(log, writer, request, page)
	})
}
