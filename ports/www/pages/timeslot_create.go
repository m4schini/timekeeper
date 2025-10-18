package pages

import (
	"fmt"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"strconv"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/middleware"
	"timekeeper/ports/www/render"
)

func CreateTimeslotPage(event model.EventModel, rooms []model.RoomModel) Node {
	roomOptions := Group{}
	for _, room := range rooms {
		roomOptions = append(roomOptions, Option(Value(fmt.Sprintf("%v", room.ID)), Text(room.Name)))
	}

	return Shell(event.Name,
		components.PageHeader(event),
		Main(
			H2(Text("Create Timeslot")),
			components.TimeslotForm(nil, event, rooms, "POST", "/_/create/timeslot", "Create"),
		),
	)
}

type CreateTimeslotPageRoute struct {
	DB *database.Database
}

func (l *CreateTimeslotPageRoute) Method() string {
	return http.MethodGet
}

func (l *CreateTimeslotPageRoute) Pattern() string {
	return "/timeslot/new"
}

func (l *CreateTimeslotPageRoute) Handler() http.Handler {
	log := components.Logger(l)
	queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isOrganizer := middleware.IsOrganizer(request)
		if !isOrganizer {
			render.Error(log, writer, http.StatusUnauthorized, "user is not authorized", nil)
			return
		}
		var (
			eventParam   = request.URL.Query().Get("event")
			eventId, err = strconv.ParseInt(eventParam, 10, 64)
		)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid eventId", err)
			return
		}

		event, err := queries.GetEvent(int(eventId))
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to retrieve event", err)
			return
		}

		rooms, _, err := queries.GetRooms(0, 100)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to retrieve rooms", err)
			return
		}

		render.HTML(log, writer, request, CreateTimeslotPage(event, rooms))
	})
}
