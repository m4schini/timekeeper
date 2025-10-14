package components

import (
	"fmt"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/middleware"
	"timekeeper/ports/www/render"
)

func CreateLocationForm() Node {
	return Form(Method("POST"), Action("/_/event"), Class("form"),
		Div(
			Div(
				Label(For("name"), Text("Name")),
				Input(Type("text"), Name("name"), Placeholder("Location Name"), Required()),
			),

			Div(
				Label(For("map_file"), Text("Link zu map file (Optional)")),
				Input(Type("text"), Name("map_file"), Placeholder("/static/betahaus2.png")),
			),

			Div(
				Label(For("osm_id"), Text("Open Streetmap ID (Optional)")),
				Input(Type("text"), Name("osm_id"), Placeholder("N1234567")),
			),

			Input(Type("submit"), Value("Location erstellen")),
		),
	)
}

type CreateLocationRoute struct {
	DB *database.Database
}

func (l *CreateLocationRoute) Method() string {
	return http.MethodPost
}

func (l *CreateLocationRoute) Pattern() string {
	return "/event"
}

func (l *CreateLocationRoute) Handler() http.Handler {
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
			nameParam    = request.PostFormValue("name")
			mapFileParam = request.PostFormValue("map_file")
			osmIdParam   = request.PostFormValue("osm_id")
		)
		model, err := ParseCreateLocationModel(nameParam, mapFileParam, osmIdParam)
		if err != nil {
			render.RenderError(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}
		log.Debug("parsed create location form", zap.Any("model", model))

		id, err := commands.CreateLocation(model)
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to create location", err)
			return
		}
		log.Debug("created location", zap.Int("id", id))

		http.Redirect(writer, request, fmt.Sprintf("/location/edit/%v", id), http.StatusSeeOther)
	})
}

func ParseCreateLocationModel(name, mapFile, osmId string) (model.CreateLocationModel, error) {
	return model.CreateLocationModel{
		Name:    name,
		MapFile: mapFile,
		OsmId:   osmId,
	}, nil
}
