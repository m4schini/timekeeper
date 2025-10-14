package pages

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"strconv"
	"timekeeper/adapters"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/config"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/middleware"
	. "timekeeper/ports/www/render"
)

func EventPublicPage(event model.EventModel, locations []model.EventLocationModel) Node {
	eventLocationsGroup := Group{}
	for _, location := range locations {
		eventLocationsGroup = append(eventLocationsGroup, components.EventLocationCard(event, location, false))
		//eventLocationsGroup = append(eventLocationsGroup, Li(Textf("%v: %v (%v)", location.Relationship, location.Name, location.Address.City)))
	}

	return Shell(
		components.PageHeader(event),
		Main(
			Div(Class("event-container"),
				Div(
					H2(Text("Zeitplan")),
					Div(Style("display: flex; justify-content: space-between"),
						Div(
							components.EventSchedule(event.ID),
						),
						Div(
							Span(Text("Export:"), Style("margin-right: 1rem")),
							components.ExportEventMarkdownButton(event.ID),
							components.ExportEventVocScheduleButton(event.ID),
						),
					),
				),
				Div(
					H2(Text("Orte")),
					Div(Style("display: flex; gap: 1rem"), eventLocationsGroup),
				),
			),
		),
	)
}

func EventOrgaPage(event model.EventModel, locations []model.LocationModel, eventLocations []model.EventLocationModel) Node {
	eventLocationItems := Group{}
	for _, location := range eventLocations {
		eventLocationItems = append(eventLocationItems, components.EventLocationCard(event, location, true))

		//eventLocationItems = append(eventLocationItems, Li(
		//	Textf("%v: %v ", location.Relationship, location.Name),
		//	Title(fmt.Sprintf("%v (id=%v, relationship=%v)", location.Name, location.ID, location.RelationshipId)),
		//	components.DeleteEventLocationButton(event.ID, location.RelationshipId)))
	}

	embedDays := Group{}
	for i := 0; i < event.TotalDays; i++ {
		embedDays = append(embedDays, Div(Style("display: flex"),
			Div(Style("max-width: 400px"), components.CopyTextBox(fmt.Sprintf("copy_embed_day_%v", i), fmt.Sprintf("Tag %v: ", i), components.IFrameCompactDay(event.ID, i))),
			Div(Style("max-width: 400px"), components.CopyTextBox(fmt.Sprintf("copy_embed_day_%v_r", i), " Für Mentor*innen: ", components.IFrameCompactDay(event.ID, i, model.RoleParticipant, model.RoleMentor))),
		))
	}

	return Shell(
		components.PageHeader(event),
		Main(
			Div(Class("event-container"),
				Div(Style("display: flex; flex-direction: column"),
					components.CopyTextBox("copy_event", "Link zum Event", fmt.Sprintf("%v/event/%v", config.BaseUrl(), event.ID)),
					components.EditEvent(event.ID),
				),

				Div(
					H2(Text("Zeitplan")),
					Div(Style("display: flex; justify-content: space-between"),
						Div(
							components.EventSchedule(event.ID),
						),
						Div(
							Span(Text("Export:"), Style("margin-right: 1rem")),
							components.ExportEventMarkdownButton(event.ID),
							components.ExportEventIcalScheduleButton(event.ID),
							components.ExportEventVocScheduleButton(event.ID),
						),
					),
					Div(Style("display: flex; flex-direction: column; margin-top: 1rem"),
						Strong(Text("Links zum teilen")),
						components.CopyTextBox("copy_tn", "Link für Teilnehmer*innen", fmt.Sprintf("%v%v", config.BaseUrl(), components.UrlScheduleWithRoles(event.ID))),
						components.CopyTextBox("copy_men", "Link für Mentor*innen", fmt.Sprintf("%v%v", config.BaseUrl(), components.UrlScheduleWithRoles(event.ID, model.RoleParticipant, model.RoleMentor))),
						components.CopyTextBox("copy_org", "Link für Orga", fmt.Sprintf("%v%v", config.BaseUrl(), components.UrlScheduleWithRoles(event.ID, model.RoleParticipant, model.RoleMentor, model.RoleOrganizer))),
						components.CopyTextBox("copy_voc", "Link für VOC/Info-Beamer", fmt.Sprintf("%v%v", config.BaseUrl(), components.UrlExportVocSchedule(event.ID))),
						components.CopyTextBox("copy_voc", "Link für Calendar", fmt.Sprintf("%v%v", config.BaseUrl(), components.UrlExportIcalSchedule(event.ID))),
					),
					Div(Style("display: flex; flex-direction: column; margin-top: 1rem"),
						Strong(Text("Links zum im Pad einbetten (einfach ins pad kopieren)")),
						embedDays,
					),
				),
				Div(
					H2(Text("Ort der Veranstaltung")),
					Div(Style("display: flex; gap: 1rem; align-items: center; justify-content: space-between"),
						components.AddLocationForm(event.ID, locations),
						components.CreateLocation(),
					),
					Div(Style("display: flex; gap: 1rem; margin-top: 2rem"), eventLocationItems),
				),
			),
		),
	)
}

type EventPageRoute struct {
	DB        *database.Database
	Nominatim *adapters.NominatimClient
}

func (l *EventPageRoute) Method() string {
	return http.MethodGet
}

func (l *EventPageRoute) Pattern() string {
	return "/event/{event}"
}

func (l *EventPageRoute) Handler() http.Handler {
	log := zap.L().Named("www").Named("event")
	queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isOrganizer := middleware.IsOrganizer(request)
		eventParam := chi.URLParam(request, "event")
		eventId, err := strconv.ParseInt(eventParam, 10, 64)
		if err != nil {
			RenderError(log, writer, http.StatusBadRequest, "invalid event id", err)
			return
		}

		event, err := queries.GetEvent(int(eventId))
		if err != nil {
			RenderError(log, writer, http.StatusInternalServerError, "failed to get event", err)
			return
		}

		eventLocations, err := queries.GetLocationsOfEvent(int(eventId))
		if err != nil {
			log.Warn("failed to get event locations", zap.Error(err))
			eventLocations = make([]model.EventLocationModel, 0)
		}

		for i, location := range eventLocations {
			resp, err := l.Nominatim.Lookup(request.Context(), location.OsmId)
			if err != nil {
				log.Warn("failed to lookup osm data", zap.Error(err))
				continue
			}

			location.Address = &resp.Address
			eventLocations[i] = location
		}

		var page Node
		if isOrganizer {
			locations, err := queries.GetLocations(0, 100)
			if err != nil {
				log.Warn("failed to get locations", zap.Error(err))
				locations = make([]model.LocationModel, 0)
			}

			page = EventOrgaPage(event, locations, eventLocations)
		} else {
			page = EventPublicPage(event, eventLocations)
		}

		Render(log, writer, request, page)
	})
}
