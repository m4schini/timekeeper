package components

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/app/auth/authz"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/app/database/query"
	"raumzeitalpaka/app/export/md"
	"raumzeitalpaka/config"
	"raumzeitalpaka/ports/www/middleware"
	"raumzeitalpaka/ports/www/render"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Separator(minutes float64) Node {
	return Div(Class("separator"), ID("separator"), Text(fmt.Sprintf("In %.0f Minuten", minutes)))
}

func Day(event, day int, date time.Time, withActions bool, timeslots []model.TimeslotModel, headHref string) Node {
	t := Group{}
	now := time.Now()
	insertedSep := false
	activeTimeSlots := Group{}

	for _, timeslot := range timeslots {
		ts := timeslot.Date()
		tsDay := ts.YearDay()
		nowDay := now.YearDay()
		until := now.Sub(ts)
		active := now.After(ts) && now.Before(ts.Add(timeslot.Duration))

		tsNode := TimeSlot(timeslot, withActions, active, until > 0 && !insertedSep)

		if tsDay == nowDay {
			if active {
				activeTimeSlots = append(activeTimeSlots, tsNode)
			}
			if until <= 0 && !insertedSep {
				until = until * (-1)
				minutes := until.Minutes()

				t = append(t,
					append(activeTimeSlots, Separator(minutes)),
				)
				insertedSep = true
			}
		}
		if !active {
			t = append(t, tsNode)
		}
	}

	return Div(Class("day-container"),
		H2(A(Text(fmt.Sprintf("Tag %v (%v)", day, md.Wochentag(date.Weekday()))), Href(headHref))),
		Div(Style("display: flex; flex-direction: column; gap: 1rem"), t),
	)
}

func CompactDay(timeslots []model.TimeslotModel) Node {
	t := Group{}
	now := time.Now()
	insertedSep := false
	activeTimeSlots := Group{}

	for _, timeslot := range timeslots {
		ts := timeslot.Date()
		tsDay := ts.YearDay()
		nowDay := now.YearDay()
		until := now.Sub(ts)
		active := now.After(ts) && now.Before(ts.Add(timeslot.Duration))

		tsNode := CompactTimeSlot(timeslot, active, until > 0 && !insertedSep)

		if tsDay == nowDay {
			if active {
				activeTimeSlots = append(activeTimeSlots, tsNode)
			}
			if until <= 0 && !insertedSep {
				until = until * (-1)
				minutes := until.Minutes()

				t = append(t,
					append(activeTimeSlots, Separator(minutes)),
				)
				insertedSep = true
			}
		}
		if !active {
			t = append(t, tsNode)
		}
	}

	return Div(
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
	GetEvent            query.GetEvent
	GetTimeslotsOfEvent query.GetTimeslotsOfEvent
	Authz               authz.Authorizer
}

func (d *DayRoute) Method() string {
	return http.MethodGet
}

func (d *DayRoute) Pattern() string {
	return "/event/{event}/{day}"
}

func (d *DayRoute) Handler() http.Handler {
	log := Logger(d)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var (
			eventParam  = strings.ToLower(chi.URLParam(request, "eventId"))
			dayParam    = strings.ToLower(chi.URLParam(request, "day"))
			isOrganizer = middleware.IsOrganizer(request, d.Authz)
		)
		log.Debug("rendering day", zap.Bool("isOrganizer", isOrganizer), zap.String("eventParam", eventParam), zap.String("dayParam", dayParam))
		eventId, err := strconv.ParseInt(eventParam, 10, 64)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid eventId", err)
			return
		}
		day, err := strconv.ParseInt(dayParam, 10, 64)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid day", err)
			return
		}

		event, err := d.GetEvent.Query(query.GetEventRequest{EventId: int(eventId)})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to retrieve event", err)
			return
		}

		response, err := d.GetTimeslotsOfEvent.Query(query.GetTimeslotsOfEventRequest{
			EventId: int(eventId),
			Roles:   []model.Role{model.RoleMentor},
			Offset:  0,
			Limit:   100,
		})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to retrieve day", err)
			return
		}
		timeslots := response.Timeslots

		dayData := make([]model.TimeslotModel, 0, len(timeslots))
		for _, timeslot := range timeslots {
			if timeslot.Day == int(day) {
				dayData = append(dayData, timeslot)
			}
		}

		render.HTML(log, writer, request, Day(event.ID, int(day), event.Start, isOrganizer, dayData, "#"))
	})
}
