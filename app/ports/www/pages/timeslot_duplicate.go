package pages

import (
	"fmt"
	"net/http"
	"strconv"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/middleware"
	"timekeeper/ports/www/render"

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func DuplicateTimeslotPage(timeslot model.TimeslotModel, parentTimeslot *model.TimeslotModel, event model.EventModel, rooms []model.RoomModel) Node {
	roomOptions := Group{}
	for _, room := range rooms {
		roomOptions = append(roomOptions, Option(Value(fmt.Sprintf("%v", room.ID)), Text(room.Name)))
	}

	return Shell(event.Name,
		components.PageHeader(event),
		Main(
			H2(Text("Duplicate Timeslot")),
			components.TimeslotForm(&timeslot, parentTimeslot, event, rooms, "POST", "/_/create/timeslot", "Duplicate"),
		),
	)
}

type DuplicateTimeslotPageRoute struct {
	DB *database.Database
}

func (l *DuplicateTimeslotPageRoute) Method() string {
	return http.MethodGet
}

func (l *DuplicateTimeslotPageRoute) Pattern() string {
	return "/timeslot/duplicate/{timeslot}"
}

func (l *DuplicateTimeslotPageRoute) Handler() http.Handler {
	log := components.Logger(l)
	queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isOrganizer := middleware.IsOrganizer(request)
		if !isOrganizer {
			render.Error(log, writer, http.StatusUnauthorized, "user is not authorized", nil)
			return
		}
		var (
			timeslotParam   = chi.URLParam(request, "timeslot")
			timeslotId, err = strconv.ParseInt(timeslotParam, 10, 64)
		)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid timeslotId", err)
			return
		}

		timeslot, err := queries.GetTimeslot(int(timeslotId))
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to retrieve timeslot", err)
			return
		}

		rooms, err := queries.GetRoomsOfEventLocations(timeslot.Event.ID)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to retrieve rooms", err)
			return
		}

		//TODO parent
		render.HTML(log, writer, request, DuplicateTimeslotPage(timeslot, nil, timeslot.Event, rooms))
	})
}
