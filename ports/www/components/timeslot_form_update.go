package components

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	"net/http"
	"strconv"
	"time"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/middleware"
	"timekeeper/ports/www/render"
)

type UpdateTimeslotRoute struct {
	DB *database.Database
}

func (l *UpdateTimeslotRoute) Method() string {
	return http.MethodPost
}

func (l *UpdateTimeslotRoute) Pattern() string {
	return "/edit/timeslot/{timeslot}"
}

func (l *UpdateTimeslotRoute) Handler() http.Handler {
	log := zap.L().Named(l.Pattern())
	commands := l.DB.Commands
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !middleware.IsOrganizer(request) {
			render.Error(log, writer, http.StatusUnauthorized, "unauthorized request detected", nil)
			return
		}

		err := request.ParseForm()
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}

		var (
			timeslotIdParam = chi.URLParam(request, "timeslot")
			eventParam      = request.PostFormValue("event")
			roleParam       = request.PostFormValue("role")
			dayParam        = request.PostFormValue("day")
			timeslotParam   = request.PostFormValue("timeslot")
			durationParam   = request.PostFormValue("duration")
			titleParam      = request.PostFormValue("title")
			noteParam       = request.PostFormValue("note")
			roomParam       = request.PostFormValue("room")
		)
		model, err := ParseUpdateTimeslotModel(timeslotIdParam, eventParam, roleParam, dayParam, timeslotParam, durationParam, titleParam, noteParam, roomParam)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}
		log.Debug("parsed create timeslot form", zap.Any("model", model))

		err = commands.UpdateTimeslot(model)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to update timeslot", err)
			return
		}
		log.Debug("updated timeslot", zap.Int("id", model.ID))

		http.Redirect(writer, request, fmt.Sprintf("/event/%v/schedule", eventParam), http.StatusSeeOther)
	})
}

func ParseUpdateTimeslotModel(timeslotId, event, role, day, timeslot, duration, title, note, room string) (model.UpdateTimeslotModel, error) {
	timeslotIdValue, err := strconv.ParseInt(timeslotId, 10, 64)
	if err != nil {
		return model.UpdateTimeslotModel{}, err
	}

	eventId, err := strconv.ParseInt(event, 10, 64)
	if err != nil {
		return model.UpdateTimeslotModel{}, err
	}
	dayValue, err := strconv.ParseInt(day, 10, 64)
	if err != nil {
		return model.UpdateTimeslotModel{}, err
	}
	timeslotValue, err := time.Parse("15:04", timeslot)
	if err != nil {
		return model.UpdateTimeslotModel{}, err
	}
	roomValue, err := strconv.ParseInt(room, 10, 64)
	if err != nil {
		return model.UpdateTimeslotModel{}, err
	}
	durationValue, err := strconv.ParseInt(duration, 10, 64)
	if err != nil {
		return model.UpdateTimeslotModel{}, err
	}

	return model.UpdateTimeslotModel{
		ID:       int(timeslotIdValue),
		Event:    int(eventId),
		Role:     model.RoleFrom(role),
		Day:      int(dayValue),
		Timeslot: timeslotValue,
		Duration: time.Duration(durationValue) * time.Minute,
		Title:    title,
		Note:     note,
		Room:     int(roomValue),
	}, nil
}
