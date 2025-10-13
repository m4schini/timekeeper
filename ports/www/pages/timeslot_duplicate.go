package pages

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"strconv"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/render"
)

func DuplicateTimeslotPage(timeslot model.TimeslotModel, event model.EventModel, rooms []model.RoomModel) Node {
	roomOptions := Group{}
	for _, room := range rooms {
		roomOptions = append(roomOptions, Option(Value(fmt.Sprintf("%v", room.ID)), Text(room.Name)))
	}

	return Shell(
		components.PageHeader(event, false),
		Main(
			H1(Textf("%v (%v)", event.Name, event.Start.Format("2006.01.02"))),
			components.TimeslotForm(&timeslot, event, rooms, "POST", "/_/create/timeslot", "Duplicate"),
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

func (l *DuplicateTimeslotPageRoute) UseCache() bool {
	return false
}

func (l *DuplicateTimeslotPageRoute) Handler() http.Handler {
	log := zap.L().Named(l.Pattern())
	queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var (
			timeslotParam   = chi.URLParam(request, "timeslot")
			timeslotId, err = strconv.ParseInt(timeslotParam, 10, 64)
		)
		if err != nil {
			render.RenderError(log, writer, http.StatusBadRequest, "invalid timeslotId", err)
			return
		}

		timeslot, err := queries.GetTimeslot(int(timeslotId))
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to retrieve timeslot", err)
			return
		}

		rooms, _, err := queries.GetRooms(0, 100)
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to retrieve rooms", err)
			return
		}

		render.Render(log, writer, request, DuplicateTimeslotPage(timeslot, timeslot.Event, rooms))
	})
}
