package command

type DeleteTimeslot Handler[DeleteTimeslotRequest]

type DeleteTimeslotRequest struct {
	TimeslotID int
}

type DeleteTimeslotHandler struct {
	DB Database
}

func (c *DeleteTimeslotHandler) Execute(request DeleteTimeslotRequest) (err error) {
	_, err = c.DB.Exec(`
DELETE FROM  raumzeitalpaka.timeslots
WHERE id = $1`, request.TimeslotID)
	return err
}
