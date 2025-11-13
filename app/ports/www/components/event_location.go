package components

import (
	"fmt"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"raumzeitalpaka/app/database/model"
)

func EventLocationCard(event model.EventModel, eventLocation model.EventLocationModel, withActions bool) Node {
	address := eventLocation.Address
	var eventRole = eventLocation.RelationshipLabel()

	var container func(children ...Node) Node
	if withActions {
		container = Div
	} else {
		container = A
	}

	return container(Class("location-card"), If(!withActions, Href(fmt.Sprintf("/event/%v/location/%v", event.ID, eventLocation.ID))),
		If(!withActions, Div(Strong(Text(eventRole)))),
		Iff(withActions, func() Node {
			return EventLocationUpdateForm(event.ID, eventLocation)
		}),
		If(!withActions, Div(Text(eventLocation.RelationshipNote))), Br(),
		Div(Style("white-space: pre-line"), Text(eventLocation.Name)),
		Div(Style("white-space: pre-line"), Iff(eventLocation.Address != nil, func() Node {
			return Textf(`%v %v
%v %v
`, address.Road, address.HouseNumber, address.Postcode, address.City)
		}),
			Iff(withActions, func() Node {
				return Div(Style("display: flex; gap: 1rem"),
					EditLocationButton(eventLocation.ID),
					DeleteEventLocationButton(event.ID, eventLocation.RelationshipId),
				)
			}),
		),
	)
}
