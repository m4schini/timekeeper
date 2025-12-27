package main

import (
	"net"
	"raumzeitalpaka/adapters"
	"raumzeitalpaka/adapters/nominatim"
	"raumzeitalpaka/app/auth"
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/config"
	"raumzeitalpaka/ports/www"

	"go.uber.org/zap"
)

var version = "dev"

func main() {
	logger := NewLogger()
	zap.ReplaceGlobals(logger)

	logger.Info("starting raumzeitalpaka", zap.String("version", version))
	config.Load()

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

	// create admin user
	//adminPassword := config.AdminPassword()
	//if adminPassword != "" {
	//	id, err := authy.CreateUser("admin", config.AdminPassword())
	//	if err != nil {
	//		logger.Debug("tried to create admin user", zap.Error(err), zap.Int("user", id))
	//	}
	//}

	port := config.Port()
	pages, components := www.NewWWWPort(db, nominatimClient, authy)
	logger.Info("serving raumzeitalpaka", zap.String("port", port), zap.Int("pages", len(pages)), zap.Int("components", len(components)))

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

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
