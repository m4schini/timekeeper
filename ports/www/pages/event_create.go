package pages

import (
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"timekeeper/app/database"
	"timekeeper/ports/www/middleware"
	. "timekeeper/ports/www/render"
)

func CreateEventPage() Node {
	return Shell(
		Main(
			Div(Text("Create Event")),
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
	return "/event/create"
}

func (l *CreateEventPageRoute) UseCache() bool {
	return false
}

func (l *CreateEventPageRoute) Handler() http.Handler {
	log := zap.L().Named("www").Named("event")
	//queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isOrganizer := middleware.IsOrganizer(request)
		roles, _ := ParseRolesQuery(request.URL.Query(), isOrganizer)

		_ = roles

		Render(log, writer, request, CreateEventPage())
	})
}
