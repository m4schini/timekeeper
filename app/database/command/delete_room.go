package command

func (c *Commands) DeleteRoom(id int) (err error) {

	_, err = c.DB.Exec(`
DELETE FROM  timekeeper.rooms
WHERE id = $1`, id)
	return err
}
