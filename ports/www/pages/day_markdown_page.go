package pages

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	export "timekeeper/app/export/md"
	"timekeeper/ports/www/render"
)

type DayMarkdownPageRoute struct {
	DB *database.Database
}

func (l *DayMarkdownPageRoute) Method() string {
	return http.MethodGet
}

func (l *DayMarkdownPageRoute) Pattern() string {
	return "/event/{event}/{day}/export/markdown"
}

func (l *DayMarkdownPageRoute) Handler() http.Handler {
	log := zap.L().Named(l.Pattern())
	queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var (
			eventParam = strings.ToLower(chi.URLParam(request, "event"))
			dayParam   = strings.ToLower(chi.URLParam(request, "day"))
		)
		log.Debug("rendering day", zap.String("eventParam", eventParam), zap.String("dayParam", dayParam))
		eventId, err := strconv.ParseInt(eventParam, 10, 64)
		if err != nil {
			render.RenderError(log, writer, http.StatusBadRequest, "invalid eventId", err)
			return
		}
		day, err := strconv.ParseInt(dayParam, 10, 64)
		if err != nil {
			render.RenderError(log, writer, http.StatusBadRequest, "invalid day", err)
			return
		}

		//event, err := queries.GetEvent(int(eventId))
		//if err != nil {
		//	render.RenderError(log, writer, http.StatusInternalServerError, "failed to get event", err)
		//	return
		//}

		timeslots, _, err := queries.GetTimeslotsOfEvent(int(eventId), 0, 100)
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to retrieve day", err)
			return
		}

		dayData := make([]model.TimeslotModel, 0, len(timeslots))
		for _, timeslot := range timeslots {
			if timeslot.Day == int(day) {
				dayData = append(dayData, timeslot)
			}
		}

		md, err := export.ExportMarkdownTimetable(dayData)
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to render markdown", err)
			return
		}

		writer.Header().Set("Content-Type", "text/markdown; charset=utf8")
		writer.Write([]byte(md))
	})
}
