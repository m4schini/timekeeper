package ical

import (
	"fmt"
	"net/url"
	"strings"
	"time"
	"timekeeper/app/database/model"
	"timekeeper/config"

	"github.com/arran4/golang-ical"
	"go.uber.org/zap"
)

func ExportEventCalendar(events []model.EventModel) (string, error) {
	log := zap.L().Named("export").Named("ical")

	log.Debug("exporting events calendar")
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
		now := time.Now()
		e := cal.AddEvent(fmt.Sprintf("%v@%v", event.ID, domain.Host))
		e.SetCreatedTime(now)
		e.SetDtStampTime(now)
		e.SetAllDayStartAt(start)
		e.SetAllDayEndAt(end)
		e.SetSummary(event.Name)
	}

	log.Debug("exported events calendar")
	return cal.Serialize(ics.WithNewLineWindows), nil
}

func ExportCalendarSchedule(event model.EventModel, timeslots []model.TimeslotModel) (string, error) {
	log := zap.L().Named("export").Named("ical").With(zap.Int("event", event.ID), zap.Int("timeslots_count", len(timeslots)))

	timeslots = model.FlattenTimeslots(timeslots)

	log.Debug("exporting schedule as calendar (ical)")
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
		timeslot.Note = config.PixelHackPlaceholderRx.ReplaceAllString(timeslot.Note, "")
		timeslot.Title = config.PixelHackPlaceholderRx.ReplaceAllString(timeslot.Title, "")

		now := time.Now()
		event := cal.AddEvent(fmt.Sprintf("%v@%v", timeslot.ID, domain.Host))
		event.SetCreatedTime(now)
		event.SetDtStampTime(now)
		event.SetStartAt(timeslot.Date())
		event.SetEndAt(timeslot.Date().Add(timeslot.Duration))
		event.SetLocation(fmt.Sprintf("%v: %v", timeslot.Room.Location.Name, timeslot.Room.Name))
		event.SetSummary(timeslot.Title)
		event.SetDescription(strings.ReplaceAll(timeslot.Note, string(ics.WithNewLineWindows), string(ics.WithNewLineUnix)))
	}

	log.Debug("exported schedule as calendar (ical)")
	calData := cal.Serialize(ics.WithNewLineWindows)
	return calData, nil
}
