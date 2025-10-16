package pages

import (
	"fmt"
	"github.com/go-chi/chi/v5"
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

func EditTimeslotPage(timeslot model.TimeslotModel, event model.EventModel, rooms []model.RoomModel) Node {
	return Shell(event.Name,
		components.PageHeader(event),
		Main(
			H2(Text("Edit Timeslot")),
			components.TimeslotForm(&timeslot, event, rooms, "POST", fmt.Sprintf("/_/edit/timeslot/%v", timeslot.ID), "Update"),
		),
	)
}

type EditTimeslotPageRoute struct {
	DB *database.Database
}

func (l *EditTimeslotPageRoute) Method() string {
	return http.MethodGet
}

func (l *EditTimeslotPageRoute) Pattern() string {
	return "/timeslot/edit/{timeslot}"
}

func (l *EditTimeslotPageRoute) Handler() http.Handler {
	log := components.Logger(l)
	queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isOrganizer := middleware.IsOrganizer(request)
		if !isOrganizer {
			render.RenderError(log, writer, http.StatusUnauthorized, "user is not authorized", nil)
			return
		}
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

		rooms, err := queries.GetRoomsOfEventLocations(timeslot.Event.ID)
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to retrieve rooms", err)
			return
		}

		render.Render(log, writer, request, EditTimeslotPage(timeslot, timeslot.Event, rooms))
	})
}
