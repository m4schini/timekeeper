package components

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
	"net/http"
	"strconv"
	"time"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/config"
	"timekeeper/ports/www/middleware"
	"timekeeper/ports/www/render"
)

func TimeSlot(t model.TimeslotModel, withActions, disabled bool) Node {
	return Div(Class("timeslot-container"), If(disabled, Style("opacity: 0.5;")),
		Div(Class("timeslot-meta"),
			timeslotTime(t.Event.Start, t.Start, t.Day),
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

func DeleteTimeslotButton(timeslotId int) Node {
	return A(Text("delete"), Href("#"),
		hx.Delete(fmt.Sprintf("/_/timeslot/%v", timeslotId)),
		hx.Target("closest .timeslot-container"),
		hx.Swap("outerHTML swap:1s"),
	)
}

func EditTimeslotButton(timeslotId int) Node {
	return A(Text("edit"), Href(fmt.Sprintf("/edit/timeslot/%v", timeslotId)))
}

func DuplicateTimeslotButton(timeslotId int) Node {
	return A(Text("duplicate"), Href(fmt.Sprintf("/duplicate/timeslot/%v", timeslotId)))
}

func FullTimeSlot(t model.TimeslotModel, disabled bool) Node {
	return Div(Class("full-timeslot-container"), If(disabled, Style("opacity: 0.5;")),
		timeslotTime(t.Event.Start, t.Start, t.Day),
		timeslotRoom(t.Event.ID, t.Room.Location.ID, t.Room),
		Div(Class("timeslot-roles"), RoleTag(t.Role)),
		Div(Class("timeslot-info"),
			Div(Class("timeslot-info-title"), Text(t.Title)),
			Div(Class("full-timeslot-info-notes"), Text(t.Note)),
		),
		Div(Class("timeslot-map")), //LocationCrop(t.Location.X, t.Location.Y, t.Location.Width, t.Location.Height, 100),
	)
}

func CompactTimeSlot(t model.TimeslotModel, disabled bool) Node {
	return Div(Class("compact-timeslot-container"), If(disabled, Style("opacity: 0.5;")),
		timeslotTime(t.Event.Start, t.Start, t.Day),
		timeslotRoom(t.Event.ID, t.Room.Location.ID, t.Room),
		Div(Class("timeslot-roles"), RoleTag(t.Role)),
		Div(Class("timeslot-info"),
			Div(Class("timeslot-info-title"), Text(t.Title)),
			Div(Class("timeslot-info-notes"), Text(t.Note)),
		),
		Div(Class("timeslot-map")), //LocationCrop(t.Location.X, t.Location.Y, t.Location.Width, t.Location.Height, 100),
	)
}

func timeslotTime(startDate, timeslot time.Time, day int) Node {
	date := time.Date(startDate.Year(), startDate.Month(), startDate.Day(),
		timeslot.Hour(), timeslot.Minute(), timeslot.Second(), timeslot.Nanosecond(),
		config.Timezone())
	offset := time.Duration(day) * 24 * time.Hour
	return Div(
		Class("timeslot-time"),
		Text(timeslot.Format("15:04")),
		Title(fmt.Sprintf("%v", date.
			Add(offset).
			Format(time.RFC1123Z))))
}

func timeslotRoom(eventId, locationId int, r model.RoomModel) Node {
	return Div(Class("timeslot-room"),
		If(
			true,
			A(Textf("%v", r.Name), Href(fmt.Sprintf("/event/%v/location/%v#%v", eventId, locationId, r.ID))),
		),
	)
}

func RoleTag(role model.Role) Node {
	switch role {
	case model.RoleOrganizer:
		return Span(Class("role role-o"), Text("Orga"))
	case model.RoleMentor:
		return Span(Class("role role-m"), Text("Mentor*innen"))
	case model.RoleParticipant:
		return Span(Class("role role-t"), Text("Teilnehmer*innen"))
	default:
		return Span(Class("role role-o"), Text("Orga"))
	}
}

type DeleteTimeslotRoute struct {
	DB *database.Database
}

func (l *DeleteTimeslotRoute) Method() string {
	return http.MethodDelete
}

func (l *DeleteTimeslotRoute) Pattern() string {
	return "/timeslot/{timeslot}"
}

func (l *DeleteTimeslotRoute) Handler() http.Handler {
	log := zap.L().Named(l.Pattern())
	commands := l.DB.Commands
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !middleware.IsOrganizer(request) {
			render.RenderError(log, writer, http.StatusUnauthorized, "unauthorized request detected", nil)
			return
		}

		var (
			timeslotParam = chi.URLParam(request, "timeslot")
			timeslot, err = strconv.ParseInt(timeslotParam, 10, 64)
		)
		if err != nil {
			render.RenderError(log, writer, http.StatusBadRequest, "invalid timeslot", err)
			return
		}

		err = commands.DeleteTimeslot(int(timeslot))
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to delete timeslot", err)
			return
		}
		log.Debug("deleted timeslot", zap.Int64("id", timeslot))
		writer.Write([]byte{})
	})
}
