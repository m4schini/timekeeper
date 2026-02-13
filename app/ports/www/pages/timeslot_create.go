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

	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func CreateTimeslotPage(event model.EventModel, parentTimeslot *model.TimeslotModel, rooms []model.RoomModel) Node {
	roomOptions := Group{}
	for _, room := range rooms {
		roomOptions = append(roomOptions, Option(Value(fmt.Sprintf("%v", room.ID)), Text(room.Name)))
	}

	return components.Shell(event.Name,
		components.PageHeader(event),
		Main(
			H2(Text("Create Timeslot")),
			components.TimeslotForm(nil, parentTimeslot, event, rooms, "POST", "/_/create/timeslot", "Create"),
		),
	)
}

type CreateTimeslotPageRoute struct {
	GetEvent                 query.GetEvent
	GetTimeslot              query.GetTimeslot
	GetRoomsOfEventLocations query.GetRoomsOfEventLocations
	Authz                    authz.Authorizer
}

func (l *CreateTimeslotPageRoute) Method() string {
	return http.MethodGet
}

func (l *CreateTimeslotPageRoute) Pattern() string {
	return "/timeslot/new"
}

func (l *CreateTimeslotPageRoute) Handler() http.Handler {
	log := components.Logger(l)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isOrganizer := middleware.IsOrganizer(request, l.Authz)
		if !isOrganizer {
			render.Error(log, writer, http.StatusUnauthorized, "user is not authorized", nil)
			return
		}
		var (
			eventParam     = request.URL.Query().Get("event")
			hasParentParam = request.URL.Query().Has("parent")
			parentParam    = request.URL.Query().Get("parent")
			parent         *model.TimeslotModel
		)
		eventId, err := strconv.ParseInt(eventParam, 10, 64)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid eventId", err)
			return
		}

		event, err := l.GetEvent.Query(query.GetEventRequest{EventId: int(eventId)})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to retrieve event", err)
			return
		}

		if hasParentParam {
			parentId, err := strconv.ParseInt(parentParam, 10, 64)
			if err != nil {
				render.Error(log, writer, http.StatusBadRequest, "invalid parentId", err)
				return
			}
			p, err := l.GetTimeslot.Query(query.GetTimeslotRequest{TimeslotId: int(parentId)})
			if err == nil {
				parent = &p
			} else {
				log.Warn("failed to get parent", zap.String("parentParam", parentParam))
			}
		}

		rooms, err := l.GetRoomsOfEventLocations.Query(query.GetRoomsOfEventLocationsRequest{EventId: event.ID})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to retrieve rooms", err)
			return
		}

		render.HTML(log, writer, request, CreateTimeslotPage(event, parent, rooms))
	})
}
