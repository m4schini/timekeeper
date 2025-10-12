package pages

import (
	"fmt"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/components"
	. "timekeeper/ports/www/render"
)

func LandingPage(events []model.EventModel) Node {
	g := Group{}
	for _, event := range events {
		g = append(g, Li(A(Href(fmt.Sprintf("/event/%v", event.ID)), Text(event.Name))))
	}
	return Shell(
		Main(
			components.PageHeader(model.EventModel{}, false),
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
	queries := l.DB.Queries
	log := zap.L().Named(l.Pattern())
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		events, err := queries.GetEvents(0, 100)
		if err != nil {
			RenderError(log, writer, http.StatusInternalServerError, "failed to get events", err)
			return
		}

		Render(log, writer, request, LandingPage(events))
	})
}
