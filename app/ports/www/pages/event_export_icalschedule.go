package pages

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/app/cache"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/app/database/query"
	export "raumzeitalpaka/app/export/ical"
	"raumzeitalpaka/ports/www/components"
	"raumzeitalpaka/ports/www/render"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type EventExportIcalScheduleRoute struct {
	GetEvent            query.GetEvent
	GetTimeslotsOfEvent query.GetTimeslotsOfEvent
}

func (v *EventExportIcalScheduleRoute) Method() string {
	return http.MethodGet
}

func (v *EventExportIcalScheduleRoute) Pattern() string {
	return "/event/{event}/export/schedule.ics"
}

func (v *EventExportIcalScheduleRoute) Handler() http.Handler {
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
			log.Debug("using cached calendar export", zap.Int64("event", eventId), zap.Any("roles", roles), zap.Duration("ttl", expiresAt.Sub(time.Now())))
			writer.Header().Set("Content-Type", "text/calendar; charset=utf-8")
			writer.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=rza_event_%v.ics", eventId))
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

		cal, err := export.ExportCalendarSchedule(event, timeslots.Timeslots)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to generate ical schedule", err)
			return
		}

		cache.Set(cacheKey, []byte(cal), 5*time.Minute)
		writer.Header().Set("Content-Type", "text/calendar; charset=utf-8")
		writer.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=raumzeitalpaka_event_%v.ics", event.ID))
		writer.Write([]byte(cal))
	})
}

func cacheKey(eventId int64, roles []model.Role) string {
	var (
		orga        bool
		mentor      bool
		participant bool
	)
	for _, role := range roles {
		switch role {
		case model.RoleOrganizer:
			orga = true
			break
		case model.RoleMentor:
			mentor = true
			break
		case model.RoleParticipant:
			participant = true
			break
		}
	}

	return fmt.Sprintf("%v:%v:%v:%v", eventId, orga, mentor, participant)
}
