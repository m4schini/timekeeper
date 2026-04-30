package command

import (
	"context"
	"time"
)

type UpdateEvent Handler[UpdateEventRequest]

type UpdateEventRequest struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Slug     string    `json:"slug"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Setup    int       `json:"setup"`
	Teardown int       `json:"teardown"`
}

type UpdateEventHandler struct {
	DB Database
}

func (c *UpdateEventHandler) Execute(ctx context.Context, request UpdateEventRequest) (err error) {
	_, err = c.DB.Exec(`
UPDATE raumzeitalpaka.events
SET
    name = $1,
    slug = $2,
    event_start = $3,
    event_end = $4,
    setup = $5,
    teardown = $6
WHERE id = $7`, request.Name, request.Slug, request.Start, request.End, request.Setup, request.Teardown, request.ID)
	return err
}
