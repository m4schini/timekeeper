package components

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/app/auth/authz"
	"raumzeitalpaka/app/database/command"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/ports/www/middleware"
	"raumzeitalpaka/ports/www/render"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UpdateTimeslotRoute struct {
	UpdateTimeslot command.UpdateTimeslot
	Authz          authz.Authorizer
}

func (l *UpdateTimeslotRoute) Method() string {
	return http.MethodPost
}

func (l *UpdateTimeslotRoute) Pattern() string {
	return "/edit/timeslot/{timeslot}"
}

func (l *UpdateTimeslotRoute) Handler() http.Handler {
	log := zap.L().Named(l.Pattern())
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !middleware.IsOrganizer(request, l.Authz) {
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
			rankParam       = request.PostFormValue("rank")
		)
		model, err := ParseUpdateTimeslotModel(timeslotIdParam, eventParam, roleParam, dayParam, timeslotParam, durationParam, titleParam, noteParam, roomParam, rankParam)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}
		log.Debug("parsed create timeslot form", zap.Any("model", model))

		err = l.UpdateTimeslot.Execute(model)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to update timeslot", err)
			return
		}
		log.Debug("updated timeslot", zap.Int("id", model.TimeslotID))

		http.Redirect(writer, request, fmt.Sprintf("/event/%v/schedule", eventParam), http.StatusSeeOther)
	})
}

func ParseUpdateTimeslotModel(timeslotId, event, role, day, timeslot, duration, title, note, room, rankRaw string) (command.UpdateTimeslotRequest, error) {
	timeslotIdValue, err := strconv.ParseInt(timeslotId, 10, 64)
	if err != nil {
		return command.UpdateTimeslotRequest{}, err
	}
	rank, err := strconv.ParseInt(rankRaw, 10, 64)
	if err != nil {
		return command.UpdateTimeslotRequest{}, err
	}

	eventId, err := strconv.ParseInt(event, 10, 64)
	if err != nil {
		return command.UpdateTimeslotRequest{}, err
	}
	dayValue, err := strconv.ParseInt(day, 10, 64)
	if err != nil {
		return command.UpdateTimeslotRequest{}, err
	}
	timeslotValue, err := time.Parse("15:04", timeslot)
	if err != nil {
		return command.UpdateTimeslotRequest{}, err
	}
	roomValue, err := strconv.ParseInt(room, 10, 64)
	if err != nil {
		return command.UpdateTimeslotRequest{}, err
	}
	durationValue, err := strconv.ParseInt(duration, 10, 64)
	if err != nil {
		return command.UpdateTimeslotRequest{}, err
	}

	return command.UpdateTimeslotRequest{
		TimeslotID: int(timeslotIdValue),
		Event:      int(eventId),
		Role:       model.RoleFrom(role),
		Day:        int(dayValue),
		Timeslot:   timeslotValue,
		Duration:   time.Duration(durationValue) * time.Minute,
		Title:      title,
		Note:       note,
		Room:       int(roomValue),
		Rank:       int(rank),
	}, nil
}
