package model

import (
	"time"
	"timekeeper/config"
)

type TimeslotModel struct {
	ID    int
	Event EventModel
	Title string
	Note  string
	Day   int
	Start time.Time
	Room  RoomModel
	Role  Role
}

func (t *TimeslotModel) Date() time.Time {
	day := t.Event.Start
	return time.Date(day.Year(), day.Month(), day.Day()+t.Day,
		t.Start.Hour(), t.Start.Minute(), t.Start.Second(), t.Start.Nanosecond(),
		config.Timezone())
}

type CreateTimeslotModel struct {
	Event    int
	Role     Role
	Day      int
	Timeslot time.Time
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
	Title    string
	Note     string
	Room     int
}
