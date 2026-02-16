package model

import (
	"fmt"
	"raumzeitalpaka/config"
	"time"
)

type TimeslotModel struct {
	ID       int
	Rank     int
	GUID     string
	Event    EventModel
	Title    string
	Note     string
	Day      int
	Start    time.Time
	Duration time.Duration
	Room     RoomModel
	Role     Role
	Children []TimeslotModel
}

func (t *TimeslotModel) Date() time.Time {
	day := t.Event.Start
	return time.Date(day.Year(), day.Month(), day.Day()+t.Day,
		t.Start.Hour(), t.Start.Minute(), t.Start.Second(), t.Start.Nanosecond(),
		config.Timezone())
}

func FlattenTimeslots(timeslots []TimeslotModel) []TimeslotModel {
	flatTimeslots := make([]TimeslotModel, 0, len(timeslots))
	for _, timeslot := range timeslots {
		if timeslot.Children == nil || len(timeslot.Children) == 0 {
			flatTimeslots = append(flatTimeslots, timeslot)
		} else {
			for _, child := range timeslot.Children {
				flatTimeslots = append(flatTimeslots, TimeslotModel{
					ID:       child.ID,
					GUID:     child.GUID,
					Event:    child.Event,
					Title:    fmt.Sprintf("%v: %v", timeslot.Title, child.Title),
					Note:     child.Note,
					Day:      child.Day,
					Start:    child.Start,
					Duration: child.Duration,
					Room:     child.Room,
					Role:     child.Role,
					Children: child.Children,
				})
			}
		}
	}

	return flatTimeslots
}

func FilterTimeslotDay(timeslots []TimeslotModel, dayIndex int) []TimeslotModel {
	filtered := make([]TimeslotModel, 0, len(timeslots))
	for _, timeslot := range timeslots {
		if timeslot.Day == dayIndex {
			filtered = append(filtered, timeslot)
		}
	}

	return filtered
}

func MapTimeslotsToDays(timeslots []TimeslotModel) map[int][]TimeslotModel {
	eventDays := make(map[int][]TimeslotModel)
	for _, timeslot := range timeslots {
		day, ok := eventDays[timeslot.Day]
		if !ok {
			day = make([]TimeslotModel, 0)
		}

		day = append(day, timeslot)

		eventDays[timeslot.Day] = day
	}
	return eventDays
}

type CreateTimeslotModel struct {
	Event    int
	Parent   *int64
	Role     Role
	Day      int
	Timeslot time.Time
	Duration time.Duration
	Title    string
	Note     string
	Room     int
}

type UpdateTimeslotModel struct {
	ID       int
	Event    int
	Role     Role
	Day      int
	Timeslot time.Time
	Duration time.Duration
	Title    string
	Note     string
	Room     int
}
