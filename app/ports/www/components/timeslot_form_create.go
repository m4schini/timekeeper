package components

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/middleware"
	"timekeeper/ports/www/render"

	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func TimeslotForm(ts *model.TimeslotModel, event model.EventModel, rooms []model.RoomModel, method, action, actionText string) Node {
	hasTs := ts != nil
	if !hasTs {
		ts = &model.TimeslotModel{Room: model.RoomModel{}, Event: model.EventModel{}}
	}
	roomOptions := Group{}
	for _, room := range rooms {
		roomOptions = append(roomOptions, Option(
			Value(fmt.Sprintf("%v", room.ID)),
			Textf("%v: %v", room.Location.Name, room.Name),
			If(hasTs && room.ID == ts.Room.ID, Selected())))
	}

	return Form(Method(method), Action(action), Class("form"),
		Input(Type("hidden"), Name("event"), Value(fmt.Sprintf("%v", event.ID))),

		Div(Class("param"),
			Label(For("role"), Text("Rolle")),
			Select(Name("role"), Required(),
				Option(Value("Organizer"), Text("Orga"), If(hasTs && ts.Role == model.RoleOrganizer, Selected())),
				Option(Value("Mentor"), Text("Mentor*innen"), If(hasTs && ts.Role == model.RoleMentor, Selected())),
				Option(Value("Participant"), Text("Teilnehmer*innen"), If(!hasTs || ts.Role == model.RoleParticipant, Selected())),
			),
		),

		Div(Class("param"),
			Label(For("day"), Text("Tag")),
			Input(Type("number"), Name("day"), Placeholder("0"), Min("0"), Required(), If(hasTs, Value(fmt.Sprintf("%v", ts.Day)))),
		),

		Div(Class("param"),
			Label(For("timeslot"), Text("Zeit")),
			Input(Type("text"), Name("timeslot"), Placeholder("08:00"), Pattern(`^([01]\d|2[0-3]):([0-5]\d)$`), Required(), If(hasTs, Value(ts.Start.Format("15:04")))),
		),

		Div(Class("param"),
			Label(For("duration"), Text("Dauer in Minuten")),
			Input(Type("number"), Name("duration"), Placeholder("60"), Min("0"), Max("1440"), Required(), If(hasTs, Value(fmt.Sprintf("%v", ts.Duration.Minutes())))),
		),

		Div(Class("param"),
			Label(For("title"), Text("Titel")),
			Input(Type("text"), Name("title"), Placeholder("Title"), Required(), If(hasTs, Value(ts.Title))),
		),

		Div(Class("param"),
			Label(For("note"), Text("Notiz")),
			Textarea(Name("note"), Placeholder("Notiz"), Rows("4"), Cols("50"), If(hasTs, Text(ts.Note))),
			Label(For("note"), Text("Pixelhack emojis: "), A(Href("/help/pixelhack"), Textf("PixelHack Overview"))),
		),

		Div(Class("param"),
			Label(For("room"), Text("Raum")),
			Select(Name("room"), Required(),
				roomOptions,
			),
		),

		Input(Type("submit"), Value(actionText)),
	)
}

type CreateTimeslotRoute struct {
	DB *database.Database
}

func (l *CreateTimeslotRoute) Method() string {
	return http.MethodPost
}

func (l *CreateTimeslotRoute) Pattern() string {
	return "/create/timeslot"
}

func (l *CreateTimeslotRoute) Handler() http.Handler {
	log := Logger(l)
	commands := l.DB.Commands
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !middleware.IsOrganizer(request) {
			render.Error(log, writer, http.StatusUnauthorized, "unauthorized request detected", nil)
			return
		}

		err := request.ParseForm()
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}

		var (
			eventParam    = request.PostFormValue("event")
			roleParam     = request.PostFormValue("role")
			dayParam      = request.PostFormValue("day")
			timeslotParam = request.PostFormValue("timeslot")
			durationParam = request.PostFormValue("duration")
			titleParam    = request.PostFormValue("title")
			noteParam     = request.PostFormValue("note")
			roomParam     = request.PostFormValue("room")
		)
		model, err := ParseCreateTimeslotModel(eventParam, roleParam, dayParam, timeslotParam, durationParam, titleParam, noteParam, roomParam)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}
		log.Debug("parsed create timeslot form", zap.Any("model", model))

		id, err := commands.CreateTimeslot(model)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to create timeslot", err)
			return
		}
		log.Debug("created timeslot", zap.Int("id", id))

		http.Redirect(writer, request, fmt.Sprintf("/event/%v/schedule", eventParam), http.StatusSeeOther)
	})
}

func ParseCreateTimeslotModel(event, role, day, timeslot, duration, title, note, room string) (model.CreateTimeslotModel, error) {
	eventId, err := strconv.ParseInt(event, 10, 64)
	if err != nil {
		return model.CreateTimeslotModel{}, err
	}
	dayValue, err := strconv.ParseInt(day, 10, 64)
	if err != nil {
		return model.CreateTimeslotModel{}, err
	}
	timeslotValue, err := time.Parse("15:04", timeslot)
	if err != nil {
		return model.CreateTimeslotModel{}, err
	}
	roomValue, err := strconv.ParseInt(room, 10, 64)
	if err != nil {
		return model.CreateTimeslotModel{}, err
	}
	durationValue, err := strconv.ParseInt(duration, 10, 64)
	if err != nil {
		return model.CreateTimeslotModel{}, err
	}

	return model.CreateTimeslotModel{
		Event:    int(eventId),
		Role:     model.RoleFrom(role),
		Day:      int(dayValue),
		Timeslot: timeslotValue,
		Duration: time.Duration(durationValue) * time.Minute,
		Title:    title,
		Note:     note,
		Room:     int(roomValue),
	}, nil
}
