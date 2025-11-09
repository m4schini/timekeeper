package components

import (
	"fmt"
	"strings"
	"time"
	"timekeeper/app/database/model"
	"timekeeper/config"

	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

func CreateTimeslotButton(eventId int, parentTimeslotId *int) Node {
	if parentTimeslotId == nil {
		return AButton(ColorDefault, fmt.Sprintf("/timeslot/new?event=%v", eventId), "Create Timeslot")
	} else {
		return AButton(ColorDefault, fmt.Sprintf("/timeslot/new?event=%v&parent=%v", eventId, *parentTimeslotId), "Create Child Timeslot")
	}
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
	return SubTimeSlot(t, nil, withActions, active, disabled, true, true, true)
}

func SubTimeSlot(t model.TimeslotModel, parent *model.TimeslotModel, withActions, active, disabled, showTime, showRoom, showRole bool) Node {
	return Div(Class("timeslot-container"), If(disabled && !active, Style("opacity: 0.5;")), If(active, Style("border-left: 8px solid var(--color-deep-green);")),
		Div(Class("timeslot-meta"),
			If(showTime, timeslotTime(t.Date(), t.Duration, false)),
			If(showRoom, timeslotRoom(t.Event.ID, t.Room.Location.ID, t.Room)),
			If(showRole, Div(Class("timeslot-roles"), RoleTag(t.Role))),
		),
		Div(Class("timeslot-info"),
			Div(Class("timeslot-info-title"), Text2(t.Title, 32)),
			Div(Class("timeslot-info-notes"), Text2(t.Note, 16)),
		),
		Iff(t.Children != nil && len(t.Children) > 0, func() Node {
			g := Group{}
			now := time.Now().In(config.Timezone()) //TODO zeigt/filtert sub events nicht richtig
			for _, child := range t.Children {
				//until := now.Before(child.Start)
				active := now.After(child.Start) && now.Before(child.Start.Add(child.Duration))

				g = append(g, SubTimeSlot(child, &t, withActions, active, false,
					child.Start != t.Start,
					child.Room.ID != t.Room.ID,
					child.Role != t.Role))

			}

			return Div(Class("timeslot-sub-events"), g)
		}),
		If(withActions, Div(Class("timeslot-action"),
			If(parent == nil, CreateTimeslotButton(t.Event.ID, &t.ID)),
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
	if withEnd {
		timeslotText := fmt.Sprintf("%v-%v", date.Format("15:04"), date.Add(duration).Format("15:04"))
		return Div(
			Class("timeslot-time"),
			Text(timeslotText),
			Title(fmt.Sprintf("%v - %v (%v)", date.Format("15:04"), date.Add(duration).Format("15:04"), duration)))
	}

	return Div(
		Class("timeslot-time"),
		Text(fmt.Sprintf("%v", date.Format("15:04"))), If(duration.Minutes() > 0, Span(Textf("%vm", duration.Minutes()))),
		Title(fmt.Sprintf("%v - %v (%v)", date.Format("15:04"), date.Add(duration).Format("15:04"), duration)))
}

func timeslotRoom(eventId, locationId int, r model.RoomModel) Node {

	desc := r.Description
	if strings.TrimSpace(desc) == "" {
		desc = r.Name
	}
	title := fmt.Sprintf("%v: %v", r.Location.Name, desc)
	if desc == r.Location.Name {
		title = r.Location.Name
	}
	return Div(Class("timeslot-room"),
		If(
			true,
			A(Textf("%v", r.Name), Href(fmt.Sprintf("/event/%v/location/%v#room-%v", eventId, locationId, r.ID)), Title(title)), //
		),
	)
}
