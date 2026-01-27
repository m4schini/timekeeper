package command

import (
	"raumzeitalpaka/app/database/model"
)

func (c *Commands) UpdateTimeslot(m model.UpdateTimeslotModel) (err error) {
	if m.Event == 0 {
		return ErrInvalidEventId
	}
	if m.Room == 0 {
		return ErrInvalidRoomId
	}

	_, err = c.DB.Exec(`
UPDATE raumzeitalpaka.timeslots
SET
    event = $1,
    title = $2,
    note = $3,
    day = $4,
    start = $5,
    room = $6,
    role = $7,
    duration = ($8 * interval '1 second')
WHERE id = $9`, m.Event, m.Title, m.Note, m.Day, m.Timeslot, m.Room, m.Role, int(m.Duration.Seconds()), m.ID)
	return err
}
