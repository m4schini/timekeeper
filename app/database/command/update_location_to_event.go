package command

import (
	"timekeeper/app/database/model"
)

func (c *Commands) UpdateLocationToEvent(m model.UpdateLocationToEventModel) (err error) {
	_, err = c.DB.Exec(`
UPDATE timekeeper.event_has_location
SET
    name = $1,
    note = $2
WHERE id = $3`, m.Name, m.Note, m.ID)
	return err
}
