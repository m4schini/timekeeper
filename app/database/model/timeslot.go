package model

import "time"

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
