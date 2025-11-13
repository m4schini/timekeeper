package components

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/ports/www/middleware"
	"raumzeitalpaka/ports/www/render"
	"strconv"
)

type DeleteTimeslotRoute struct {
	DB *database.Database
}

func (l *DeleteTimeslotRoute) Method() string {
	return http.MethodDelete
}

func (l *DeleteTimeslotRoute) Pattern() string {
	return "/timeslot/{timeslot}"
}

func (l *DeleteTimeslotRoute) Handler() http.Handler {
	log := Logger(l)
	commands := l.DB.Commands
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !middleware.IsOrganizer(request) {
			render.Error(log, writer, http.StatusUnauthorized, "unauthorized request detected", nil)
			return
		}

		var (
			timeslotParam = chi.URLParam(request, "timeslot")
			timeslot, err = strconv.ParseInt(timeslotParam, 10, 64)
		)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid timeslot", err)
			return
		}

		err = commands.DeleteTimeslot(int(timeslot))
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to delete timeslot", err)
			return
		}
		log.Debug("deleted timeslot", zap.Int64("id", timeslot))
		writer.Write([]byte{})
	})
}
