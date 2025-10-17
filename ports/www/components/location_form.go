package components

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"strconv"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/middleware"
	"timekeeper/ports/www/render"
)

func AddLocationForm(eventId int, locations []model.LocationModel) Node {
	locationOptions := Group{}
	for _, location := range locations {
		locationOptions = append(locationOptions, Option(Value(fmt.Sprintf("%v", location.ID)), Text(location.Name)))
	}

	return Form(Method("POST"), Action("/_/event/{event}/location"),
		Div(
			Input(Type("hidden"), Name("event"), Value(fmt.Sprintf("%v", eventId))),

			Label(For("location"), Text("Ort")),
			Select(Name("location"), Required(),
				locationOptions,
			),

			Label(For("location_role"), Text("als")),
			Select(Name("location_role"), Required(),
				Option(Value("event_location"), Text("Eventort")),
				Option(Value("sleep_location"), Text("Übernachtungsort")),
			),

			Input(Type("submit"), Value("Location hinzufügen")),
		),
	)
}

type AddLocationToEventRoute struct {
	DB *database.Database
}

func (l *AddLocationToEventRoute) Method() string {
	return http.MethodPost
}

func (l *AddLocationToEventRoute) Pattern() string {
	return "/event/{event}/location"
}

func (l *AddLocationToEventRoute) Handler() http.Handler {
	log := Logger(l)
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
			eventParam        = request.PostFormValue("event")
			locationRoleParam = request.PostFormValue("location_role")
			locationParam     = request.PostFormValue("location")
		)
		model, err := ParseAddLocationToEventModel(eventParam, locationRoleParam, locationParam)
		if err != nil {
			render.RenderError(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}
		log.Debug("parsed add location to event form", zap.Any("model", model))

		id, err := commands.AddLocationToEvent(model)
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to add location to event", err)
			return
		}
		log.Debug("added location to event", zap.Int("id", id))

		http.Redirect(writer, request, fmt.Sprintf("/event/%v", eventParam), http.StatusSeeOther)
	})
}

func ParseAddLocationToEventModel(event, name, location string) (model.AddLocationToEventModel, error) {
	eventId, err := strconv.ParseInt(event, 10, 64)
	if err != nil {
		return model.AddLocationToEventModel{}, err
	}
	locationId, err := strconv.ParseInt(location, 10, 64)
	if err != nil {
		return model.AddLocationToEventModel{}, err
	}

	return model.AddLocationToEventModel{
		Name:       name,
		EventId:    int(eventId),
		LocationId: int(locationId),
	}, nil
}

type DeleteLocationFromEventRoute struct {
	DB *database.Database
}

func (l *DeleteLocationFromEventRoute) Method() string {
	return http.MethodDelete
}

func (l *DeleteLocationFromEventRoute) Pattern() string {
	return "/event/{event}/location/{event_location}"
}

func (l *DeleteLocationFromEventRoute) Handler() http.Handler {
	log := zap.L().Named(l.Pattern())
	commands := l.DB.Commands
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !middleware.IsOrganizer(request) {
			render.RenderError(log, writer, http.StatusUnauthorized, "unauthorized request detected", nil)
			return
		}

		var (
			eventLocationParam   = chi.URLParam(request, "event_location")
			eventLocationId, err = strconv.ParseInt(eventLocationParam, 10, 64)
		)
		if err != nil {
			render.RenderError(log, writer, http.StatusBadRequest, "invalid event_location", err)
			return
		}

		err = commands.DeleteLocationFromEvent(int(eventLocationId))
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to delete location from event", err)
			return
		}
		log.Debug("deleted location from event", zap.Int64("id", eventLocationId))
		writer.Write([]byte{})
	})
}
