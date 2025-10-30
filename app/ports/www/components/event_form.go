package components

import (
	"fmt"
	"regexp"
	"timekeeper/app/database/model"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

const (
	EventSlugPattern = `(?i)^[a-z0-9-]+$` // (?i) makes it case-insensitive
)

var (
	EventSlugRegex = regexp.MustCompile(EventSlugPattern)
)

func eventForm(event *model.EventModel, method, action, actionText string) Node {
	hasModel := event != nil
	eventId := -1
	if hasModel {
		eventId = event.ID
	}
	return Form(Method(method), Action(action), Class("form"),
		Input(Type("hidden"), Name("event"), Value(fmt.Sprintf("%v", eventId))),

		Div(Class("param"),
			Label(For("name"), Text("Name")),
			Input(Type("text"), Name("name"), Placeholder("Jugend hackt 2042"), Required(), Iff(hasModel, func() Node {
				return Value(event.Name)
			})),
		),

		Div(Class("param"),
			Label(For("slug"), Text("URL Slug")),
			Input(Type("text"), Name("slug"), Placeholder("jh42"), Pattern("^[A-Za-z0-9-]+$"), Required(), Iff(hasModel, func() Node {
				return Value(event.Slug)
			})),
		),

		Div(Class("param"),
			Label(For("start"), Text("Erster Tag")),
			Input(Type("text"), Name("start"),
				Placeholder("13.12.2042"),
				Pattern(`^(0[1-9]|[12][0-9]|3[01])\.(0[1-9]|1[0-2])\.(\d{4})$`),
				Required(), Iff(hasModel, func() Node {
					return Value(event.Start.Format("02.01.2006"))
				})),
		),

		Input(Type("submit"), Value(actionText)),
	)
}
