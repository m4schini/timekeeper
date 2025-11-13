package command

import (
	"raumzeitalpaka/app/database/model"
)

func (c *Commands) CreateEvent(m model.CreateEventModel) (id int, err error) {

	row := c.DB.QueryRow(`
INSERT INTO timekeeper.events (name, slug, start) 
VALUES ($1, $2, $3)
RETURNING id`, m.Name, m.Slug, m.Start)
	if err = row.Err(); err != nil {
		return -1, err
	}

	err = row.Scan(&id)
	return id, err
}
