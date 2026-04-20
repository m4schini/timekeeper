package main

import (
	"context"
	"raumzeitalpaka/adapters"
	"raumzeitalpaka/adapters/nominatim"
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/config"
	"raumzeitalpaka/ports"

	"go.uber.org/zap"
)

var version = "dev"

func main() {
	ctx := context.TODO()
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
	_, err = database.InitSchema(dbAdapter, database.LatestSchema)
	if err != nil {
		logger.Fatal("failed to initiate database schema", zap.Error(err))
	}
	db := database.New(dbAdapter)
	if err != nil {
		logger.Fatal("failed to initiate login server")
	}

	port := config.Port()
	err = ports.Serve(ctx, port, db, nominatimClient)
	logger.Error("serving raumzeitalpaka", zap.String("port", port))
}

func NewLogger() *zap.Logger {
	if config.TelemetryEnabled() {
		logger, _ := zap.NewProduction()
		return logger
	}

	logger, _ := zap.NewDevelopment()
	return logger
}
