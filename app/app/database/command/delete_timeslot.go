package command

func (c *Commands) DeleteTimeslot(id int) (err error) {

	_, err = c.DB.Exec(`
DELETE FROM  timekeeper.timeslots
WHERE id = $1`, id)
	return err
}
