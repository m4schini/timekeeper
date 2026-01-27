package command

import (
	"raumzeitalpaka/app/database/model"
)

func (c *Commands) CreateRoom(m model.CreateRoomModel) (id int, err error) {
	row := c.DB.QueryRow(`
INSERT INTO raumzeitalpaka.rooms (location, name, description, location_x, location_y, location_w, location_h) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id`, m.Location, m.Name, m.Description, 0, 0, 0, 0)
	if err = row.Err(); err != nil {
		return -1, err
	}

	err = row.Scan(&id)
	return id, err
}
