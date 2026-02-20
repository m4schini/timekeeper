package command

import "context"

type CreateRoom InsertHandler[CreateRoomRequest, int]

type CreateRoomRequest struct {
	Location    int
	Name        string
	Description string
}

type CreateRoomHandler struct {
	DB Database
}

func (c *CreateRoomHandler) Execute(ctx context.Context, m CreateRoomRequest) (id int, err error) {
	row := c.DB.QueryRow(`
INSERT INTO raumzeitalpaka.rooms (location, name, description, location_x, location_y, location_w, location_h) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id`, m.Location, m.Name, m.Description, 0, 0, 0, 0)
	if err = row.Err(); err != nil {
		return -1, err
	}

	err = row.Scan(&id)
	return id, err
}
