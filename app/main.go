package main

import (
	"net"
	"timekeeper/adapters"
	"timekeeper/adapters/nominatim"
	"timekeeper/app/auth"
	"timekeeper/app/database"
	"timekeeper/config"
	"timekeeper/ports/www"
	c "timekeeper/ports/www/components"
	p "timekeeper/ports/www/pages"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

var version = "dev"

func main() {
	logger := NewLogger()
	zap.ReplaceGlobals(logger)

	logger.Info("starting timekeeper", zap.String("version", version))

	// init adapters
	nominatimClient := nominatim.New()
	dbAdapter, err := adapters.NewPostgresqlDatabase()
	if err != nil {
		logger.Fatal("failed to create postgresql adapter", zap.Error(err))
	}
	defer dbAdapter.Close()

	// init app
	db := database.New(dbAdapter)
	authy := auth.NewAuthenticator(db)
	pixelHack := www.PixelHackItems()

	// create admin user
	adminPassword := config.AdminPassword()
	if adminPassword != "" {
		id, err := authy.CreateUser("admin", config.AdminPassword())
		if err != nil {
			logger.Debug("tried to create admin user", zap.Error(err), zap.Int("user", id))
		}
	}

	// init ports/frontend
	pages := []www.Route{
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

		&www.ShortEventHandler{DB: db},
		&www.ShortEventScheduleHandler{DB: db},
		&www.ShortEventScheduleMHandler{DB: db},

		www.StaticFileRoute{},
		www.FontFileRoute{},
		www.PixelhackFileRoute{},
	}
	c.SetAvailablePixelHackIcons(pixelHack)
	components := []www.Route{
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

	l, err := net.Listen("tcp", ":"+config.Port())
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	logger.Info("serving timekeeper :" + config.Port())
	err = www.Serve(l, authy, pages, components)
	if err != nil {
		logger.Warn("failed to serve www", zap.Error(err))
	}
}

func NewLogger() *zap.Logger {
	if config.TelemetryEnabled() {
		logger, _ := zap.NewProduction()
		return logger
	}

	logger, _ := zap.NewDevelopment()
	return logger
}
