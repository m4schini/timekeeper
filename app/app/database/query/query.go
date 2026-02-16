package query

import (
	"database/sql"
)

type Handler[R, T any] interface {
	Query(request R) (T, error)
}

type Queries struct {
	Event          GetEvent
	Events         GetEvents
	EventBySlug    GetEventBySlug
	EventLocation  GetEventLocation
	EventLocations GetEventLocations

	Location  GetLocation
	Locations GetLocations

	Room                  GetRoom
	Rooms                 GetRooms
	RoomsOfLocation       GetRoomsOfLocation
	RoomsOfEventLocations GetRoomsOfEventLocations

	Timeslot         GetTimeslot
	TimeslotsOfEvent GetTimeslotsOfEvent

	User             GetUser
	UserByLoginName  GetUserByLoginName
	UserHasRole      UserHasRole
	UserHasGroupRole UserHasGroupRole
	UserHasEventRole UserHasEventRole
	GroupBySlug      GetGroupBySlug
}

func NewQueries(db Database) Queries {
	return Queries{
		Event:          &GetEventHandler{DB: db},
		Events:         &GetEventsHandler{DB: db},
		EventBySlug:    &GetEventBySlugHandler{DB: db},
		EventLocation:  &GetEventLocationHandler{DB: db},
		EventLocations: &GetEventLocationsHandler{DB: db},

		Location:  &GetLocationHandler{DB: db},
		Locations: &GetLocationsHandler{DB: db},

		Room:                  &GetRoomHandler{DB: db},
		Rooms:                 &GetRoomsHandler{DB: db},
		RoomsOfLocation:       &GetRoomsOfLocationHandler{DB: db},
		RoomsOfEventLocations: &GetRoomsOfEventLocationsHandler{DB: db},

		Timeslot:         &GetTimeslotHandler{DB: db},
		TimeslotsOfEvent: &GetTimeslotsOfEventHandler{DB: db},

		User:             &GetUserHandler{DB: db},
		UserByLoginName:  &GetUserByLoginNameHandler{DB: db},
		UserHasRole:      &UserHasRoleHandler{DB: db},
		UserHasGroupRole: &UserHasGroupRoleHandler{DB: db},
		UserHasEventRole: &UserHasEventRoleHandler{DB: db},
	}
}

type Database interface {
	QueryRow(query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
}
