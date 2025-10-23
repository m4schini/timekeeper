package command

import (
	"timekeeper/app/database/model"
)

func (c *Commands) UpdateLocation(m model.UpdateLocationModel) (err error) {
	_, err = c.DB.Exec(`
UPDATE timekeeper.locations
SET
    name = $1,
    file = $2,
    osm_id = $3
WHERE id = $4`, m.Name, m.MapFile, m.OsmId, m.ID)
	return err
}
