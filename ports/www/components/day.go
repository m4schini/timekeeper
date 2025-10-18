package components

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"strconv"
	"strings"
	"time"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/app/export/md"
	"timekeeper/config"
	"timekeeper/ports/www/middleware"
	"timekeeper/ports/www/render"
)

func Day(event, day int, date time.Time, withActions bool, timeslots []model.TimeslotModel) Node {
	t := Group{}
	now := time.Now()
	insertedSep := false
	for _, timeslot := range timeslots {
		ts := timeslot.Date()
		tsDay := ts.YearDay()
		nowDay := now.YearDay()
		until := now.Sub(ts)
		active := now.After(ts) && now.Before(ts.Add(timeslot.Duration))
		if tsDay == nowDay {
			if until <= 0 && !insertedSep {
				until = until * (-1)
				minutes := until.Minutes()

				t = append(t, Div(Class("separator"), ID("separator"), Text(fmt.Sprintf("In %.0f Minuten", minutes))))
				insertedSep = true
			}
		}
		t = append(t, TimeSlot(timeslot, withActions, active, until > 0 && !insertedSep))
	}

	return Div(Class("day-container"), //hx.Get("/_/day/"+day), hx.Trigger("load delay:60s"), hx.Swap("outerHTML"),
		H2(A(Text(fmt.Sprintf("Tag %v (%v)", day, md.Wochentag(date.Weekday()))), Href(fmt.Sprintf("/event/%v/schedule/%v?compact", event, day)))),
		//If(withActions, Div(Style("display: flex; gap: 1rem"), CreateTimeslotButton(event), ExportEventDayMarkdownButton(event, day))),
		Div(Style("display: flex; flex-direction: column; gap: 1rem"), t),
	)
}

func FullDay(day int, date time.Time, timeslots []model.TimeslotModel) Node {
	t := Group{}
	now := time.Now()
	insertedSep := false
	for _, timeslot := range timeslots {
		ts := timeslot.Date()
		tsDay := ts.YearDay()
		nowDay := now.YearDay()
		until := now.Sub(ts)
		if tsDay == nowDay {
			if until <= 0 && !insertedSep {
				until = until * (-1)
				minutes := until.Minutes()

				t = append(t, Div(Class("separator"), ID("separator"), Text(fmt.Sprintf("In %.0f Minuten", minutes))))
				insertedSep = true
			}
		}
		t = append(t, FullTimeSlot(timeslot, until > 0 && !insertedSep))
	}

	return Div( //hx.Get("/_/day/"+day), hx.Trigger("load delay:60s"), hx.Swap("outerHTML"),
		H2(Text(fmt.Sprintf("Tag %v (%v)", day, date.Weekday()))),
		Div(Style("display: flex; flex-direction: column; gap: 1rem"),
			t,
		),
	)
}

func CompactDay(timeslots []model.TimeslotModel) Node {
	//log := zap.L().Named("day")
	t := Group{}
	now := time.Now()
	insertedSep := false
	for _, timeslot := range timeslots {
		ts := timeslot.Date()
		tsDay := ts.YearDay()
		nowDay := now.YearDay()
		until := now.Sub(ts)
		active := now.After(ts) && now.Before(ts.Add(timeslot.Duration))
		if tsDay == nowDay {
			if until <= 0 && !insertedSep {
				until = until * (-1)
				minutes := until.Minutes()

				t = append(t, Div(Class("separator"), ID("separator"), Text(fmt.Sprintf("In %.0f Minuten", minutes))))
				insertedSep = true
			}
		}
		t = append(t, CompactTimeSlot(timeslot, active, until > 0 && !insertedSep))
	}

	return Div( //hx.Get("/_/day/"+day), hx.Trigger("load delay:60s"), hx.Swap("outerHTML"),
		Div(Style("display: flex; flex-direction: column; gap: 1rem"),
			t,
		),
		P(Textf("Generated: %v", NowInTimezone(config.Timezone()).Format(time.RFC822)), Style("opacity: 0.5; font-size: x-small")),
	)
}

func NowInTimezone(location *time.Location) time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), location)
}

type DayRoute struct {
	DB *database.Database
}

func (d *DayRoute) Method() string {
	return http.MethodGet
}

func (d *DayRoute) Pattern() string {
	return "/event/{event}/{day}"
}

func (d *DayRoute) UseCache() bool {
	return false
}

func (d *DayRoute) Handler() http.Handler {
	log := Logger(d)
	queries := d.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var (
			eventParam  = strings.ToLower(chi.URLParam(request, "eventId"))
			dayParam    = strings.ToLower(chi.URLParam(request, "day"))
			isOrganizer = middleware.IsOrganizer(request)
		)
		log.Debug("rendering day", zap.Bool("isOrganizer", isOrganizer), zap.String("eventParam", eventParam), zap.String("dayParam", dayParam))
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

		event, err := queries.GetEvent(int(eventId))
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to retrieve event", err)
			return
		}

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

		render.Render(log, writer, request, Day(event.ID, int(day), event.Start, isOrganizer, dayData))
	})
}
