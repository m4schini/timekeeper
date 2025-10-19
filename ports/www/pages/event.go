package pages

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"strconv"
	"time"
	"timekeeper/adapters/nominatim"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/config"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/middleware"
	"timekeeper/ports/www/render"
)

func EventPublicPage(event model.EventModel, locations []model.EventLocationModel) Node {
	return Shell(event.Name,
		components.PageHeader(event),
		Main(
			Div(Class("event-container"),
				EventSectionHeader(event),
				EventSectionSchedule(event, false),
				EventSectionPublicLocations(event, locations),
			),
		),
	)
}

func EventOrgaPage(event model.EventModel, locations []model.LocationModel, eventLocations []model.EventLocationModel) Node {
	return Shell(event.Name,
		components.PageHeader(event),
		Main(
			Div(Class("event-container"),
				Div(Style("display: flex; flex-direction: column"),
					components.CopyTextBox("copy_event", "Link zum Event", fmt.Sprintf("%v/event/%v", config.BaseUrl(), event.ID)),
					components.EditEvent(event.ID),
					EventDateRange(event.Start, event.TotalDays),
				),
				EventSectionSchedule(event, true),
				EventSectionLocations(event, locations, eventLocations),
			),
		),
	)
}

func EventDateRange(start time.Time, totalDays int) Node {
	return H3(Style("margin: 0"), Textf("%v - %v", start.Format("02.01.2006"), start.AddDate(0, 0, totalDays-1).Format("02.01.2006")))
}

func EventSectionHeader(event model.EventModel) Node {
	return Div(Style("display: flex; flex-direction: column"),
		EventDateRange(event.Start, event.TotalDays),
	)
}

func EventSectionSchedule(event model.EventModel, withCopyBoxes bool) Node {
	return Div(
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
		Iff(withCopyBoxes, func() Node {
			embedDays := Group{}
			for i := 0; i < event.TotalDays; i++ {
				embedDays = append(embedDays, Div(Style("display: flex"),
					Div(Style("max-width: 400px"), components.CopyTextBox(fmt.Sprintf("copy_embed_day_%v", i), fmt.Sprintf("Tag %v: ", i), components.IFrameCompactDay(event.ID, i))),
					Div(Style("max-width: 400px"), components.CopyTextBox(fmt.Sprintf("copy_embed_day_%v_r", i), " Für Mentor*innen: ", components.IFrameCompactDay(event.ID, i, model.RoleParticipant, model.RoleMentor))),
				))
			}

			return Group{
				Div(Style("display: flex; flex-direction: column; margin-top: 1rem"),
					Strong(Text("Links zum teilen")),
					components.CopyTextBox("copy_tn", "Link für Teilnehmer*innen", fmt.Sprintf("%v%v", config.BaseUrl(), components.UrlScheduleWithRoles(event.ID))),
					components.CopyTextBox("copy_men", "Link für Mentor*innen", fmt.Sprintf("%v%v", config.BaseUrl(), components.UrlScheduleWithRoles(event.ID, model.RoleParticipant, model.RoleMentor))),
					components.CopyTextBox("copy_org", "Link für Orga", fmt.Sprintf("%v%v", config.BaseUrl(), components.UrlScheduleWithRoles(event.ID, model.RoleParticipant, model.RoleMentor, model.RoleOrganizer))),
					components.CopyTextBox("copy_voc", "Link für VOC/Info-Beamer", fmt.Sprintf("%v%v", config.BaseUrl(), components.UrlExportVocSchedule(event.ID))),
					components.CopyTextBox("copy_ical", "Link für Calendar", fmt.Sprintf("%v%v", config.BaseUrl(), components.UrlExportIcalSchedule(event.ID))),
				),
				Div(Style("display: flex; flex-direction: column; margin-top: 1rem"),
					Strong(Text("Links zum im Pad einbetten (einfach ins pad kopieren)")),
					embedDays,
				),
			}
		}),
	)
}

func EventSectionPublicLocations(event model.EventModel, locations []model.EventLocationModel) Node {
	eventLocationsGroup := Group{}
	for _, location := range locations {
		if !location.Visible {
			continue
		}
		eventLocationsGroup = append(eventLocationsGroup, components.EventLocationCard(event, location, false))
	}

	return Div(
		H2(Text("Orte der Veranstaltung")),
		Div(Style("display: flex; flex-basis: 100%; flex-wrap: wrap; gap: 1rem"), eventLocationsGroup),
	)
}

func EventSectionLocations(event model.EventModel, locations []model.LocationModel, eventLocations []model.EventLocationModel) Node {
	eventLocationItems := Group{}
	for _, location := range eventLocations {
		eventLocationItems = append(eventLocationItems, components.EventLocationCard(event, location, true))
	}

	return Div(
		H2(Text("Orte der Veranstaltung")),
		Div(Style("display: flex; gap: 1rem; align-items: center; justify-content: space-between"),
			components.AddLocationForm(event.ID, locations),
			components.CreateLocation(),
		),
		Div(Style("display: flex; gap: 1rem; margin-top: 2rem"), eventLocationItems),
	)
}

type EventPageRoute struct {
	DB        *database.Database
	Nominatim *nominatim.Client
}

func (l *EventPageRoute) Method() string {
	return http.MethodGet
}

func (l *EventPageRoute) Pattern() string {
	return "/event/{event}"
}

func (l *EventPageRoute) Handler() http.Handler {
	log := components.Logger(l)
	queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isOrganizer := middleware.IsOrganizer(request)
		eventParam := chi.URLParam(request, "event")
		eventId, err := strconv.ParseInt(eventParam, 10, 64)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid event id", err)
			return
		}

		event, err := queries.GetEvent(int(eventId))
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get event", err)
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
				log.Warn("failed to lookup osm data", zap.Error(err), zap.String("osm_id", location.OsmId))
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

		render.HTML(log, writer, request, page)
	})
}
