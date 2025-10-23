package command

import (
	"timekeeper/app/database/model"
)

func (c *Commands) CreateLocation(m model.CreateLocationModel) (id int, err error) {

	row := c.DB.QueryRow(`
INSERT INTO timekeeper.locations (name, file, osm_id) 
VALUES ($1, $2, $3)
RETURNING id`, m.Name, m.MapFile, m.OsmId)
	if err = row.Err(); err != nil {
		return -1, err
	}

	err = row.Scan(&id)
	return id, err
}
