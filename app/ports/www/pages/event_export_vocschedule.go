package pages

import (
	"net/http"
	"raumzeitalpaka/app/cache"
	"raumzeitalpaka/app/database/query"
	export "raumzeitalpaka/app/export/voc"
	"raumzeitalpaka/ports/www/components"
	"raumzeitalpaka/ports/www/render"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type EventExportVocScheduleRoute struct {
	GetEvent            query.GetEvent
	GetTimeslotsOfEvent query.GetTimeslotsOfEvent
}

func (v *EventExportVocScheduleRoute) Method() string {
	return http.MethodGet
}

func (v *EventExportVocScheduleRoute) Pattern() string {
	return "/event/{event}/export/schedule.json"
}

func (v *EventExportVocScheduleRoute) Handler() http.Handler {
	log := components.Logger(v)
	cache := cache.NewInMemory()

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		eventParam := chi.URLParam(request, "event")
		eventId, err := strconv.ParseInt(eventParam, 10, 64)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid event id", err)
			return
		}
		roles, _ := ParseRolesQuery(request.URL.Query(), false)

		cacheKey := cacheKey(eventId, roles)
		cachedExport, expiresAt, valid := cache.Get(cacheKey)
		if valid {
			log.Debug("using cached schedule.json export", zap.Int64("event", eventId), zap.Any("roles", roles), zap.Duration("ttl", expiresAt.Sub(time.Now())))
			writer.Header().Set("Content-Type", "application/json")
			writer.Write(cachedExport)
			return
		}

		event, err := v.GetEvent.Query(query.GetEventRequest{EventId: int(eventId)})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get event", err)
			return
		}

		timeslots, err := v.GetTimeslotsOfEvent.Query(query.GetTimeslotsOfEventRequest{
			EventId: int(eventId),
			Roles:   roles,
			Offset:  0,
			Limit:   1000,
		})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get timeslots of event", err)
			return
		}

		export, err := export.ExportVocSchedule(event, timeslots.Timeslots)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to generate voc schedule", err)
			return
		}

		cache.Set(cacheKey, export, 60*time.Second)
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(export)
	})
}
