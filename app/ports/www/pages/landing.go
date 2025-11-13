package pages

import (
	"fmt"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/ports/www/components"
	"raumzeitalpaka/ports/www/middleware"
	"raumzeitalpaka/ports/www/render"
)

func LandingPage(events []model.EventModel) Node {
	g := Group{}
	for _, event := range events {
		g = append(g, Li(A(Href(fmt.Sprintf("/event/%v", event.ID)), Text(event.Name))))
	}
	return Shell("",
		Main(
			components.PageHeader(model.EventModel{}),
			components.AButton(components.ColorDefault, "/event/new", "New Event"),
			Ul(g),
		),
	)
}

type LandingPageRoute struct {
	DB *database.Database
}

func (l *LandingPageRoute) Method() string {
	return http.MethodGet
}

func (l *LandingPageRoute) Pattern() string {
	return "/"
}

func (l *LandingPageRoute) Handler() http.Handler {
	log := components.Logger(l)
	queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !middleware.IsOrganizer(request) {
			http.Redirect(writer, request, "/login", http.StatusTemporaryRedirect)
			return
		}

		events, err := queries.GetEvents(0, 100)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get events", err)
			return
		}

		render.HTML(log, writer, request, LandingPage(events))
	})
}
