package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"time"
	"timekeeper/app/database/model"
	"timekeeper/config"
)

func PageHeader(event model.EventModel) Node {
	now := time.Now()
	time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), config.Timezone())

	return Header(Class("page-header"),
		Logo(event.Name, event.ID),
		Div(Class("last-change"),
			Text("Generated:"),
			Br(),
			Text(now.Format(time.RFC822)),
		),
	)
}
