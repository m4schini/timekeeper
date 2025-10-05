package model

type RoomModel struct {
	ID       int
	Location LocationModel
	Name     string

	LocationX int
	LocationY int
	LocationW int
	LocationH int
}
