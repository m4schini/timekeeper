package command

func (c *Commands) DeleteEvent(id int) (err error) {

	_, err = c.DB.Exec(`
DELETE FROM  raumzeitalpaka.events
WHERE id = $1`, id)
	return err
}
