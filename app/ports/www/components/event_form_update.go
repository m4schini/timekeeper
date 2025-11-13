package components

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/ports/www/middleware"
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
	DB *database.Database
}

func (l *UpdateEventRoute) Method() string {
	return http.MethodPost
}

func (l *UpdateEventRoute) Pattern() string {
	return "/event/edit"
}

func (l *UpdateEventRoute) Handler() http.Handler {
	log := Logger(l)
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
			eventParam = request.PostFormValue("event")
			nameParam  = request.PostFormValue("name")
			slugParam  = request.PostFormValue("slug")
			startParam = request.PostFormValue("start")
		)
		model, err := ParseUpdateEventModel(eventParam, nameParam, slugParam, startParam)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}
		log.Debug("parsed create event form", zap.Any("model", model))

		err = commands.UpdateEvent(model)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to create event", err)
			return
		}
		log.Debug("updated event", zap.String("id", eventParam))

		http.Redirect(writer, request, fmt.Sprintf("/event/%v", eventParam), http.StatusSeeOther)
	})
}

func ParseUpdateEventModel(event, name, slug, start string) (model.UpdateEventModel, error) {
	eventId, err := strconv.ParseInt(event, 10, 64)
	if err != nil {
		return model.UpdateEventModel{}, err
	}

	startDate, err := time.Parse("02.01.2006", start)
	if err != nil {
		return model.UpdateEventModel{}, err
	}

	if !EventSlugRegex.MatchString(slug) {
		return model.UpdateEventModel{}, fmt.Errorf("invalid slug")
	}

	return model.UpdateEventModel{
		ID:    int(eventId),
		Name:  name,
		Slug:  slug,
		Start: startDate,
	}, nil
}
