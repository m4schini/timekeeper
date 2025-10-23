package command

func (c *Commands) DeleteEvent(id int) (err error) {

	_, err = c.DB.Exec(`
DELETE FROM  timekeeper.events
WHERE id = $1`, id)
	return err
}
