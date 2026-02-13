package command

type UpdateRoom Handler[UpdateRoomRequest]

type UpdateRoomRequest struct {
	ID          int
	Name        string
	Description string
}

type UpdateRoomHandler struct {
	DB Database
}

func (c *UpdateRoomHandler) Execute(m UpdateRoomRequest) (err error) {
	_, err = c.DB.Exec(`
UPDATE raumzeitalpaka.rooms
SET
    name = $1,
    description = $2
WHERE id = $3`, m.Name, m.Description, m.ID)
	return err
}
