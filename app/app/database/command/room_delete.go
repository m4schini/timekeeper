package command

type DeleteRoom Handler[DeleteRoomRequest]

type DeleteRoomRequest struct {
	RoomID int
}

type DeleteRoomHandler struct {
	DB Database
}

func (c *DeleteRoomHandler) Execute(request DeleteRoomRequest) (err error) {
	_, err = c.DB.Exec(`
DELETE FROM raumzeitalpaka.rooms
WHERE id = $1`, request.RoomID)
	return err
}
