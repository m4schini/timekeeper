package command

import (
	"timekeeper/app/database/model"
)

func (c *Commands) UpdateEvent(m model.UpdateEventModel) (err error) {
	_, err = c.DB.Exec(`
UPDATE timekeeper.events
SET
    name = $1,
    start = $2
WHERE id = $3`, m.Name, m.Start, m.ID)
	return err
}
