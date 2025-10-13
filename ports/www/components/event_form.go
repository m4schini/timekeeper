package components

import (
	"fmt"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"time"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/middleware"
	"timekeeper/ports/www/render"
)

func EventForm(event *model.EventModel, method, action, actionText string) Node {
	hasModel := event != nil
	eventId := -1
	if hasModel {
		eventId = event.ID
	}
	return Form(Method(method), Action(action), Class("form"),
		Input(Type("hidden"), Name("event"), Value(fmt.Sprintf("%v", eventId))),

		Div(
			Label(For("name"), Text("Name")),
			Input(Type("text"), Name("name"), Placeholder("Jugend hackt 2042"), Required(), Iff(hasModel, func() Node {
				return Value(event.Name)
			})),
		),

		Div(
			Label(For("start"), Text("Erster Tag")),
			Input(Type("text"), Name("start"), Placeholder("03.10.2042"), Required(), Iff(hasModel, func() Node {
				return Value(event.Start.Format("02.01.2006"))
			})),
		),

		Input(Type("submit"), Value(actionText)),
	)
}

type CreateEventRoute struct {
	DB *database.Database
}

func (l *CreateEventRoute) Method() string {
	return http.MethodPost
}

func (l *CreateEventRoute) Pattern() string {
	return "/create/event"
}

func (l *CreateEventRoute) Handler() http.Handler {
	log := zap.L().Named(l.Pattern())
	commands := l.DB.Commands
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !middleware.IsOrganizer(request) {
			render.RenderError(log, writer, http.StatusUnauthorized, "unauthorized request detected", nil)
			return
		}

		err := request.ParseForm()
		if err != nil {
			render.RenderError(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}

		var (
			nameParam  = request.PostFormValue("name")
			startParam = request.PostFormValue("start")
		)
		model, err := ParseCreateEventModel(nameParam, startParam)
		if err != nil {
			render.RenderError(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}
		log.Debug("parsed create event form", zap.Any("model", model))

		id, err := commands.CreateEvent(model)
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to create event", err)
			return
		}
		log.Debug("created event", zap.Int("id", id))

		http.Redirect(writer, request, fmt.Sprintf("/event/%v", id), http.StatusSeeOther)
	})
}

func ParseCreateEventModel(name, start string) (model.CreateEventModel, error) {
	startDate, err := time.Parse("02.01.2006", start)
	if err != nil {
		return model.CreateEventModel{}, err
	}

	return model.CreateEventModel{
		Name:  name,
		Start: startDate,
	}, nil
}
