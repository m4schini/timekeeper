package command

import (
	"context"
	"time"
)

type CreateEvent InsertHandler[CreateEventRequest, int]

type CreateEventRequest struct {
	Name     string    `json:"name"`
	Slug     string    `json:"slug"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Setup    int       `json:"setup"`
	Teardown int       `json:"teardown"`
}

type CreateEventHandler struct {
	DB Database
}

func (c *CreateEventHandler) Execute(ctx context.Context, request CreateEventRequest) (id int, err error) {
	row := c.DB.QueryRow(`
INSERT INTO raumzeitalpaka.events (name, slug, event_start, event_end, setup, teardown) 
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id`, request.Name, request.Slug, request.Start, request.End, request.Setup, request.Teardown)
	if err = row.Err(); err != nil {
		return -1, err
	}

	err = row.Scan(&id)
	return id, err
}
