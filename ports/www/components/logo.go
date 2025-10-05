package components

import (
	"fmt"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Logo(name string, eventId int) Node {
	if name == "" {
		name = "Timekeeper"
	}
	href := "/"
	if eventId != 0 {
		href = fmt.Sprintf("/event/%v", eventId)
	}
	return H1(Class("logo"),
		A(Text(name), Href(href)),
	)
}
