package ical

import (
	"fmt"
	"github.com/arran4/golang-ical"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"time"
	"timekeeper/app/database/model"
	"timekeeper/config"
)

func ExportEventCalendar(events []model.EventModel) (string, error) {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)
	cal.SetCalscale("GREGORIAN")
	cal.SetXWRTimezone(config.Timezone().String())
	cal.SetXWRCalName("Jugend hackt Events")
	cal.SetRefreshInterval("PT30M")
	cal.SetXPublishedTTL("PT30M")

	domain, err := url.Parse(config.BaseUrl())
	if err != nil {
		return "", err
	}

	for _, event := range events {
		start := event.Start
		end := event.Start.AddDate(0, 0, event.TotalDays)
		zap.L().Debug("event", zap.Any("start", start), zap.Any("end", end), zap.Int("total_days", event.TotalDays))
		now := time.Now()
		e := cal.AddEvent(fmt.Sprintf("%v@%v", event.ID, domain.Host))
		e.SetCreatedTime(now)
		e.SetDtStampTime(now)
		e.SetAllDayStartAt(start)
		e.SetAllDayEndAt(end)
		e.SetSummary(event.Name)
	}

	return cal.Serialize(ics.WithNewLineWindows), nil
}

func ExportCalendarSchedule(event model.EventModel, timeslots []model.TimeslotModel) (string, error) {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)
	cal.SetCalscale("GREGORIAN")
	cal.SetXWRCalDesc(fmt.Sprintf("Zeitplan siehe %v", event.ScheduleURL()))
	cal.SetXWRTimezone(config.Timezone().String())
	cal.SetXWRCalName(event.Name)
	cal.SetRefreshInterval("PT30M")
	cal.SetXPublishedTTL("PT30M")

	domain, err := url.Parse(config.BaseUrl())
	if err != nil {
		return "", err
	}

	for _, timeslot := range timeslots {
		now := time.Now()
		event := cal.AddEvent(fmt.Sprintf("%v@%v", timeslot.ID, domain.Host))
		event.SetCreatedTime(now)
		event.SetDtStampTime(now)
		event.SetStartAt(timeslot.Date())
		event.SetEndAt(timeslot.Date().Add(timeslot.Duration))
		event.SetLocation(timeslot.Room.Name)
		event.SetSummary(timeslot.Title)
		event.SetDescription(strings.ReplaceAll(timeslot.Note, string(ics.WithNewLineWindows), string(ics.WithNewLineUnix)))
	}

	calData := cal.Serialize(ics.WithNewLineWindows)

	return calData, nil
}
