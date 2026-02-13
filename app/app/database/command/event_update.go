package command

import "time"

type UpdateEvent Handler[UpdateEventRequest]

type UpdateEventRequest struct {
	ID    int
	Name  string
	Slug  string
	Start time.Time
}

type UpdateEventHandler struct {
	DB Database
}

func (c *UpdateEventHandler) Execute(request UpdateEventRequest) (err error) {
	_, err = c.DB.Exec(`
UPDATE raumzeitalpaka.events
SET
    name = $1,
    start = $2,
    slug = $3
WHERE id = $4`, request.Name, request.Start, request.Slug, request.ID)
	return err
}
