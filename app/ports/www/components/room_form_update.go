package components

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/app/auth/authz"
	"raumzeitalpaka/app/database/command"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/ports/www/middleware"
	"raumzeitalpaka/ports/www/render"
	"strconv"

	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func UpdateRoomForm(room model.RoomModel) Node {
	return Form(Method("POST"), Action("/_/room/edit"), Class("form"),
		Input(Type("hidden"), Name("room"), Value(fmt.Sprintf("%v", room.ID))),
		Input(Type("hidden"), Name("location"), Value(fmt.Sprintf("%v", room.Location.ID))),
		Input(Type("text"), Name("room_name"), Value(room.Name)),
		Textarea(Name("description"), Placeholder("Raumbeschreibung (wo ist der Raum?)"), Text(room.Description)),
		Div(Style("display: flex"),
			Input(Type("submit"), Value("Speichern")),
			DeleteRoomButton(room.ID),
		),
	)
}

type UpdateRoomRoute struct {
	UpdateRoom command.UpdateRoom
	Authz      authz.Authorizer
}

func (l *UpdateRoomRoute) Method() string {
	return http.MethodPost
}

func (l *UpdateRoomRoute) Pattern() string {
	return "/room/edit"
}

func (l *UpdateRoomRoute) Handler() http.Handler {
	log := Logger(l)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !middleware.IsOrganizer(request, l.Authz) {
			render.Error(log, writer, http.StatusUnauthorized, "unauthorized request detected", nil)
			return
		}

		err := request.ParseForm()
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}

		var (
			roomParam        = request.PostFormValue("room")
			locationParam    = request.PostFormValue("location")
			nameParam        = request.PostFormValue("room_name")
			descriptionParam = request.PostFormValue("description")
		)
		model, err := ParseEditRoomModel(roomParam, nameParam, descriptionParam)
		if err != nil {
			log.Sugar().Debug("failed to parse", roomParam, locationParam, nameParam, descriptionParam)
			render.Error(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}
		log.Debug("parsed update room form", zap.Any("model", model))

		err = l.UpdateRoom.Execute(model)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to update room", err)
			return
		}
		log.Debug("updated room", zap.String("id", roomParam))

		http.Redirect(writer, request, fmt.Sprintf("/location/edit/%v", locationParam), http.StatusSeeOther)
	})
}

func ParseEditRoomModel(room, name, description string) (command.UpdateRoomRequest, error) {
	roomId, err := strconv.ParseInt(room, 10, 64)
	if err != nil {
		return command.UpdateRoomRequest{}, err
	}

	return command.UpdateRoomRequest{
		ID:          int(roomId),
		Name:        name,
		Description: description,
	}, nil
}
