package command

func (c *Commands) DeleteLocationFromEvent(id int) (err error) {
	row := c.DB.QueryRow(`
DELETE FROM  raumzeitalpaka.event_has_location
WHERE id = $1`, id)

	return row.Err()
}
