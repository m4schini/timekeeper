package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"time"
	"timekeeper/app/database/model"
)

func PageHeader(event model.EventModel, withActions bool) Node {
	return Header(Class("page-header"),
		Logo(event.Name, event.ID),
		Div(Class("last-change"),
			Text("Generated:"),
			Br(),
			Text(time.Now().Format(time.RFC822)),
		),
	)
}
