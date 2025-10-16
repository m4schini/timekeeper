package pages

import (
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/middleware"
	. "timekeeper/ports/www/render"
)

func CreateUserPage() Node {
	return Shell(
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
	log := zap.L().Named("www").Named("user")
	//queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isOrganizer := middleware.IsOrganizer(request)
		if !isOrganizer {
			RenderError(log, writer, http.StatusUnauthorized, "user is not authorized", nil)
			return
		}

		Render(log, writer, request, CreateUserPage())
	})
}
