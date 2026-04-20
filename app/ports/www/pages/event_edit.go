package pages

import (
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/ports/www/components"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func EditEventPage(event model.EventModel) Node {
	return components.Shell(event.Name,
		components.PageHeader(event),
		Main(
			H2(Text("Edit Event")),
			components.EventUpdateForm(event),
		),
	)
}
