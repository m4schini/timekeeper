package main

import (
	"go.uber.org/zap"
	"net"
	"timekeeper/adapters"
	"timekeeper/app/auth"
	"timekeeper/app/database"
	"timekeeper/config"
	"timekeeper/ports/www"
	c "timekeeper/ports/www/components"
	p "timekeeper/ports/www/pages"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	l, err := net.Listen("tcp", ":"+config.Port())
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	//memcached, err := memory.NewAdapter(
	//	memory.AdapterWithAlgorithm(memory.LRU),
	//	memory.AdapterWithCapacity(10000000),
	//)
	//if err != nil {
	//	logger.Fatal("failed to create memcached adapter", zap.Error(err))
	//}
	//
	//cacheClient, err := cache.NewClient(
	//	cache.ClientWithAdapter(memcached),
	//	cache.ClientWithTTL(1*time.Minute),
	//	cache.ClientWithRefreshKey("opn"),
	//)
	//if err != nil {
	//	logger.Fatal("failed to create cache client", zap.Error(err))
	//}

	nominatimClient := adapters.NewNominatimClient()

	dbAdapter, err := adapters.NewPostgresqlDatabase()
	if err != nil {
		logger.Fatal("failed to create postgresql adapter", zap.Error(err))
	}
	db := database.New(dbAdapter)
	authy := auth.NewAuthenticator(db)

	id, err := authy.CreateUser("admin", config.AdminPassword())
	if err != nil {
		logger.Debug("tried to create admin user", zap.Error(err), zap.Int("user", id))
	}

	pages := []www.Route{
		&p.LandingPageRoute{DB: db},
		&p.LocationPageRoute{DB: db},
		&p.CreateLocationPageRoute{DB: db},
		&p.UpdateLocationPageRoute{DB: db},
		&p.EventScheduleDayRoute{DB: db},
		&p.EventPageRoute{DB: db, Nominatim: nominatimClient},
		&p.SchedulePageRoute{DB: db},
		&p.CreateEventPageRoute{DB: db},
		&p.LoginPageRoute{Auth: authy},
		&p.LogoutRoute{},
		&p.CreateTimeslotPageRoute{DB: db},
		&p.EditTimeslotPageRoute{DB: db},
		&p.DuplicateTimeslotPageRoute{DB: db},
		&p.EventExportVocScheduleRoute{DB: db},
		&p.EventScheduleExportMarkdownRoute{DB: db},
		www.StaticFileRoute{},
	}
	components := []www.Route{
		&c.DayRoute{DB: db},
		&c.CreateEventRoute{DB: db},
		&c.CreateLocationRoute{DB: db},
		&c.EditLocationRoute{DB: db},
		&c.CreateTimeslotRoute{DB: db},
		&c.UpdateTimeslotRoute{DB: db},
		&c.DeleteTimeslotRoute{DB: db},
		&c.AddLocationToEventRoute{DB: db},
		&c.DeleteLocationFromEventRoute{DB: db},
		&c.UpdateEventLocationRoute{DB: db},
		&c.CreateRoomRoute{DB: db},
		&c.DeleteRoomRoute{DB: db},
		&p.LoginRoute{Auth: authy},
	}

	logger.Debug("serving timekeeper :" + config.Port())
	logger.Warn("failed to serve", zap.Error(www.Serve(l, authy, pages, components)))
}
