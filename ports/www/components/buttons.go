package components

import (
	"fmt"
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

func ExportMarkdownButton(eventId, day int) Node {
	return A(Href(fmt.Sprintf("/event/%v/%v/export/markdown", eventId, day)), Text("Markdown"))
}

func CreateTimeslotButton(eventId int) Node {
	return A(Href(fmt.Sprintf("/create/timeslot?event=%v", eventId)), Text("Create Timeslot"))
}

func EditTimeslotButton(timeslotId int) Node {
	return A(Text("edit"), Href(fmt.Sprintf("/edit/timeslot/%v", timeslotId)))
}

func DuplicateTimeslotButton(timeslotId int) Node {
	return A(Text("duplicate"), Href(fmt.Sprintf("/duplicate/timeslot/%v", timeslotId)))
}

func DeleteTimeslotButton(timeslotId int) Node {
	return A(Text("delete"), Href("#"),
		hx.Delete(fmt.Sprintf("/_/timeslot/%v", timeslotId)),
		hx.Target("closest .timeslot-container"),
		hx.Swap("outerHTML swap:1s"),
	)
}
