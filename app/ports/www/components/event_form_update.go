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
	"time"

	"go.uber.org/zap"
	"maragu.dev/gomponents"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func EditEvent(eventId int) Node {
	return A(Class("button"), Href(fmt.Sprintf("/event/%v/edit", eventId)), Text("Event bearbeiten"))
}

func EventUpdateForm(event model.EventModel) gomponents.Node {
	return eventForm(&event, "POST", "/_/event/edit", "Update")
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
		userId, authenticated := auth.UserFrom(request)
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
			eventParam = request.PostFormValue("event")
			nameParam  = request.PostFormValue("name")
			slugParam  = request.PostFormValue("slug")
			startParam = request.PostFormValue("start")
		)
		eventModel, err := ParseUpdateEventModel(eventParam, nameParam, slugParam, startParam)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}
		log.Debug("parsed create event form", zap.Any("eventModel", eventModel))

		if _, isAuthorized := l.Authz.HasEventRole(userId, eventModel.ID, model.RoleOrganizer); !isAuthorized {
			render.Error(log, writer, http.StatusUnauthorized, "unauthorized request detected", nil)
			return
		}

		err = l.UpdateEvent.Execute(eventModel)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to create event", err)
			return
		}
		log.Debug("updated event", zap.String("id", eventParam))

		http.Redirect(writer, request, fmt.Sprintf("/event/%v", eventParam), http.StatusSeeOther)
	})
}

func ParseUpdateEventModel(event, name, slug, start string) (command.UpdateEventRequest, error) {
	eventId, err := strconv.ParseInt(event, 10, 64)
	if err != nil {
		return command.UpdateEventRequest{}, err
	}

	startDate, err := time.Parse("02.01.2006", start)
	if err != nil {
		return command.UpdateEventRequest{}, err
	}

	if !EventSlugRegex.MatchString(slug) {
		return command.UpdateEventRequest{}, fmt.Errorf("invalid slug")
	}

	return command.UpdateEventRequest{
		ID:    int(eventId),
		Name:  name,
		Slug:  slug,
		Start: startDate,
	}, nil
}
