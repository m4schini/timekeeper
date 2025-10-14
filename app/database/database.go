package database

import (
	"database/sql"
	"timekeeper/app/database/command"
	"timekeeper/app/database/model"
	"timekeeper/app/database/query"
)

type Database struct {
	Queries  Queries
	Commands Commands
}

func New(db *sql.DB) *Database {
	return &Database{
		Queries:  &query.Queries{DB: db},
		Commands: &command.Commands{DB: db},
	}
}

type Queries interface {
	GetLocation(id int) (l model.LocationModel, err error)
	GetLocations(offset int, limit int) (ls []model.LocationModel, err error)
	GetLocationsOfEvent(eventId int) (ls []model.EventLocationModel, err error)
	GetRoom(id int) (r model.RoomModel, err error)
	GetRooms(offset int, limit int) (rs []model.RoomModel, total int, err error)
	GetRoomsOfLocation(location int, offset int, limit int) (rs []model.RoomModel, total int, err error)
	GetRoomsOfEventLocations(event int) (rs []model.RoomModel, err error)
	GetEvent(id int) (r model.EventModel, err error)
	GetEvents(offset int, limit int) (es []model.EventModel, err error)
	GetTimeslot(id int) (r model.TimeslotModel, err error)
	GetTimeslotsOfEvent(event int, offset int, limit int) (ts []model.TimeslotModel, total int, err error)
	GetUserByLoginName(loginName string) (u model.UserModel, err error)
}

type Commands interface {
	CreateTimeslot(m model.CreateTimeslotModel) (id int, err error)
	DeleteTimeslot(id int) (err error)
	UpdateTimeslot(m model.UpdateTimeslotModel) (err error)

	CreateEvent(m model.CreateEventModel) (id int, err error)
	DeleteEvent(id int) (err error)
	UpdateEvent(m model.UpdateEventModel) (err error)

	CreateLocation(m model.CreateLocationModel) (id int, err error)
	UpdateLocation(m model.UpdateLocationModel) (err error)
	AddLocationToEvent(m model.AddLocationToEventModel) (id int, err error)
	UpdateLocationToEvent(m model.UpdateLocationToEventModel) (err error)
	DeleteLocationFromEvent(id int) (err error)

	CreateRoom(m model.CreateRoomModel) (id int, err error)
	DeleteRoom(id int) (err error)

	CreateUser(m model.CreateUserModel) (id int, err error)
}
