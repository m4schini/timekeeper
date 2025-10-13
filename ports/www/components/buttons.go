package components

import (
	"fmt"
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

func ExportEventDayMarkdownButton(eventId, day int) Node {
	return A(Href(fmt.Sprintf("/event/%v/%v/export/schedule.md", eventId, day)), Text("Markdown"))
}

func EventActions(eventId int) Node {
	return Div(Class("menu"),
		Text("Export: "),
		ExportEventMarkdownButton(eventId),
		ExportEventVocScheduleButton(eventId),
		Text("View: "),
		EventViewOrganizer(eventId),
		EventViewMentor(eventId),
		EventViewParticipant(eventId, true),
		Text("Share: "),
		EventViewOrganizer(eventId),
		EventViewMentor(eventId),
		EventViewParticipant(eventId, false),
	)
}

func ExportEventActions(eventId int) Node {
	return Div(Class("menu"), Text("Export: "),
		ExportEventMarkdownButton(eventId),
		ExportEventVocScheduleButton(eventId),
	)
}

func EventViewActions(eventId int) Node {
	return Div(Class("menu"), Text("View: "),
		EventViewOrganizer(eventId),
		EventViewMentor(eventId),
		EventViewParticipant(eventId, true),
	)
}

func ShareActions(eventId int) Node {
	return Div(Class("menu"), Text("Share: "),
		EventViewOrganizer(eventId),
		EventViewMentor(eventId),
		EventViewParticipant(eventId, false),
	)
}

func EventViewParticipant(eventId int, withFilter bool) Node {
	if withFilter {
		return A(Href(fmt.Sprintf("/event/%v?role=Participant", eventId)), Text("Teilnehmer*innen"))
	} else {
		return A(Href(fmt.Sprintf("/event/%v", eventId)), Text("Teilnehmer*innen"))
	}
}

func EventViewMentor(eventId int) Node {
	return A(Href(fmt.Sprintf("/event/%v?role=Mentor,Participant", eventId)), Text("Mentor*innen"))
}

func EventViewOrganizer(eventId int) Node {
	return A(Href(fmt.Sprintf("/event/%v?role=Organizer,Mentor,Participant", eventId)), Text("Orga"))
}

func ExportEventMarkdownButton(eventId int) Node {
	return A(Href(fmt.Sprintf("/event/%v/export/schedule.md", eventId)), Text("Markdown"))
}

func ExportEventVocScheduleButton(eventId int) Node {
	return A(Href(fmt.Sprintf("/event/%v/export/schedule.json", eventId)), Text("VOC Schedule"))
}

func CreateTimeslotButton(eventId int) Node {
	return A(Href(fmt.Sprintf("/timeslot/create?event=%v", eventId)), Text("Create Timeslot"))
}

func EditTimeslotButton(timeslotId int) Node {
	return A(Text("edit"), Href(fmt.Sprintf("/timeslot/edit/%v", timeslotId)))
}

func DuplicateTimeslotButton(timeslotId int) Node {
	return A(Text("duplicate"), Href(fmt.Sprintf("/timeslot/duplicate/%v", timeslotId)))
}

func DeleteTimeslotButton(timeslotId int) Node {
	return A(Text("delete"), Href("#"),
		hx.Delete(fmt.Sprintf("/_/timeslot/%v", timeslotId)),
		hx.Target("closest .timeslot-container"),
		hx.Swap("outerHTML swap:1s"),
	)
}
