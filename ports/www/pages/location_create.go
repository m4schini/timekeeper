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

func CreateLocationPage() Node {
	return Shell(
		components.PageHeader(model.EventModel{}),
		Main(
			Div(Text("Create Location")),
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
	return "/location/create"
}

func (l *CreateLocationPageRoute) Handler() http.Handler {
	log := zap.L().Named("www").Named("event")
	//queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isOrganizer := middleware.IsOrganizer(request)
		roles, _ := ParseRolesQuery(request.URL.Query(), isOrganizer)

		_ = roles

		Render(log, writer, request, CreateLocationPage())
	})
}
