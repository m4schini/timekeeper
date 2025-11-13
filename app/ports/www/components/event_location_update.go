package components

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/ports/www/middleware"
	"raumzeitalpaka/ports/www/render"
	"strconv"
)

func EventLocationUpdateForm(eventId int, eventLocation model.EventLocationModel) Node {
	return Form(Style("display: flex; flex-direction: column"), Method("POST"), Action(fmt.Sprintf("/_/event/%v/location/%v/edit", eventId, eventLocation.RelationshipId)),
		Input(Type("hidden"), Name("relationship"), Value(fmt.Sprintf("%v", eventLocation.RelationshipId))),
		//Input(Type("text"), Name("relationship_name"), Value(eventLocation.Relationship)),
		Select(Name("relationship_name"),
			Option(Value("sleep_location"), Text("Übernachtungsort"), If(eventLocation.Relationship == "sleep_location", Selected())),
			Option(Value("event_location"), Text("Eventort"), If(eventLocation.Relationship == "event_location", Selected())),
		),

		Div(
			Input(Type("checkbox"), Name("visible"), If(eventLocation.Visible, Checked())),
			Label(Text(" Sichtbar"), For("visible")),
		),

		Input(Type("text"), Name("relationship_note"), Value(eventLocation.RelationshipNote), Placeholder("kurze Anmerkung")),
		Input(Type("submit"), Value("Speichern")),
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
			eventId          = chi.URLParam(request, "event")
			relationshipId   = chi.URLParam(request, "event_location")
			visible          = request.PostFormValue("visible")
			relationshipName = request.PostFormValue("relationship_name")
			relationshipNote = request.PostFormValue("relationship_note")
		)
		model, err := ParseUpdateEventLocationModel(relationshipId, visible, relationshipName, relationshipNote)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}
		log.Debug("parsed update event location form", zap.Any("model", model))

		err = commands.UpdateLocationToEvent(model)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to update event location", err)
			return
		}
		log.Debug("updated timeslot", zap.Int("id", model.ID))

		http.Redirect(writer, request, fmt.Sprintf("/event/%v", eventId), http.StatusSeeOther)
	})
}

func ParseUpdateEventLocationModel(relationship, visible, relationshipName, relationshipNote string) (model.UpdateLocationToEventModel, error) {
	relationshipId, err := strconv.ParseInt(relationship, 10, 64)
	if err != nil {
		return model.UpdateLocationToEventModel{}, err
	}

	return model.UpdateLocationToEventModel{
		ID:      int(relationshipId),
		Name:    relationshipName,
		Note:    relationshipNote,
		Visible: visible == "on",
	}, nil
}
