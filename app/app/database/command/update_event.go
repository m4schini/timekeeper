package command

import (
	"timekeeper/app/database/model"
)

func (c *Commands) UpdateEvent(m model.UpdateEventModel) (err error) {
	_, err = c.DB.Exec(`
UPDATE timekeeper.events
SET
    name = $1,
    start = $2,
    slug = $3
WHERE id = $4`, m.Name, m.Start, m.Slug, m.ID)
	return err
}
