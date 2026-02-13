package command

type UpdateLocationFromEvent Handler[UpdateLocationFromEventRequest]

type UpdateLocationFromEventRequest struct {
	ID      int
	Name    string
	Note    string
	Visible bool
}

type UpdateLocationFromEventHandler struct {
	DB Database
}

func (c *UpdateLocationFromEventHandler) Execute(m UpdateLocationFromEventRequest) (err error) {
	_, err = c.DB.Exec(`
UPDATE raumzeitalpaka.event_has_location
SET
    name = $1,
    note = $2,
    visible = $3
WHERE id = $4`, m.Name, m.Note, m.Visible, m.ID)
	return err
}
