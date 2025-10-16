package pages

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"strconv"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/middleware"
	. "timekeeper/ports/www/render"
)

func EditLocationPage(locationModel model.LocationModel, rooms []model.RoomModel) Node {
	roomsList := Group{}
	for _, room := range rooms {
		roomsList = append(roomsList, Li(components.UpdateRoomForm(room),
			components.DeleteRoomButton(room.ID),
		))
	}

	return Shell(
		components.PageHeader(model.EventModel{}),
		Main(
			Div(Text("Location bearbeiten")),
			components.EditLocationForm(locationModel),
			H3(Text("Räume")),
			Ul(Class("rooms"), roomsList),
			components.CreateRoomForm(locationModel),
		),
	)
}

type UpdateLocationPageRoute struct {
	DB *database.Database
}

func (l *UpdateLocationPageRoute) Method() string {
	return http.MethodGet
}

func (l *UpdateLocationPageRoute) Pattern() string {
	return "/location/edit/{location}"
}

func (l *UpdateLocationPageRoute) Handler() http.Handler {
	log := zap.L().Named("www").Named("event")
	queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isOrganizer := middleware.IsOrganizer(request)
		if !isOrganizer {
			RenderError(log, writer, http.StatusUnauthorized, "user is not authorized", nil)
			return
		}
		var (
			locationParam   = chi.URLParam(request, "location")
			locationId, err = strconv.ParseInt(locationParam, 10, 64)
		)
		if err != nil {
			RenderError(log, writer, http.StatusBadRequest, "invalid locationId", err)
			return
		}

		location, err := queries.GetLocation(int(locationId))
		if err != nil {
			RenderError(log, writer, http.StatusInternalServerError, "failed to get location", err)
			return
		}
		log.Debug("retrieved location", zap.Any("model", location))

		rooms, total, err := queries.GetRoomsOfLocation(location.ID, 0, 100)
		if err != nil {
			RenderError(log, writer, http.StatusInternalServerError, "failed to get rooms of location", err)
			return
		}
		log.Debug("retrieved rooms of location", zap.Int("total", total), zap.Int("rooms", len(rooms)))

		Render(log, writer, request, EditLocationPage(location, rooms))
	})
}
