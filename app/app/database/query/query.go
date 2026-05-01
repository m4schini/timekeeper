package query

import (
	"context"
	"database/sql"
)

type Handler[R, T any] interface {
	Query(ctx context.Context, request R) (T, error)
}

type Queries struct {
	Event                GetEvent
	Events               GetEvents
	EventsByOrganisation GetEventsByOrganisation
	EventBySlug          GetEventBySlug
	EventLocation        GetEventLocation
	EventLocations       GetEventLocations

	Location  GetLocation
	Locations GetLocations

	Room                  GetRoom
	Rooms                 GetRooms
	RoomsOfLocation       GetRoomsOfLocation
	RoomsOfEventLocations GetRoomsOfEventLocations

	Timeslot         GetTimeslot
	TimeslotsOfEvent GetTimeslotsOfEvent

	User                    GetUser
	UserByLoginName         GetUserByLoginName
	UserHasRole             UserHasRole
	UserHasOrganisationRole UserHasOrganisationRole
	UserHasEventRole        UserHasEventRole
	UserOrgs                GetUserOrganisations

	Organisation        GetOrganisation
	OrganisationMembers GetOrganisationMembers
	OrganisationBySlug  GetOrganisationBySlug
}

func NewQueries(db Database) Queries {
	return Queries{
		Event:                &GetEventHandler{DB: db},
		Events:               &GetEventsHandler{DB: db},
		EventsByOrganisation: &GetEventsByOrganisationHandler{DB: db},
		EventBySlug:          &GetEventBySlugHandler{DB: db},
		EventLocation:        &GetEventLocationHandler{DB: db},
		EventLocations:       &GetEventLocationsHandler{DB: db},

		Location:  &GetLocationHandler{DB: db},
		Locations: &GetLocationsHandler{DB: db},

		Room:                  &GetRoomHandler{DB: db},
		Rooms:                 &GetRoomsHandler{DB: db},
		RoomsOfLocation:       &GetRoomsOfLocationHandler{DB: db},
		RoomsOfEventLocations: &GetRoomsOfEventLocationsHandler{DB: db},

		Timeslot:         &GetTimeslotHandler{DB: db},
		TimeslotsOfEvent: &GetTimeslotsOfEventHandler{DB: db},

		User:                    &GetUserHandler{DB: db},
		UserByLoginName:         &GetUserByLoginNameHandler{DB: db},
		UserHasRole:             &UserHasRoleHandler{DB: db},
		UserHasOrganisationRole: &UserHasOrganisationRoleHandler{DB: db},
		UserHasEventRole:        &UserHasEventRoleHandler{DB: db},
		UserOrgs:                &GetUserOrganisationsHandler{DB: db},

		Organisation:        &GetOrganisationHandler{DB: db},
		OrganisationMembers: &GetOrganisationMembersHandler{DB: db},
		OrganisationBySlug:  &GetOrganisationBySlugHandler{DB: db},
	}
}

type Database interface {
	QueryRow(query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
}
