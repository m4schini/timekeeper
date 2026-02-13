package command

type AddLocationToEvent InsertHandler[AddLocationToEventRequest, int]

type AddLocationToEventRequest struct {
	Name       string
	EventId    int
	LocationId int
	Note       string
}

type AddLocationToEventHandler struct {
	DB Database
}

func (c *AddLocationToEventHandler) Execute(m AddLocationToEventRequest) (id int, err error) {
	row := c.DB.QueryRow(`
INSERT INTO raumzeitalpaka.event_has_location (name, event, location, note) 
VALUES ($1, $2, $3, $4)
RETURNING id`, m.Name, m.EventId, m.LocationId, m.Note)
	if err = row.Err(); err != nil {
		return -1, err
	}

	err = row.Scan(&id)
	return id, err
}
