package pages

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
	"timekeeper/app/cache"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	export "timekeeper/app/export/ical"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/render"
)

type EventExportIcalScheduleRoute struct {
	DB *database.Database
}

func (v *EventExportIcalScheduleRoute) Method() string {
	return http.MethodGet
}

func (v *EventExportIcalScheduleRoute) Pattern() string {
	return "/event/{event}/export/schedule.ics"
}

func (v *EventExportIcalScheduleRoute) Handler() http.Handler {
	log := components.Logger(v)
	queries := v.DB.Queries

	log.Debug("starting cache")
	calcache := cache.NewInMemory()
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		eventParam := chi.URLParam(request, "event")
		eventId, err := strconv.ParseInt(eventParam, 10, 64)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid event id", err)
			return
		}
		roles, _ := ParseRolesQuery(request.URL.Query(), false)

		cacheKey := cacheKey(eventId, roles)
		cachedCal, expiresAt, valid := calcache.Get(cacheKey)
		if valid {
			log.Debug("using cached calendar export", zap.Int64("event", eventId), zap.Any("roles", roles), zap.Duration("ttl", expiresAt.Sub(time.Now())))
			writer.Header().Set("Content-Type", "text/calendar; charset=utf-8")
			writer.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=timekeeper_event_%v.ics", eventId))
			writer.Write([]byte(cachedCal))
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

		cal, err := export.ExportCalendarSchedule(event, timeslots)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to generate ical schedule", err)
			return
		}

		calcache.Set(cacheKey, cal, 5*time.Minute)
		writer.Header().Set("Content-Type", "text/calendar; charset=utf-8")
		writer.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=timekeeper_event_%v.ics", event.ID))
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
