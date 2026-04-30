package components

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/app/auth"
	"raumzeitalpaka/app/auth/authz"
	"raumzeitalpaka/app/database/command"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/ports/www/render"
	"time"

	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func EditEvent(eventId int) Node {
	return A(Class("button"), Href(fmt.Sprintf("/event/%v/edit", eventId)), Text("Event bearbeiten"))
}

type UpdateEventRoute struct {
	UpdateEvent command.UpdateEvent
	Authz       authz.Authorizer
}

func (l *UpdateEventRoute) Method() string {
	return http.MethodPost
}

func (l *UpdateEventRoute) Pattern() string {
	return "/event/edit"
}

func (l *UpdateEventRoute) Handler() http.Handler {
	log := Logger(l)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		userId, authenticated := auth.UserFrom(request)
		if !authenticated {
			render.Error(log, writer, http.StatusUnauthorized, "unauthorized request detected", nil)
			return
		}

		form, err := DecodeEventForm(request, true)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to decode form", err)
			return
		}
		log.Info("update event", zap.Any("form", form))

		if _, isAuthorized := l.Authz.HasEventRole(userId, form.Event, model.RoleOrganizer); !isAuthorized {
			render.Error(log, writer, http.StatusUnauthorized, "unauthorized request detected", nil)
			return
		}

		err = l.UpdateEvent.Execute(ctx, command.UpdateEventRequest{
			ID:    form.Event,
			Name:  form.Name,
			Slug:  form.Slug,
			Start: time.Time(form.Start),
		})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to create event", err)
			return
		}
		log.Debug("updated event", zap.Int("id", form.Event))

		http.Redirect(writer, request, fmt.Sprintf("/event/%v", form.Event), http.StatusSeeOther)
	})
}
