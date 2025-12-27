package www

import (
	"raumzeitalpaka/adapters/nominatim"
	"raumzeitalpaka/app/auth"
	"raumzeitalpaka/app/database"
	c "raumzeitalpaka/ports/www/components"
	p "raumzeitalpaka/ports/www/pages"

	"golang.org/x/time/rate"
)

func NewWWWPort(
	db *database.Database,
	nominatimClient *nominatim.Client,
	authy auth.Authenticator,
) (pages []Route, components []Route) {
	pixelHack := PixelHackItems()

	// init ports/frontend
	pages = []Route{
		&p.LandingPageRoute{DB: db},

		&p.CreateEventPageRoute{DB: db},
		&p.EditEventPageRoute{DB: db},
		&p.EventPageRoute{DB: db, Nominatim: nominatimClient},

		&p.SchedulePageRoute{DB: db},
		&p.CreateTimeslotPageRoute{DB: db},
		&p.EditTimeslotPageRoute{DB: db},
		&p.DuplicateTimeslotPageRoute{DB: db},

		&p.EventScheduleDayRoute{DB: db},
		&p.EventExportVocScheduleRoute{DB: db},
		&p.EventExportIcalScheduleRoute{DB: db},
		&p.EventsExportIcalRoute{DB: db},
		&p.EventScheduleExportMarkdownRoute{DB: db},

		&p.LocationPageRoute{DB: db, Nominatim: nominatimClient},
		&p.CreateLocationPageRoute{DB: db},
		&p.UpdateLocationPageRoute{DB: db},

		&p.CreateUserPageRoute{DB: db},
		&p.LoginPageRoute{Auth: authy},
		&p.LogoutRoute{},

		&p.PixelHackPageRoute{},
		&p.AttributionsPageRoute{},

		&ShortEventHandler{DB: db},
		&ShortEventScheduleHandler{DB: db},
		&ShortEventScheduleMHandler{DB: db},

		StaticFileRoute{},
		FontFileRoute{},
		PixelhackFileRoute{},
	}
	c.SetAvailablePixelHackIcons(pixelHack)
	components = []Route{
		&c.CreateEventRoute{DB: db},
		&c.UpdateEventRoute{DB: db},
		&c.DayRoute{DB: db},

		&c.CreateLocationRoute{DB: db},
		&c.EditLocationRoute{DB: db},
		&c.AddLocationToEventRoute{DB: db},
		&c.DeleteLocationFromEventRoute{DB: db},
		&c.UpdateEventLocationRoute{DB: db},

		&c.CreateTimeslotRoute{DB: db},
		&c.UpdateTimeslotRoute{DB: db},
		&c.DeleteTimeslotRoute{DB: db},

		&c.CreateRoomRoute{DB: db},
		&c.UpdateRoomRoute{DB: db},
		&c.DeleteRoomRoute{DB: db},

		&c.CreateUserRoute{Auth: authy},
		&p.LoginRoute{Auth: authy, RateLimiter: rate.NewLimiter(1, 1)},
	}

	return pages, components
}
