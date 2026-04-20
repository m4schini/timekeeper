package ports

import (
	"context"
	"net"
	"net/http"
	"raumzeitalpaka/adapters/nominatim"
	"raumzeitalpaka/app/auth/dev"
	"raumzeitalpaka/app/auth/local"
	"raumzeitalpaka/app/auth/oidc"
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/config"
	"raumzeitalpaka/ports/api"
	"raumzeitalpaka/ports/www"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Serve(ctx context.Context, port string, db *database.Database, nominatim *nominatim.Client) error {
	logger := zap.L()
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	defer l.Close()

	// web ui
	logger.Debug("initiating auth provider")
	var authHandler chi.Router
	oidcCfg, oidcEnabled := config.OIDCProviderConfig()
	switch {
	case config.DevAuthEnabled():
		logger.Info("using dev auth provider")
		authy := local.NewAuthenticator(db)
		authHandler, err = dev.NewHandler(db.Commands.InsertUser, authy)
		break
	case oidcEnabled:
		logger.Info("using oidc auth provider", zap.Any("issuer", oidcCfg.IssuerURL), zap.String("callbackPath", oidc.CallbackPath))
		syncer := oidc.NewAlpakaSyncer(db)
		authHandler, err = oidc.NewHandler(ctx, oidcCfg, syncer, db.Commands.UpdateLastLogin)
		break
	default:
		logger.Info("using local auth provider")
		authy := local.NewAuthenticator(db)
		authHandler, err = local.NewHandler(authy)
	}

	pages, components := www.NewWWWPort(db, nominatim)
	wwwRouter, err := www.NewRouter(authHandler, pages, components)
	if err != nil {
		return err
	}

	// api
	apiRouter := api.NewRouter(db)

	// router
	router := chi.NewRouter()
	router.Mount("/", wwwRouter)
	router.Mount("/api", apiRouter)
	return http.Serve(l, router)
}
