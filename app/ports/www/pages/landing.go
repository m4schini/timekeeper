package pages

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/app/auth/authz"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/app/database/query"
	"raumzeitalpaka/ports/www/components"
	"raumzeitalpaka/ports/www/render"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func LandingPage(events []model.EventModel) Node {
	g := Group{}
	for _, event := range events {
		g = append(g, Li(A(Href(fmt.Sprintf("/event/%v", event.ID)), Text(event.Name))))
	}
	return components.Shell("",
		Main(
			components.PageHeader(model.EventModel{}),
			components.AButton(components.ColorDefault, "/event/new", "New Event"),
			Ul(g),
		),
	)
}

type LandingPageRoute struct {
	GetEvents query.GetEvents
	Authz     authz.Authorizer
}

func (l *LandingPageRoute) Method() string {
	return http.MethodGet
}

func (l *LandingPageRoute) Pattern() string {
	return "/"
}

func (l *LandingPageRoute) Handler() http.Handler {
	log := components.Logger(l)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		//if !middleware.IsOrganizer(request) {
		//	http.Redirect(writer, request, "/login", http.StatusTemporaryRedirect)
		//	return
		//}

		events, err := l.GetEvents.Query(query.GetEventsRequest{
			Offset: 0,
			Limit:  1000,
		})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get events", err)
			return
		}

		render.HTML(log, writer, request, LandingPage(events))
	})
}
