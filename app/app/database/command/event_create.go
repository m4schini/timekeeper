package command

import "time"

type CreateEvent InsertHandler[CreateEventRequest, int]

type CreateEventRequest struct {
	Name  string
	Slug  string
	Start time.Time
}

type CreateEventHandler struct {
	DB Database
}

func (c *CreateEventHandler) Execute(request CreateEventRequest) (id int, err error) {
	row := c.DB.QueryRow(`
INSERT INTO raumzeitalpaka.events (name, slug, start) 
VALUES ($1, $2, $3)
RETURNING id`, request.Name, request.Slug, request.Start)
	if err = row.Err(); err != nil {
		return -1, err
	}

	err = row.Scan(&id)
	return id, err
}
