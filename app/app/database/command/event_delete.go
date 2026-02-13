package command

type DeleteEvent Handler[DeleteEventRequest]

type DeleteEventRequest struct {
	EventID int
}

type DeleteEventHandler struct {
	DB Database
}

func (c *DeleteEventHandler) Execute(request DeleteEventRequest) (err error) {
	_, err = c.DB.Exec(`
DELETE FROM  raumzeitalpaka.events
WHERE id = $1`, request.EventID)
	return err
}
