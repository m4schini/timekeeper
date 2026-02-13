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

	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
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
	AddLocationToEvent command.AddLocationToEvent
	Authz              authz.Authorizer
}

func (l *AddLocationToEventRoute) Method() string {
	return http.MethodPost
}

func (l *AddLocationToEventRoute) Pattern() string {
	return "/event/{event}/location"
}

func (l *AddLocationToEventRoute) Handler() http.Handler {
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
			eventParam        = request.PostFormValue("event")
			locationRoleParam = request.PostFormValue("location_role")
			locationParam     = request.PostFormValue("location")
		)
		dbmodel, err := ParseAddLocationToEventModel(eventParam, locationRoleParam, locationParam)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}
		log.Debug("parsed add location to event form", zap.Any("dbmodel", dbmodel))

		if _, isAuthorized := l.Authz.HasEventRole(userId, dbmodel.EventId, model.RoleOrganizer); !isAuthorized {
			render.Error(log, writer, http.StatusUnauthorized, "unauthorized request detected", nil)
			return
		}

		id, err := l.AddLocationToEvent.Execute(dbmodel)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to add location to event", err)
			return
		}
		log.Debug("added location to event", zap.Int("id", id))

		http.Redirect(writer, request, fmt.Sprintf("/event/%v", eventParam), http.StatusSeeOther)
	})
}

func ParseAddLocationToEventModel(event, name, location string) (command.AddLocationToEventRequest, error) {
	eventId, err := strconv.ParseInt(event, 10, 64)
	if err != nil {
		return command.AddLocationToEventRequest{}, err
	}
	locationId, err := strconv.ParseInt(location, 10, 64)
	if err != nil {
		return command.AddLocationToEventRequest{}, err
	}

	return command.AddLocationToEventRequest{
		Name:       name,
		EventId:    int(eventId),
		LocationId: int(locationId),
	}, nil
}
