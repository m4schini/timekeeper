package command

import (
	"timekeeper/app/database/model"
)

func (c *Commands) AddLocationToEvent(m model.AddLocationToEventModel) (id int, err error) {
	row := c.DB.QueryRow(`
INSERT INTO timekeeper.event_has_location (name, event, location) 
VALUES ($1, $2, $3)
RETURNING id`, m.Name, m.EventId, m.LocationId)
	if err = row.Err(); err != nil {
		return -1, err
	}

	err = row.Scan(&id)
	return id, err
}
