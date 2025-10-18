package ical

import (
	"bytes"
	"fmt"
	"github.com/arran4/golang-ical"
	"io"
	"net/url"
	"time"
	"timekeeper/app/database/model"
	"timekeeper/config"
)

func ExportCalendarScheduleTo(event model.EventModel, timeslots []model.TimeslotModel, writer io.Writer) error {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetCalscale("GREGORIAN")
	cal.SetXWRTimezone(config.Timezone().String())
	cal.SetName(event.Name)
	cal.SetUrl(event.ScheduleURL())
	cal.SetRefreshInterval("PT5M")
	cal.SetXPublishedTTL("PT5M")

	domain, err := url.Parse(config.BaseUrl())
	if err != nil {
		return err
	}

	for _, timeslot := range timeslots {
		event := cal.AddEvent(fmt.Sprintf("%v@%v", timeslot.ID, domain.Host))
		event.SetCreatedTime(time.Now())
		event.SetStartAt(timeslot.Date())
		event.SetEndAt(timeslot.Date().Add(timeslot.Duration))
		event.SetLocation(timeslot.Room.Name)
		event.SetSummary(timeslot.Title)
		event.SetDescription(timeslot.Note)
	}

	return cal.SerializeTo(writer)
}

func ExportCalendarSchedule(event model.EventModel, timeslots []model.TimeslotModel) (string, error) {
	var buf bytes.Buffer
	err := ExportCalendarScheduleTo(event, timeslots, &buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
