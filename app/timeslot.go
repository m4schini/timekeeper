package app

import (
	"time"
	"timekeeper/app/database/model"
)

type Timeslot struct {
	ID    int
	Title string

	Day   int
	Start time.Time

	Role model.Role
}

func TimeslotFromDatabase(m model.TimeslotModel) Timeslot {
	return Timeslot{
		ID:    m.ID,
		Title: m.Title,
		Day:   m.Day,
		Start: m.Start,
		Role:  m.Role,
	}
}

type Room struct {
	ID int
}
