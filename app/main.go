package main

import (
	"context"
	"net"
	"raumzeitalpaka/adapters"
	"raumzeitalpaka/adapters/nominatim"
	"raumzeitalpaka/app/auth/local"
	"raumzeitalpaka/app/auth/oidc"
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/config"
	"raumzeitalpaka/ports/www"

	"github.com/go-chi/chi/v5"
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
	_, err = database.InitSchema(dbAdapter)
	if err != nil {
		logger.Fatal("failed to initiate database schema", zap.Error(err))
	}
	db := database.New(dbAdapter)
	//authy := local.NewAuthenticator(db)
	if err != nil {
		logger.Fatal("failed to initiate login server")
	}

	// auth
	logger.Debug("initiating auth provider")
	var authHandler chi.Router
	oidcCfg, oidcEnabled := config.OIDCProviderConfig()
	if oidcEnabled {
		logger.Info("using oidc auth provider", zap.Any("issuer", oidcCfg.IssuerURL), zap.String("callbackPath", oidc.CallbackPath))
		authHandler, err = oidc.NewHandler(ctx, oidcCfg)
	} else {
		logger.Info("using local auth provider")
		authy := local.NewAuthenticator(db)
		authHandler, err = local.NewHandler(authy)
	}
	if err != nil {
		logger.Fatal("failed to initiate auth", zap.Error(err), zap.Bool("oidc", oidcEnabled))
	}

	port := config.Port()
	pages, components := www.NewWWWPort(db, nominatimClient)
	logger.Info("serving raumzeitalpaka", zap.String("port", port), zap.Int("pages", len(pages)), zap.Int("components", len(components)))

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	err = www.Serve(l, authHandler, pages, components)
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
