package pages

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"strconv"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/config"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/middleware"
	. "timekeeper/ports/www/render"
)

func EventPublicPage(event model.EventModel) Node {
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
				),
			),
		),
	)
}

func EventOrgaPage(event model.EventModel) Node {
	return Shell(
		components.PageHeader(event),
		Main(
			Div(Class("event-container"),
				Div(
					H2(Text("Zeitplan")),
					components.EventSchedule(event.ID),
					Div(Style("display: flex; flex-direction: column; margin-top: 1rem"),
						Strong(Text("Links zum teilen")),
						components.CopyTextBox("copy_tn", "Link für Teilnehmer*innen", fmt.Sprintf("%v%v", config.BaseUrl(), components.UrlScheduleWithRoles(event.ID))),
						components.CopyTextBox("copy_men", "Link für Mentor*innen", fmt.Sprintf("%v%v", config.BaseUrl(), components.UrlScheduleWithRoles(event.ID, model.RoleParticipant, model.RoleMentor))),
						components.CopyTextBox("copy_org", "Link für Orga", fmt.Sprintf("%v%v", config.BaseUrl(), components.UrlScheduleWithRoles(event.ID, model.RoleParticipant, model.RoleMentor, model.RoleOrganizer))),
						components.CopyTextBox("copy_voc", "Link für VOC/Info-Beamer", fmt.Sprintf("%v%v", config.BaseUrl(), components.UrlExportVocSchedule(event.ID))),
					),
				),
				Div(
					H2(Text("Ort der Veranstaltung")),
					Div(Style("display: flex; gap: 1rem; align-items: center; justify-content: space-between"),
						Label(For("location"), Text("Ort")),
						Select(Name("location"),
							Option(Text("Betahaus | Schanze")),
							Option(Text("Pyjama Park")),
							Option(Text("Theater an der Parkaue")),
						),
						Label(For("relationship"), Text("Wofür?")),
						Select(Name("relationship"),
							Option(Text("Event Location")),
							Option(Text("Übernachtung")),
						),
						components.AButton(components.ColorDefault, "#", "Location hinzufügen"),
						components.AButton(components.ColorSoftGrey, "#", "Neue Location erstellen"),
					),
				),
			),
		),
	)
}

type EventPageRoute struct {
	DB *database.Database
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

		var page Node
		if isOrganizer {
			page = EventOrgaPage(event)
		} else {
			page = EventPublicPage(event)
		}

		Render(log, writer, request, page)
	})
}
