package components

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/app/auth"
	"raumzeitalpaka/app/auth/authz"
	"raumzeitalpaka/app/database/command"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/ports/www/render"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

func DeleteEventLocationButton(eventId, relationshipId int) Node {
	return A(Class("button"), Style("background-color: var(--color-soft-red)"), Text("remove"), Href("#"),
		hx.Delete(fmt.Sprintf("/_/event/%v/location/%v", eventId, relationshipId)),
		hx.Target("closest .location-card"),
		hx.Swap("outerHTML swap:1s"),
	)
}

type DeleteLocationFromEventRoute struct {
	RemoveLocationFromEvent command.RemoveLocationFromEvent
	Authz                   authz.Authorizer
}

func (l *DeleteLocationFromEventRoute) Method() string {
	return http.MethodDelete
}

func (l *DeleteLocationFromEventRoute) Pattern() string {
	return "/event/{event}/location/{event_location}"
}

func (l *DeleteLocationFromEventRoute) Handler() http.Handler {
	log := zap.L().Named(l.Pattern())
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userId, authenticated := auth.UserFrom(request)
		if !authenticated {
			render.Error(log, writer, http.StatusUnauthorized, "unauthorized request detected", nil)
			return
		}

		var (
			eventLocationParam   = chi.URLParam(request, "event_location")
			eventLocationId, err = strconv.ParseInt(eventLocationParam, 10, 64)
		)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid event_location", err)
			return
		}

		if isAuthorized := l.Authz.HasRole(userId, model.RoleOrganizer); !isAuthorized {
			render.Error(log, writer, http.StatusUnauthorized, "unauthorized request detected", nil)
			return
		}

		err = l.RemoveLocationFromEvent.Execute(command.RemoveLocationFromEventRequest{EventLocationRelationID: int(eventLocationId)})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to delete location from event", err)
			return
		}
		log.Debug("deleted location from event", zap.Int64("id", eventLocationId))
		writer.Write([]byte{})
	})
}
