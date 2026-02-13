package command

import (
	"fmt"
	"raumzeitalpaka/app/database/model"
	"time"
)

var (
	ErrInvalidEventId = fmt.Errorf("invalid eventId")
	ErrInvalidRoomId  = fmt.Errorf("invalid roomId")
)

type CreateTimeslot InsertHandler[CreateTimeslotRequest, int]

type CreateTimeslotRequest struct {
	Event    int
	Parent   *int64
	Role     model.Role
	Day      int
	Timeslot time.Time
	Duration time.Duration
	Title    string
	Note     string
	Room     int
}

type CreateTimeslotHandler struct {
	DB Database
}

func (c *CreateTimeslotHandler) Execute(m CreateTimeslotRequest) (id int, err error) {
	if m.Event == 0 {
		return 0, ErrInvalidEventId
	}
	if m.Room == 0 {
		return 0, ErrInvalidRoomId
	}

	row := c.DB.QueryRow(`
INSERT INTO raumzeitalpaka.timeslots (event, parent_id, title, note, day, start, room, role, duration) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, ($9 * interval '1 second'))
RETURNING id`, m.Event, m.Parent, m.Title, m.Note, m.Day, m.Timeslot, m.Room, m.Role, int(m.Duration.Seconds()))
	if err = row.Err(); err != nil {
		return -1, err
	}

	err = row.Scan(&id)
	return id, err
}
