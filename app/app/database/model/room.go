package model

type RoomModel struct {
	ID          int
	GUID        string
	Location    LocationModel
	Name        string
	Description string

	LocationX int
	LocationY int
	LocationW int
	LocationH int
}

type CreateRoomModel struct {
	Location    int
	Name        string
	Description string
}

type UpdateRoomModel struct {
	ID          int
	Name        string
	Description string
}
