package pages

import (
	"net/http"
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/ports/www/components"
	"raumzeitalpaka/ports/www/middleware"
	"raumzeitalpaka/ports/www/render"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func CreateUserPage() Node {
	return components.Shell("",
		components.PageHeader(model.EventModel{}),
		Main(
			H2(Text("Create User")),
			components.UserForm(),
		),
	)
}

type CreateUserPageRoute struct {
	DB *database.Database
}

func (l *CreateUserPageRoute) Method() string {
	return http.MethodGet
}

func (l *CreateUserPageRoute) Pattern() string {
	return "/user/new"
}

func (l *CreateUserPageRoute) Handler() http.Handler {
	log := components.Logger(l)
	page := CreateUserPage()
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isOrganizer := middleware.IsOrganizer(request)
		if !isOrganizer {
			render.Error(log, writer, http.StatusUnauthorized, "user is not authorized", nil)
			return
		}

		render.HTML(log, writer, request, page)
	})
}
