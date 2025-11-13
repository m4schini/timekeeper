package command

import (
	"raumzeitalpaka/app/database/model"
)

func (c *Commands) UpdateLocationToEvent(m model.UpdateLocationToEventModel) (err error) {
	_, err = c.DB.Exec(`
UPDATE timekeeper.event_has_location
SET
    name = $1,
    note = $2,
    visible = $3
WHERE id = $4`, m.Name, m.Note, m.Visible, m.ID)
	return err
}
