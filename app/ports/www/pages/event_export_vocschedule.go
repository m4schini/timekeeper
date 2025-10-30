package pages

import (
	"net/http"
	"strconv"
	"time"
	"timekeeper/app/cache"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	export "timekeeper/app/export/voc"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/render"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type EventExportVocScheduleRoute struct {
	DB *database.Database
}

func (v *EventExportVocScheduleRoute) Method() string {
	return http.MethodGet
}

func (v *EventExportVocScheduleRoute) Pattern() string {
	return "/event/{event}/export/schedule.json"
}

func (v *EventExportVocScheduleRoute) Handler() http.Handler {
	log := components.Logger(v)
	queries := v.DB.Queries
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

		event, err := queries.GetEvent(int(eventId))
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get event", err)
			return
		}

		timeslots, _, err := queries.GetTimeslotsOfEvent(int(eventId), 0, 1000)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get timeslots of event", err)
			return
		}
		timeslots = model.FilterTimeslotRoles(timeslots, roles)

		export, err := export.ExportVocSchedule(event, timeslots)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to generate voc schedule", err)
			return
		}

		cache.Set(cacheKey, export, 60*time.Second)
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(export)
	})
}
