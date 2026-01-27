package pages

import (
	"net/http"
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/ports/www/components"
	"raumzeitalpaka/ports/www/middleware"
	"raumzeitalpaka/ports/www/render"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func EditLocationPage(locationModel model.LocationModel, rooms []model.RoomModel) Node {
	roomsList := Group{}
	for _, room := range rooms {
		roomsList = append(roomsList, Li(components.UpdateRoomForm(room)))
	}

	return components.Shell("",
		components.PageHeader(model.EventModel{}),
		Main(
			H2(Text("Location bearbeiten")),
			components.EditLocationForm(locationModel),
			H2(Text("Räume")),
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
	log := components.Logger(l)
	queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isOrganizer := middleware.IsOrganizer(request)
		if !isOrganizer {
			render.Error(log, writer, http.StatusUnauthorized, "user is not authorized", nil)
			return
		}
		var (
			locationParam   = chi.URLParam(request, "location")
			locationId, err = strconv.ParseInt(locationParam, 10, 64)
		)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid locationId", err)
			return
		}

		location, err := queries.GetLocation(int(locationId))
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get location", err)
			return
		}
		log.Debug("retrieved location", zap.Any("model", location))

		rooms, total, err := queries.GetRoomsOfLocation(location.ID, 0, 100)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get rooms of location", err)
			return
		}
		log.Debug("retrieved rooms of location", zap.Int("total", total), zap.Int("rooms", len(rooms)))

		render.HTML(log, writer, request, EditLocationPage(location, rooms))
	})
}
