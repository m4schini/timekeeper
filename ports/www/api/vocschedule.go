package api

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/app/export"
	"timekeeper/ports/www/render"
)

type VocScheduleRoute struct {
	DB *database.Database
}

func (v *VocScheduleRoute) Method() string {
	return http.MethodGet
}

func (v *VocScheduleRoute) Pattern() string {
	return "/event/{event}"
}

func (v *VocScheduleRoute) Handler() http.Handler {
	queries := v.DB.Queries
	log := zap.L().Named("api")
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		eventParam := chi.URLParam(request, "event")
		eventId, err := strconv.ParseInt(eventParam, 10, 64)
		if err != nil {
			render.RenderError(log, writer, http.StatusBadRequest, "invalid event id", err)
			return
		}
		hasFilter := request.URL.Query().Has("role")
		roles := strings.Split(request.URL.Query().Get("role"), ",")

		if !hasFilter {
			roles = []string{string(model.RoleParticipant)}
		}

		event, err := queries.GetEvent(int(eventId))
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to get event", err)
			return
		}

		timeslots, _, err := queries.GetTimeslotsOfEvent(int(eventId), 0, 1000)
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to get timeslots of event", err)
			return
		}

		filterRoles := make([]model.Role, len(roles))
		for i, role := range roles {
			filterRoles[i] = model.RoleFrom(role)
		}

		eventData := make([]model.TimeslotModel, 0, len(timeslots))
		for _, timeslot := range timeslots {
			for _, role := range filterRoles {
				if timeslot.Role == role {
					eventData = append(eventData, timeslot)
					break
				}
			}
		}

		scheduleJson, err := export.ExportVocSchedule(event, eventData)
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to generate voc schedule", err)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(scheduleJson)
	})
}
