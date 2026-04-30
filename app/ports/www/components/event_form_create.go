package components

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/app/auth"
	"raumzeitalpaka/app/auth/authz"
	"raumzeitalpaka/app/database/command"
	"raumzeitalpaka/ports/www/render"
	"time"

	"go.uber.org/zap"
)

type CreateEventRoute struct {
	CreateEvent command.CreateEvent
	Authz       authz.Authorizer
}

func (l *CreateEventRoute) Method() string {
	return http.MethodPost
}

func (l *CreateEventRoute) Pattern() string {
	return "/new/event"
}

func (l *CreateEventRoute) Handler() http.Handler {
	log := zap.L().Named(l.Pattern())
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		_, authenticated := auth.UserFrom(request)
		if !authenticated {
			render.Error(log, writer, http.StatusUnauthorized, "unauthorized request detected", nil)
			return
		}

		form, err := DecodeEventForm(request, false)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to decode form", err)
			return
		}
		zap.L().Debug("decoded event form", zap.Any("form", form))

		id, err := l.CreateEvent.Execute(ctx, command.CreateEventRequest{
			Name:     form.Name,
			Slug:     form.Slug,
			Start:    time.Time(form.Start),
			End:      time.Time(form.End),
			Setup:    form.Setup,
			Teardown: form.Teardown,
		})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to create event", err)
			return
		}
		log.Debug("created event", zap.Int("id", id))

		http.Redirect(writer, request, fmt.Sprintf("/event/%v", id), http.StatusSeeOther)
	})
}
