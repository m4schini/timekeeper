package command

import "database/sql"

type InsertHandler[R, T any] interface {
	Execute(request R) (T, error)
}

type Handler[R any] interface {
	Execute(request R) error
}

type Commands struct {
	CreateEvent CreateEvent
	UpdateEvent UpdateEvent
	DeleteEvent DeleteEvent

	AddLocationToEvent      AddLocationToEvent
	UpdateLocationFromEvent UpdateLocationFromEvent
	RemoveLocationFromEvent RemoveLocationFromEvent

	CreateLocation CreateLocation
	UpdateLocation UpdateLocation

	CreateRoom CreateRoom
	DeleteRoom DeleteRoom
	UpdateRoom UpdateRoom

	CreateTimeslot CreateTimeslot
	UpdateTimeslot UpdateTimeslot
	DeleteTimeslot DeleteTimeslot

	InsertUser                     UpsertUser
	CreateUser                     CreateUser
	UpdateLastLogin                UpdateLastLogin
	CreateGroup                    CreateGroup
	GroupAddUser                   GroupAddUser
	UpdateManagedGroupsAssignments UpdateManagedGroupsAssignments
}

func NewCommands(db Database) Commands {
	return Commands{
		CreateEvent: &CreateEventHandler{DB: db},
		UpdateEvent: &UpdateEventHandler{DB: db},
		DeleteEvent: &DeleteEventHandler{DB: db},

		AddLocationToEvent:      &AddLocationToEventHandler{DB: db},
		UpdateLocationFromEvent: &UpdateLocationFromEventHandler{DB: db},
		RemoveLocationFromEvent: &RemoveLocationFromEventHandler{DB: db},

		CreateLocation: &CreateLocationHandler{DB: db},
		UpdateLocation: &UpdateLocationHandler{DB: db},

		CreateRoom: &CreateRoomHandler{DB: db},
		DeleteRoom: &DeleteRoomHandler{DB: db},
		UpdateRoom: &UpdateRoomHandler{DB: db},

		CreateTimeslot: &CreateTimeslotHandler{DB: db},
		UpdateTimeslot: &UpdateTimeslotHandler{DB: db},
		DeleteTimeslot: &DeleteTimeslotHandler{DB: db},

		InsertUser:      &UpsertUserHandler{DB: db},
		CreateUser:      &CreateUserHandler{DB: db},
		UpdateLastLogin: &UpdateLastLoginHandler{DB: db},
		CreateGroup:     &CreateGroupHandler{DB: db},
		GroupAddUser:    &GroupAddUserHandler{DB: db},
		UpdateManagedGroupsAssignments: &UpdateManagedGroupsAssignmentsHandler{
			DB: db,
		},
	}
}

type Database interface {
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	Begin() (*sql.Tx, error)
}
