package components

import (
	"fmt"
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
	"strings"
	"timekeeper/app/database/model"
)

type PaletteColor string

const (
	ColorDefault   = ""
	ColorDeepBlue  = "var(--color-deep-blue)"
	ColorDeepGreen = "var(--color-deep-green)"
	ColorSoftGrey  = "var(--color-soft-grey)"
)

func AButton(color PaletteColor, href, text string, a ...any) Node {
	return A(If(color != ColorDefault, Style("background-color: "+string(color))), Class("button"), Href(href), Textf(text, a...))

}

func EventActions(eventId int) Node {
	return Div(Class("menu"),
		CreateTimeslotButton(eventId),
		Text("Export: "),
		ExportEventMarkdownButton(eventId),
		ExportEventVocScheduleButton(eventId),
		Text("View: "),
		EventViewOnlyRole("Alle", eventId, model.RoleOrganizer, model.RoleMentor, model.RoleParticipant),
		EventViewOnlyRole("Orga", eventId, model.RoleOrganizer),
		EventViewOnlyRole("Mentor*innen", eventId, model.RoleMentor),
		EventViewOnlyRole("Teilnehmer*innen", eventId, model.RoleParticipant),
	)
}

func CopyTextBox(name, label, value string) Node {
	return Div(Style("max-width: 800px; width: 100%; display: flex; justify-content: space-between"),
		Label(For(name), Text(label)),
		Input(Type("text"), Name(name), Value(value), Attr("onclick", "this.select()"), Style("width: 600px")),
	)
}

func UrlScheduleWithRoles(eventId int, roles ...model.Role) string {
	if roles == nil || len(roles) == 0 {
		return fmt.Sprintf("/event/%v/schedule", eventId)
	}

	roleStrs := make([]string, len(roles))
	for i, role := range roles {
		roleStrs[i] = string(role)
	}

	return fmt.Sprintf("/event/%v/schedule?role=%v", eventId, strings.Join(roleStrs, ","))
}

func UrlExportVocSchedule(eventId int) string {
	return fmt.Sprintf("/event/%v/export/schedule.json", eventId)
}

func UrlExportMdSchedule(eventId int) string {
	return fmt.Sprintf("/event/%v/export/schedule.md", eventId)
}

func EventSchedule(eventId int) Node {
	return A(Class("button"), Href(UrlScheduleWithRoles(eventId)), Text("Zeitplan öffnen"))
}

func EditEvent(eventId int) Node {
	return A(Class("button"), Href(fmt.Sprintf("/event/edit/%v", eventId)), Text("Event bearbeiten"))
}

func CreateLocation() Node {
	return AButton(ColorSoftGrey, "/location/create", "Location Erstellen")
}

func EventViewOnlyRole(text string, eventId int, roles ...model.Role) Node {
	return A(Href(UrlScheduleWithRoles(eventId, roles...)), Text(text))
}

func ExportEventMarkdownButton(eventId int) Node {
	return AButton(ColorDefault, UrlExportMdSchedule(eventId), "Markdown")
}

func ExportEventVocScheduleButton(eventId int) Node {
	return AButton(ColorDefault, UrlExportVocSchedule(eventId), "VOC Schedule (Info Beamer)")
}

func CreateTimeslotButton(eventId int) Node {
	return AButton(ColorDefault, fmt.Sprintf("/timeslot/create?event=%v", eventId), "Create Timeslot")
	//return A(Href(fmt.Sprintf("/timeslot/create?event=%v", eventId)), Text("Create Timeslot"))
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

func EditLocationButton(locationId int) Node {
	return A(Text("edit"), Href(fmt.Sprintf("/location/edit/%v", locationId)))
}

func DeleteEventLocationButton(eventId, relationshipId int) Node {
	return A(Text("remove"), Href("#"),
		hx.Delete(fmt.Sprintf("/_/event/%v/location/%v", eventId, relationshipId)),
		hx.Target("closest .location-card"),
		hx.Swap("outerHTML swap:1s"),
	)
}
