package components

import (
	"fmt"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"strings"
	"timekeeper/app/database/model"
)

const locationMapSrc = "/static/betahaus.png"

func EventLocationCard(event model.EventModel, eventLocation model.EventLocationModel, withActions bool) Node {
	address := eventLocation.Address
	var eventRole string
	switch strings.TrimSpace(strings.ToLower(eventLocation.Relationship)) {
	case "sleep_location":
		eventRole = "Übernachtung"
		break
	case "event_location":
		eventRole = "Event Location"
		break
	}

	return Div(
		Div(Strong(Text(eventRole))), Br(),
		Div(Text(eventLocation.Name)),
		Div(Style("white-space: pre-line"), Iff(eventLocation.Address != nil, func() Node {
			return Textf(`%v %v
%v %v
`, address.Road, address.HouseNumber, address.Postcode, address.City)
		}),
			Iff(withActions, func() Node {
				return Div(Style("display: flex; gap: 1rem"),
					EditLocationButton(eventLocation.RelationshipId),
					DeleteEventLocationButton(event.ID, eventLocation.RelationshipId),
				)
			}),
		),
	)
}

func LocationMap() Node {
	return Div(Class("location-map"),
		Img(Src(locationMapSrc)),
	)
}

func LocationCrop(x, y, width, height, targetH int) Node {
	scale := float64(targetH) / float64(height)
	targetW := int(scale * float64(width))

	return Div(
		Style(fmt.Sprintf("width:%dpx; height:%dpx; overflow:hidden;", targetW, targetH)),
		Div(
			Style(fmt.Sprintf(`
				width:%dpx;
				height:%dpx;
				background-image:url('%s');
				background-position:-%dpx -%dpx;
				background-repeat:no-repeat;
				transform:scale(%f);
				transform-origin:top left;
			`,
				width, height, locationMapSrc, x, y, scale,
			)),
		),
	)

}
