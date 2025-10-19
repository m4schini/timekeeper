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
	"timekeeper/app/database"
	"timekeeper/ports/www/middleware"
	"timekeeper/ports/www/render"
)

func DeleteRoomButton(roomId int) Node {
	return A(Class("button"), Style("background-color: var(--color-soft-red)"), Text("remove"), Href("#"),
		hx.Delete(fmt.Sprintf("/_/room/%v", roomId)),
		hx.Target("closest li"),
		hx.Swap("outerHTML swap:1s"),
	)
}

type DeleteRoomRoute struct {
	DB *database.Database
}

func (l *DeleteRoomRoute) Method() string {
	return http.MethodDelete
}

func (l *DeleteRoomRoute) Pattern() string {
	return "/room/{room}"
}

func (l *DeleteRoomRoute) Handler() http.Handler {
	log := zap.L().Named(l.Pattern())
	commands := l.DB.Commands
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !middleware.IsOrganizer(request) {
			render.Error(log, writer, http.StatusUnauthorized, "unauthorized request detected", nil)
			return
		}

		roomId, err := strconv.ParseInt(chi.URLParam(request, "room"), 10, 64)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid roomid", err)
			return
		}

		err = commands.DeleteRoom(int(roomId))
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to delete room", err)
			return
		}
		log.Debug("deleted room to event", zap.Int64("room", roomId))
	})
}
