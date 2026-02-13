package pages

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/app/auth/authz"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/app/database/query"
	"raumzeitalpaka/ports/www/components"
	"raumzeitalpaka/ports/www/middleware"
	"raumzeitalpaka/ports/www/render"
	"strconv"

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func EditTimeslotPage(timeslot model.TimeslotModel, parentTimeslot *model.TimeslotModel, event model.EventModel, rooms []model.RoomModel) Node {
	return components.Shell(event.Name,
		components.PageHeader(event),
		Main(
			H2(Text("Edit Timeslot")),
			components.TimeslotForm(&timeslot, parentTimeslot, event, rooms, "POST", fmt.Sprintf("/_/edit/timeslot/%v", timeslot.ID), "Update"),
		),
	)
}

type EditTimeslotPageRoute struct {
	GetTimeslot              query.GetTimeslot
	GetRoomsOfEventLocations query.GetRoomsOfEventLocations
	Authz                    authz.Authorizer
}

func (l *EditTimeslotPageRoute) Method() string {
	return http.MethodGet
}

func (l *EditTimeslotPageRoute) Pattern() string {
	return "/timeslot/edit/{timeslot}"
}

func (l *EditTimeslotPageRoute) Handler() http.Handler {
	log := components.Logger(l)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isOrganizer := middleware.IsOrganizer(request, l.Authz)
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

		timeslot, err := l.GetTimeslot.Query(query.GetTimeslotRequest{TimeslotId: int(timeslotId)})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to retrieve timeslot", err)
			return
		}

		rooms, err := l.GetRoomsOfEventLocations.Query(query.GetRoomsOfEventLocationsRequest{EventId: timeslot.Event.ID})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to retrieve rooms", err)
			return
		}

		//TODO parent
		render.HTML(log, writer, request, EditTimeslotPage(timeslot, nil, timeslot.Event, rooms))
	})
}
