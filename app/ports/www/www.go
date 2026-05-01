package www

import (
	"raumzeitalpaka/adapters/nominatim"
	authz "raumzeitalpaka/app/auth/authz"
	"raumzeitalpaka/app/database"
	c "raumzeitalpaka/ports/www/components"
	p "raumzeitalpaka/ports/www/pages"
	r "raumzeitalpaka/ports/www/routes"
)

func NewWWWPort(
	db *database.Database,
	nominatimClient *nominatim.Client,
) (pages []Route, components []Route) {
	pixelHack := PixelHackItems()
	q := db.Queries
	cmd := db.Commands
	az := authz.NewDatabaseAuthz(db)

	// init ports/frontend
	pages = []Route{
		&p.LandingPageRoute{GetEvents: q.Events, Authz: az},

		r.EventCreateFormHandler(q.Event, cmd.CreateEvent),
		r.EventUpdateFormHandler(q.Event, q.UserHasRole, cmd.UpdateEvent),
		//&p.CreateEventPageRoute{Authz: az},
		//&p.EditEventPageRoute{GetEvent: q.Event, Authz: az},
		&p.EventPageRoute{
			GetEvent:          q.Event,
			GetEventLocations: q.EventLocations,
			GetLocations:      q.Locations,
			Nominatim:         nominatimClient,
			Authz:             az,
		},

		&p.SchedulePageRoute{
			GetEvent:            q.Event,
			GetTimeslotsOfEvent: q.TimeslotsOfEvent,
			Authz:               az,
		},
		&p.CreateTimeslotPageRoute{
			GetEvent:                 q.Event,
			GetTimeslot:              q.Timeslot,
			GetRoomsOfEventLocations: q.RoomsOfEventLocations,
			Authz:                    az,
		},
		&p.EditTimeslotPageRoute{
			GetTimeslot:              q.Timeslot,
			GetRoomsOfEventLocations: q.RoomsOfEventLocations,
			Authz:                    az,
		},
		&p.DuplicateTimeslotPageRoute{
			GetTimeslot:              q.Timeslot,
			GetRoomsOfEventLocations: q.RoomsOfEventLocations,
			Authz:                    az,
		},

		&p.EventScheduleDayRoute{
			GetEvent:            q.Event,
			GetTimeslotsOfEvent: q.TimeslotsOfEvent,
			Authz:               az,
		},
		&p.EventExportVocScheduleRoute{
			GetEvent:            q.Event,
			GetTimeslotsOfEvent: q.TimeslotsOfEvent,
		},
		&p.EventExportIcalScheduleRoute{
			GetEvent:            q.Event,
			GetTimeslotsOfEvent: q.TimeslotsOfEvent,
		},
		&p.EventsExportIcalRoute{GetEvents: q.Events},
		&p.EventScheduleExportMarkdownRoute{
			GetEvent:            q.Event,
			GetTimeslotsOfEvent: q.TimeslotsOfEvent,
		},

		&p.LocationPageRoute{
			GetEvent:           q.Event,
			GetEventLocation:   q.EventLocation,
			GetRoomsOfLocation: q.RoomsOfLocation,
			Nominatim:          nominatimClient,
		},
		&p.CreateLocationPageRoute{Authz: az},
		&p.UpdateLocationPageRoute{
			GetLocation:        q.Location,
			GetRoomsOfLocation: q.RoomsOfLocation,
			Authz:              az,
		},

		&p.CreateUserPageRoute{Authz: az},
		&p.UserAccountPageRoute{GetUser: q.User, GetMemberships: q.UserOrgs},
		&p.MeAccountPageRoute{},

		&p.PixelHackPageRoute{},
		&p.AttributionsPageRoute{},

		&ShortEventHandler{GetEventBySlug: q.EventBySlug},
		&ShortEventScheduleHandler{GetEventBySlug: q.EventBySlug},
		&ShortEventScheduleMHandler{GetEventBySlug: q.EventBySlug},

		StaticFileRoute{},
		FontFileRoute{},
		PixelhackFileRoute{},
	}
	c.SetAvailablePixelHackIcons(pixelHack)
	components = []Route{
		&c.DayRoute{
			GetEvent:            q.Event,
			GetTimeslotsOfEvent: q.TimeslotsOfEvent,
			Authz:               az,
		},

		&c.CreateLocationRoute{CreateLocation: cmd.CreateLocation,
			Authz: az},
		&c.EditLocationRoute{UpdateLocation: cmd.UpdateLocation,
			Authz: az},
		&c.AddLocationToEventRoute{AddLocationToEvent: db.Commands.AddLocationToEvent,
			Authz: az},
		&c.DeleteLocationFromEventRoute{RemoveLocationFromEvent: db.Commands.RemoveLocationFromEvent,
			Authz: az},
		&c.UpdateEventLocationRoute{UpdateLocationFromEvent: cmd.UpdateLocationFromEvent,
			Authz: az},

		&c.CreateTimeslotRoute{CreateTimeslot: cmd.CreateTimeslot,
			Authz: az},
		&c.UpdateTimeslotRoute{UpdateTimeslot: cmd.UpdateTimeslot,
			Authz: az},
		&c.DeleteTimeslotRoute{DeleteTimeslot: cmd.DeleteTimeslot,
			Authz: az},

		&c.CreateRoomRoute{CreateRoom: cmd.CreateRoom,
			Authz: az},
		&c.UpdateRoomRoute{UpdateRoom: cmd.UpdateRoom,
			Authz: az},
		&c.DeleteRoomRoute{DeleteRoom: cmd.DeleteRoom,
			Authz: az},

		&c.UserHeaderWidgetRoute{User: q.User},
	}

	return pages, components
}
