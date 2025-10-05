package command

import (
	"fmt"
	"timekeeper/app/database/model"
)

var (
	ErrInvalidEventId = fmt.Errorf("invalid eventId")
	ErrInvalidRoomId  = fmt.Errorf("invalid roomId")
)

func (c *Commands) CreateTimeslot(m model.CreateTimeslotModel) (id int, err error) {
	if m.Event == 0 {
		return 0, ErrInvalidEventId
	}
	if m.Room == 0 {
		return 0, ErrInvalidRoomId
	}

	row := c.DB.QueryRow(`
INSERT INTO timekeeper.timeslots (event, title, note, day, start, room) 
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id`, m.Event, m.Title, m.Note, m.Day, m.Timeslot, m.Room)
	if err = row.Err(); err != nil {
		return -1, err
	}

	err = row.Scan(&id)
	return id, err
}
