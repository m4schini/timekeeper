package ical

import (
	"fmt"
	"github.com/arran4/golang-ical"
	"net/url"
	"strings"
	"time"
	"timekeeper/app/database/model"
	"timekeeper/config"
)

func ExportCalendarSchedule(event model.EventModel, timeslots []model.TimeslotModel) (string, error) {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetCalscale("GREGORIAN")
	cal.SetXWRCalDesc(fmt.Sprintf("Zeitplan siehe %v", event.ScheduleURL()))
	cal.SetXWRTimezone(config.Timezone().String())
	cal.SetXWRCalName(event.Name)
	cal.SetRefreshInterval("PT5M")
	cal.SetXPublishedTTL("PT5M")

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
