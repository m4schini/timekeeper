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
	. "timekeeper/ports/www/render"
)

func EditLocationPage(locationModel model.LocationModel) Node {
	return Shell(
		components.PageHeader(model.EventModel{}),
		Main(
			Div(Text("Location bearbeiten")),
			components.EditLocationForm(locationModel),
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

		Render(log, writer, request, EditLocationPage(location))
	})
}
