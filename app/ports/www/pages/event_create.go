package pages

import (
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/ports/www/components"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func CreateEventPage() Node {
	return components.Shell("",
		components.PageHeader(model.EventModel{}),
		Main(
			H2(Text("Create Event")),
			components.EventCreateForm(),
		),
	)
}
