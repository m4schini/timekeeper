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
	"maragu.dev/gomponents"
)

func EventCreateForm() gomponents.Node {
	return eventForm(nil, "POST", "/_/event", "Create")
}

type CreateEventRoute struct {
	CreateEvent command.CreateEvent
	Authz       authz.Authorizer
}

func (l *CreateEventRoute) Method() string {
	return http.MethodPost
}

func (l *CreateEventRoute) Pattern() string {
	return "/event"
}

func (l *CreateEventRoute) Handler() http.Handler {
	log := zap.L().Named(l.Pattern())
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, authenticated := auth.UserFrom(request)
		if !authenticated {
			render.Error(log, writer, http.StatusUnauthorized, "unauthorized request detected", nil)
			return
		}

		err := request.ParseForm()
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}

		var (
			nameParam  = request.PostFormValue("name")
			startParam = request.PostFormValue("start")
			slugParam  = request.PostFormValue("slug")
		)
		model, err := ParseCreateEventModel(nameParam, startParam, slugParam)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}
		log.Debug("parsed create event form", zap.Any("model", model))

		id, err := l.CreateEvent.Execute(model)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to create event", err)
			return
		}
		log.Debug("created event", zap.Int("id", id))

		http.Redirect(writer, request, fmt.Sprintf("/event/%v", id), http.StatusSeeOther)
	})
}

func ParseCreateEventModel(name, start, slug string) (command.CreateEventRequest, error) {
	startDate, err := time.Parse("02.01.2006", start)
	if err != nil {
		return command.CreateEventRequest{}, err
	}

	if !EventSlugRegex.MatchString(slug) {
		return command.CreateEventRequest{}, fmt.Errorf("invalid slug")
	}

	return command.CreateEventRequest{
		Name:  name,
		Start: startDate,
		Slug:  slug,
	}, nil
}
