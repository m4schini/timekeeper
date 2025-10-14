package components

import (
	"fmt"
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

func EditLocationForm(location model.LocationModel) Node {
	return Form(Method("POST"), Action("/_/location/edit"), Class("form"),
		Div(
			Input(Type("hidden"), Name("location"), Value(fmt.Sprintf("%v", location.ID))),

			Div(
				Label(For("name"), Text("Name")),
				Input(Type("text"), Name("name"), Placeholder("Location Name"), Required(), Value(location.Name)),
			),

			//Div(
			//	Label(For("map_file"), Text("Link zu map file (Optional)")),
			//	Input(Type("text"), Name("map_file"), Placeholder("/static/betahaus2.png"), Value(location.File)),
			//),

			Div(
				Label(For("osm_id"), Text("Open Streetmap ID (Optional)")),
				Input(Type("text"), Name("osm_id"), Placeholder("N1234567"), Value(location.OsmId)),
			),

			Input(Type("submit"), Value("Speichern")),
		),
	)
}

type EditLocationRoute struct {
	DB *database.Database
}

func (l *EditLocationRoute) Method() string {
	return http.MethodPost
}

func (l *EditLocationRoute) Pattern() string {
	return "/location/edit"
}

func (l *EditLocationRoute) Handler() http.Handler {
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
			locationParam = request.PostFormValue("location")
			nameParam     = request.PostFormValue("name")
			mapFileParam  = "" //request.PostFormValue("map_file")
			osmIdParam    = request.PostFormValue("osm_id")
		)
		model, err := ParseUpdateLocationModel(locationParam, nameParam, mapFileParam, osmIdParam)
		if err != nil {
			render.RenderError(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}
		log.Debug("parsed update location form", zap.Any("model", model))

		err = commands.UpdateLocation(model)
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to update location", err)
			return
		}
		log.Debug("updated location", zap.Int("id", model.ID))

		http.Redirect(writer, request, fmt.Sprintf("/location/edit/%v", model.ID), http.StatusSeeOther)
	})
}

func ParseUpdateLocationModel(location, name, mapFile, osmId string) (model.UpdateLocationModel, error) {
	locationId, err := strconv.ParseInt(location, 10, 64)
	if err != nil {
		return model.UpdateLocationModel{}, err
	}

	return model.UpdateLocationModel{
		ID:      int(locationId),
		Name:    name,
		MapFile: mapFile,
		OsmId:   osmId,
	}, nil
}
