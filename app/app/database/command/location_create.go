package command

type CreateLocation InsertHandler[CreateLocationRequest, int]

type CreateLocationRequest struct {
	Name    string
	MapFile string
	OsmId   string
}

type CreateLocationHandler struct {
	DB Database
}

func (c *CreateLocationHandler) Execute(m CreateLocationRequest) (id int, err error) {
	row := c.DB.QueryRow(`
INSERT INTO raumzeitalpaka.locations (name, file, osm_id) 
VALUES ($1, $2, $3)
RETURNING id`, m.Name, m.MapFile, m.OsmId)
	if err = row.Err(); err != nil {
		return -1, err
	}

	err = row.Scan(&id)
	return id, err
}
