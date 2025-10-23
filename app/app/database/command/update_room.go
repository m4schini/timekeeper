package command

import (
	"timekeeper/app/database/model"
)

func (c *Commands) UpdateRoom(m model.UpdateRoomModel) (err error) {
	_, err = c.DB.Exec(`
UPDATE timekeeper.rooms
SET
    name = $1,
    description = $2
WHERE id = $3`, m.Name, m.Description, m.ID)
	return err
}
