package command

type UpdateLocation Handler[UpdateLocationRequest]

type UpdateLocationRequest struct {
	ID      int
	Name    string
	MapFile string
	OsmId   string
}

type UpdateLocationHandler struct {
	DB Database
}

func (c *UpdateLocationHandler) Execute(m UpdateLocationRequest) (err error) {
	_, err = c.DB.Exec(`
UPDATE raumzeitalpaka.locations
SET
    name = $1,
    file = $2,
    osm_id = $3
WHERE id = $4`, m.Name, m.MapFile, m.OsmId, m.ID)
	return err
}
