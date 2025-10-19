package components

import (
	"fmt"
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
	"time"
	"timekeeper/app/database/model"
)

func CreateTimeslotButton(eventId int) Node {
	return AButton(ColorDefault, fmt.Sprintf("/timeslot/new?event=%v", eventId), "Create Timeslot")
}

func EditTimeslotButton(timeslotId int) Node {
	return A(Class("button"), Text("edit"), Href(fmt.Sprintf("/timeslot/edit/%v", timeslotId)))
}

func DuplicateTimeslotButton(timeslotId int) Node {
	return A(Class("button"), Text("duplicate"), Href(fmt.Sprintf("/timeslot/duplicate/%v", timeslotId)))
}

func DeleteTimeslotButton(timeslotId int) Node {
	return A(Class("button"), Style("background-color: var(--color-soft-red)"), Text("delete"), Href("#"),
		hx.Delete(fmt.Sprintf("/_/timeslot/%v", timeslotId)),
		hx.Target("closest .timeslot-container"),
		hx.Swap("outerHTML swap:1s"),
	)
}

func TimeSlot(t model.TimeslotModel, withActions, active, disabled bool) Node {
	return Div(Class("timeslot-container"), If(disabled && !active, Style("opacity: 0.5;")), If(active, Style("border-left: 8px solid var(--color-deep-green);")),
		Div(Class("timeslot-meta"),
			timeslotTime(t.Date(), t.Duration, false),
			timeslotRoom(t.Event.ID, t.Room.Location.ID, t.Room),
			Div(Class("timeslot-roles"), RoleTag(t.Role)),
		),
		Div(Class("timeslot-info"),
			Div(Class("timeslot-info-title"), Text(t.Title)),
			Div(Class("timeslot-info-notes"), Text(t.Note)),
		),
		If(withActions, Div(Class("timeslot-action"),
			EditTimeslotButton(t.ID),
			DuplicateTimeslotButton(t.ID),
			DeleteTimeslotButton(t.ID),
		)),
	)
}

func CompactTimeSlot(t model.TimeslotModel, active, disabled bool) Node {
	return Div(Class("compact-timeslot-container"), If(disabled && !active, Style("opacity: 0.5;")), If(active, Style("border-left: 8px solid var(--color-deep-green);")),
		timeslotTime(t.Date(), t.Duration, true),
		timeslotRoom(t.Event.ID, t.Room.Location.ID, t.Room),
		Div(Class("timeslot-roles"), RoleTag(t.Role)),
		Div(Class("timeslot-info"),
			Div(Class("timeslot-info-title"), Text(t.Title)),
			Div(Class("timeslot-info-notes"), Text(t.Note)),
		),
		Div(Class("timeslot-map")), //LocationCrop(t.Location.X, t.Location.Y, t.Location.Width, t.Location.Height, 100),
	)
}

func timeslotTime(date time.Time, duration time.Duration, withEnd bool) Node {
	var timeslotText string
	if withEnd {
		timeslotText = fmt.Sprintf("%v-%v", date.Format("15:04"), date.Add(duration).Format("15:04"))
	} else {
		timeslotText = fmt.Sprintf("%v", date.Format("15:04"))
	}

	return Div(
		Class("timeslot-time"),
		Text(timeslotText),
		Title(fmt.Sprintf("%v - %v (%v)", date.Format("15:04"), date.Add(duration).Format("15:04"), duration)))
}

func timeslotRoom(eventId, locationId int, r model.RoomModel) Node {
	title := fmt.Sprintf("%v: %v", r.Location.Name, r.Name)
	return Div(Class("timeslot-room"),
		If(
			true,
			A(Textf("%v", r.Name), Href(fmt.Sprintf("/event/%v/location/%v#room-%v", eventId, locationId, r.ID)), Title(title)), //
		),
	)
}
