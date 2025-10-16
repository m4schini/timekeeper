package components

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"strconv"
	"strings"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/middleware"
	"timekeeper/ports/www/render"
)

func SetEventLocationNote(eventId int, eventLocation model.EventLocationModel) Node {
	return Form(Method("POST"), Action(fmt.Sprintf("/_/event/%v/location/%v/edit", eventId, eventLocation.RelationshipId)),
		Input(Type("hidden"), Name("relationship"), Value(fmt.Sprintf("%v", eventLocation.RelationshipId))),
		//Input(Type("text"), Name("relationship_name"), Value(eventLocation.Relationship)),
		Select(Name("relationship_name"),
			Option(Value("sleep_location"), Text("Übernachtungsort"), If(eventLocation.Relationship == "sleep_location", Selected())),
			Option(Value("event_location"), Text("Eventort"), If(eventLocation.Relationship == "event_location", Selected())),
		),
		Input(Type("text"), Name("relationship_note"), Value(eventLocation.RelationshipNote), Placeholder("kurze Anmerkung")),
		Input(Type("submit"), Value("Speichern")),
	)
}

func EventLocationCard(event model.EventModel, eventLocation model.EventLocationModel, withActions bool) Node {
	address := eventLocation.Address
	var eventRole string
	switch strings.TrimSpace(strings.ToLower(eventLocation.Relationship)) {
	case "sleep_location":
		eventRole = "Übernachtung"
		break
	case "event_location":
		eventRole = "Event Location"
		break
	}

	return Div(Class("location-card"),
		If(!withActions, Div(Strong(Text(eventRole)))),
		Iff(withActions, func() Node {
			return SetEventLocationNote(event.ID, eventLocation)
		}),
		If(!withActions, Div(Text(eventLocation.RelationshipNote))), Br(),
		Div(Style("white-space: pre-line"), Text(eventLocation.Name)),
		Div(Style("white-space: pre-line"), Iff(eventLocation.Address != nil, func() Node {
			return Textf(`%v %v
%v %v
`, address.Road, address.HouseNumber, address.Postcode, address.City)
		}),
			Iff(withActions, func() Node {
				return Div(Style("display: flex; gap: 1rem"),
					EditLocationButton(eventLocation.ID),
					DeleteEventLocationButton(event.ID, eventLocation.RelationshipId),
				)
			}),
		),
	)
}

type UpdateEventLocationRoute struct {
	DB *database.Database
}

func (l *UpdateEventLocationRoute) Method() string {
	return http.MethodPost
}

func (l *UpdateEventLocationRoute) Pattern() string {
	return "/event/{event}/location/{event_location}/edit"
}

func (l *UpdateEventLocationRoute) Handler() http.Handler {
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
			eventId          = chi.URLParam(request, "event")
			relationshipId   = chi.URLParam(request, "event_location")
			relationshipName = request.PostFormValue("relationship_name")
			relationshipNote = request.PostFormValue("relationship_note")
		)
		model, err := ParseUpdateEventLocationModel(relationshipId, relationshipName, relationshipNote)
		if err != nil {
			render.RenderError(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}
		log.Debug("parsed update event location form", zap.Any("model", model))

		err = commands.UpdateLocationToEvent(model)
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to update event location", err)
			return
		}
		log.Debug("updated timeslot", zap.Int("id", model.ID))

		http.Redirect(writer, request, fmt.Sprintf("/event/%v", eventId), http.StatusSeeOther)
	})
}

func ParseUpdateEventLocationModel(relationship, relationshipName, relationshipNote string) (model.UpdateLocationToEventModel, error) {
	relationshipId, err := strconv.ParseInt(relationship, 10, 64)
	if err != nil {
		return model.UpdateLocationToEventModel{}, err
	}

	return model.UpdateLocationToEventModel{
		ID:   int(relationshipId),
		Name: relationshipName,
		Note: relationshipNote,
	}, nil
}
