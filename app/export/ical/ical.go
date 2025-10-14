package ical

import (
	"bytes"
	"fmt"
	"github.com/arran4/golang-ical"
	"io"
	"time"
	"timekeeper/app/database/model"
)

func ExportCalendarScheduleTo(event model.EventModel, timeslots []model.TimeslotModel, writer io.Writer) error {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)

	for _, timeslot := range timeslots {
		event := cal.AddEvent(fmt.Sprintf("%v@timekeeper", timeslot.ID))
		event.SetStartAt(timeslot.Date())
		event.SetStartAt(timeslot.Date().Add(1 * time.Hour))
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
