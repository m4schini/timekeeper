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

	dbAdapter, err := adapters.NewPostgresqlDatabase()
	if err != nil {
		logger.Fatal("failed to create postgresql adapter", zap.Error(err))
	}
	db := database.New(dbAdapter)
	authy := auth.NewAuthenticator()

	pages := []www.Route{
		&p.LandingPageRoute{DB: db},
		&p.LocationPageRoute{DB: db},
		&p.DayPageRoute{DB: db},
		&p.DayMarkdownPageRoute{DB: db},
		&p.EventPageRoute{DB: db},
		&p.LoginPageRoute{Auth: authy},
		&p.LogoutRoute{},
		&p.CreateTimeslotPageRoute{DB: db},
		&p.EditTimeslotPageRoute{DB: db},
		&p.DuplicateTimeslotPageRoute{DB: db},
		&p.VocScheduleRoute{DB: db},
		&p.MarkdownPageRoute{DB: db},
		www.StaticFileRoute{},
	}
	components := []www.Route{
		&c.DayRoute{DB: db},
		&c.CreateTimeslotRoute{DB: db},
		&c.UpdateTimeslotRoute{DB: db},
		&c.DeleteTimeslotRoute{DB: db},
		&p.LoginRoute{Auth: authy},
	}

	logger.Debug("serving timekeeper :" + config.Port())
	logger.Warn("failed to serve", zap.Error(www.Serve(l, authy, pages, components)))
}
