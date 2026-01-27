package components

import (
	"fmt"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/config"
	"time"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func PageHeader(event model.EventModel) Node {
	now := time.Now()
	time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), config.Timezone())

	return Header(Class("page-header"),
		Logo(event.Name, event.ID),
		Div(Class("last-change"),
			A(Href("/login"), Text("Login")),
			A(Href("/logout"), Text("Logout")),
		),
	)
}

func Logo(name string, eventId int) Node {
	if name == "" {
		name = "Raumzeitalpaka"
	}
	href := "/"
	if eventId != 0 {
		href = fmt.Sprintf("/event/%v", eventId)
	}
	return H1(Class("logo"),
		A(Text(name), Href(href)),
	)
}
