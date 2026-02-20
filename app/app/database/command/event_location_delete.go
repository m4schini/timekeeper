package command

import "context"

type RemoveLocationFromEvent Handler[RemoveLocationFromEventRequest]

type RemoveLocationFromEventRequest struct {
	EventLocationRelationID int
}

type RemoveLocationFromEventHandler struct {
	DB Database
}

func (c *RemoveLocationFromEventHandler) Execute(ctx context.Context, request RemoveLocationFromEventRequest) (err error) {
	row := c.DB.QueryRow(`
DELETE FROM  raumzeitalpaka.event_has_location
WHERE id = $1`, request.EventLocationRelationID)

	return row.Err()
}
