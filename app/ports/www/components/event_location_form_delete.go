package components

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
	"net/http"
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/ports/www/middleware"
	"raumzeitalpaka/ports/www/render"
	"strconv"
)

func DeleteEventLocationButton(eventId, relationshipId int) Node {
	return A(Class("button"), Style("background-color: var(--color-soft-red)"), Text("remove"), Href("#"),
		hx.Delete(fmt.Sprintf("/_/event/%v/location/%v", eventId, relationshipId)),
		hx.Target("closest .location-card"),
		hx.Swap("outerHTML swap:1s"),
	)
}

type DeleteLocationFromEventRoute struct {
	DB *database.Database
}

func (l *DeleteLocationFromEventRoute) Method() string {
	return http.MethodDelete
}

func (l *DeleteLocationFromEventRoute) Pattern() string {
	return "/event/{event}/location/{event_location}"
}

func (l *DeleteLocationFromEventRoute) Handler() http.Handler {
	log := zap.L().Named(l.Pattern())
	commands := l.DB.Commands
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !middleware.IsOrganizer(request) {
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

		err = commands.DeleteLocationFromEvent(int(eventLocationId))
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to delete location from event", err)
			return
		}
		log.Debug("deleted location from event", zap.Int64("id", eventLocationId))
		writer.Write([]byte{})
	})
}
