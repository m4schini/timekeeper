package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
)

func EditLocationRooms(rooms []model.RoomModel) Node {
	//g := Group{}

	return Div()
}

type EditLocationRoomRoute struct {
	DB *database.Database
}
