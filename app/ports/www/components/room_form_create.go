package components

import (
	"fmt"
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

func CreateRoomForm(location model.LocationModel) Node {
	return Form(Method("POST"), Action("/_/room"), Class("form"),
		Div(
			Input(Type("hidden"), Name("location"), Value(fmt.Sprintf("%v", location.ID))),

			Input(Type("text"), Name("name"), Placeholder("Raumname"), Required()),

			Input(Type("submit"), Value("Neuen Raum erstellen")),
		),
	)
}

type CreateRoomRoute struct {
	DB *database.Database
}

func (l *CreateRoomRoute) Method() string {
	return http.MethodPost
}

func (l *CreateRoomRoute) Pattern() string {
	return "/room"
}

func (l *CreateRoomRoute) Handler() http.Handler {
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
			locationParam = request.PostFormValue("location")
			nameParam     = request.PostFormValue("name")
		)
		model, err := ParseCreateRoomModel(locationParam, nameParam)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}
		log.Debug("parsed create room form", zap.Any("model", model))

		id, err := commands.CreateRoom(model)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to create room", err)
			return
		}
		log.Debug("created room to event", zap.Int("id", id), zap.Int("location", model.Location))

		http.Redirect(writer, request, fmt.Sprintf("/location/edit/%v", locationParam), http.StatusSeeOther)
	})
}

func ParseCreateRoomModel(location, name string) (model.CreateRoomModel, error) {
	locationId, err := strconv.ParseInt(location, 10, 64)
	if err != nil {
		return model.CreateRoomModel{}, err
	}

	return model.CreateRoomModel{
		Location: int(locationId),
		Name:     name,
	}, nil
}
