package command

import (
	"context"
	"time"
)

type CreateEvent InsertHandler[CreateEventRequest, int]

type CreateEventRequest struct {
	Name  string    `json:"name"`
	Slug  string    `json:"slug"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type CreateEventHandler struct {
	DB Database
}

func (c *CreateEventHandler) Execute(ctx context.Context, request CreateEventRequest) (id int, err error) {
	row := c.DB.QueryRow(`
INSERT INTO raumzeitalpaka.events (name, slug, event_start, event_end) 
VALUES ($1, $2, $3, $4)
RETURNING id`, request.Name, request.Slug, request.Start, request.End)
	if err = row.Err(); err != nil {
		return -1, err
	}

	err = row.Scan(&id)
	return id, err
}
